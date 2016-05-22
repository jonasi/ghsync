package ghsync

import (
	"github.com/google/go-github/github"
	"time"
)

func (s *Syncer) syncIssues() error {
	last, err := s.data.LastUpdatedIssue()

	if err != nil {
		return err
	}

	var (
		page  = 1
		size  = 100
		since time.Time
	)

	if last != nil {
		since = *last.UpdatedAt
	}

	for {
		issues, resp, err := s.client.Issues.ListByOrg(s.conf.Organization, &github.IssueListOptions{
			Filter:    "all",
			State:     "all",
			Sort:      "updated",
			Direction: "asc",
			Since:     since,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: size,
			},
		})

		if err != nil {
			return err
		}

		for i := range issues {
			s.logger.Debug().Log("msg", "updating issue", "resp", *issues[i].Repository.Name, "issue", *issues[i].Number)

			if err := s.data.UpdateIssue(&issues[i]); err != nil {
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
