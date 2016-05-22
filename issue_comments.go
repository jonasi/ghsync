package ghsync

import (
	"github.com/google/go-github/github"
	"time"
)

func (s *Syncer) syncIssuesComments() error {
	return s.data.ForEachRepo(s.syncIssuesCommentsByRepo)
}

func (s *Syncer) syncIssuesCommentsByRepo(r *github.Repository) error {
	last, err := s.data.LastUpdatedIssueComment(*r.Name)

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
		comments, resp, err := s.client.Issues.ListComments(s.conf.Organization, *r.Name, 0, &github.IssueListCommentsOptions{
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

		for i := range comments {
			s.logger.Debug().Log("msg", "updating issue comment", "repo", *r.Name, "id", comments[i].ID)

			if err := s.data.UpdateIssueComment(&comments[i]); err != nil {
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
