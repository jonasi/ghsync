package ghsync

import (
	"github.com/google/go-github/github"
	"time"
)

func (s *Syncer) syncReviewComments() error {
	return s.data.ForEachRepo(s.syncReviewCommentsByRepo)
}

func (s *Syncer) syncReviewCommentsByRepo(r *github.Repository) error {
	last, err := s.data.LastUpdatedReviewComment(*r.Name)

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
		comments, resp, err := s.client.PullRequests.ListComments(s.conf.Organization, *r.Name, 0, &github.PullRequestListCommentsOptions{
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
			s.logger.Debug().Log("msg", "updating review comment", "repo", *r.Name, "id", comments[i].ID)

			if err := s.data.UpdateReviewComment(&comments[i]); err != nil {
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
