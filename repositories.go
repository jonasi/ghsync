package ghsync

import (
	"github.com/google/go-github/github"
)

func (s *Syncer) syncRepos() error {
	var (
		page = 1
		size = 100
	)

	for {
		repos, resp, err := s.client.Repositories.ListByOrg(s.conf.Organization, &github.RepositoryListByOrgOptions{
			Type: "all",
			ListOptions: github.ListOptions{
				PerPage: size,
				Page:    page,
			},
		})

		if err != nil {
			return err
		}

		for i := range repos {
			s.logger.Debug().Log("msg", "updating repo", "repo", *repos[i].Name)

			if err := s.data.UpdateRepo(&repos[i]); err != nil {
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
