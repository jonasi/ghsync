package ghsync

import (
	"github.com/google/go-github/github"
)

func (s *Syncer) syncMembers() error {
	var (
		page = 1
		size = 100
	)

	for {
		members, resp, err := s.client.Organizations.ListMembers(s.conf.Organization, &github.ListMembersOptions{
			PublicOnly: false,
			Filter:     "all",
			Role:       "all",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: size,
			},
		})

		if err != nil {
			return err
		}

		for i := range members {
			s.logger.Debug().Log("msg", "updating member", "member", *members[i].Login)

			if err := s.data.UpdateUser(&members[i]); err != nil {
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
