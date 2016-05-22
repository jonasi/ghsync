package ghsync

import (
	"github.com/google/go-github/github"
)

func (s *Syncer) syncTeams() error {
	var (
		page = 1
		size = 100
	)

	for {
		teams, resp, err := s.client.Organizations.ListTeams(s.conf.Organization, &github.ListOptions{
			Page:    page,
			PerPage: size,
		})

		if err != nil {
			return err
		}

		for i := range teams {
			s.logger.Debug().Log("msg", "updating team", "team", *teams[i].Name)

			if err := s.data.UpdateTeam(&teams[i]); err != nil {
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
