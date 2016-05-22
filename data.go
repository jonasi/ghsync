package ghsync

import (
	"github.com/google/go-github/github"
)

type DataStore interface {
	UpdateRepo(*github.Repository) error
	ForEachRepo(func(*github.Repository) error) error
	UpdateUser(*github.User) error
	UpdateTeam(*github.Team) error
	UpdateIssue(*github.Issue) error
	LastUpdatedIssue() (*github.Issue, error)
	LastUpdatedIssueComment(repo string) (*github.IssueComment, error)
	UpdateIssueComment(*github.IssueComment) error
	LastUpdatedReviewComment(repo string) (*github.PullRequestComment, error)
	UpdateReviewComment(*github.PullRequestComment) error
	UpdateCommitComment(*github.RepositoryComment) error
}
