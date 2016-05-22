package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"regexp"
)

func (b *boltSource) LastUpdatedIssueComment(repo string) (*github.IssueComment, error) {
	var id *int

	err := b.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(bktIssueCommentsLastUpdated).Get([]byte(repo))

		if v == nil {
			return nil
		}

		id2 := decodeInt(v)
		id = &id2
		return nil
	})

	if err != nil {
		return nil, err
	}

	if id == nil {
		return nil, nil
	}

	var c github.IssueComment

	err = b.db.View(func(tx *bolt.Tx) error {
		k := encodeInt(*id)
		v := tx.Bucket(bktIssueComments).Get(k)

		return decodeVal(v, &c)
	})

	if err != nil {
		return nil, err
	}

	return &c, nil
}

var repoURLRegexp = regexp.MustCompile("^https:\\/\\/api\\.github\\.com\\/repos\\/[^\\/]+\\/([^\\/]+)\\/.*$")

func (b *boltSource) UpdateIssueComment(c *github.IssueComment) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		k := encodeInt(*c.ID)

		v, err := encodeVal(c)

		if err != nil {
			return err
		}

		if err := tx.Bucket(bktIssueComments).Put(k, v); err != nil {
			return err
		}

		repo := repoURLRegexp.FindStringSubmatch(*c.IssueURL)[1]

		l, err := b.LastUpdatedIssueComment(repo)

		if err != nil {
			return err
		}

		if l == nil || l.UpdatedAt.Before(*c.UpdatedAt) {
			err := tx.Bucket(bktIssueCommentsLastUpdated).Put([]byte(repo), k)

			if err != nil {
				return err
			}
		}

		return nil
	})
}
