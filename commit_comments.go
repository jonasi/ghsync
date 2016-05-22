package ghsync

import (
	"github.com/google/go-github/github"
)

func (s *Syncer) syncCommitComments() error {
	return s.data.ForEachRepo(s.syncCommitCommentsByRepo)
}

func (s *Syncer) syncCommitCommentsByRepo(r *github.Repository) error {
	var (
		page = 1
		size = 100
	)

	for {
		comments, resp, err := s.client.Repositories.ListComments(s.conf.Organization, *r.Name, &github.ListOptions{
			Page:    page,
			PerPage: size,
		})

		if err != nil {
			return err
		}

		for i := range comments {
			s.logger.Debug().Log("msg", "updating commit comment", "repo", *r.Name, "id", comments[i].ID)

			if err := s.data.UpdateCommitComment(&comments[i]); err != nil {
				return err
			}
		}

		if resp.NextPage == 0 {
			break
		}

		page++
	}

	return nil
}
