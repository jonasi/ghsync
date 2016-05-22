package ghsync

import (
	"github.com/facebookgo/httpdown"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/levels"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httputil"
	"os"
)

const (
	HookName = "ghsync_hook"
	hookPath = "/webhook"
)

var httpConfig = httpdown.HTTP{}

type Config struct {
	Token        string
	Organization string
	ListenAddr   string
	PublicAddr   string
}

func New(conf Config, data DataStore) *Syncer {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: conf.Token,
	})

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return &Syncer{
		conf:   conf,
		client: github.NewClient(tc),
		data:   data,
		logger: levels.New(log.NewContext(log.NewLogfmtLogger(os.Stderr)).With("ts", log.DefaultTimestamp)),
	}
}

// Syncer dog
type Syncer struct {
	conf   Config
	client *github.Client
	http   httpdown.Server
	data   DataStore
	logger levels.Levels
}

func (s *Syncer) Start() error {
	if err := s.syncRepos(); err != nil {
		return err
	}

	if err := s.syncMembers(); err != nil {
		return err
	}

	if err := s.syncTeams(); err != nil {
		return err
	}

	if err := s.syncIssues(); err != nil {
		return err
	}

	if err := s.syncIssuesComments(); err != nil {
		return err
	}

	if err := s.syncReviewComments(); err != nil {
		return err
	}

	if err := s.syncCommitComments(); err != nil {
		return err
	}

	srv := http.Server{
		Addr:    s.conf.ListenAddr,
		Handler: http.HandlerFunc(s.handleHook),
	}

	hsrv, err := httpConfig.ListenAndServe(&srv)

	if err != nil {
		return err
	}

	s.logger.Info().Log("msg", "listening", "addr", s.conf.ListenAddr)

	s.http = hsrv

	return nil
}

func (s *Syncer) Wait() error {
	return s.http.Wait()
}

func (s *Syncer) Stop() error {
	return s.http.Stop()
}

func (s *Syncer) handleHook(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != hookPath {
		http.NotFound(w, r)
		return
	}

	var (
		ev = r.Header.Get("X-Github-Event")
		// dest interface{}
	)

	switch ev {
	default:
		b, _ := httputil.DumpRequest(r, true)
		s.logger.Warn().Log("msg", "Unhandle event", "event", ev, "request", string(b))
	}
}
