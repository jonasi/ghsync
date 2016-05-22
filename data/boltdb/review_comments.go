package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
)

func (b *boltSource) LastUpdatedReviewComment(repo string) (*github.PullRequestComment, error) {
	var id *int

	err := b.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(bktReviewCommentsLastUpdated).Get([]byte(repo))

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

	var c github.PullRequestComment

	err = b.db.View(func(tx *bolt.Tx) error {
		k := encodeInt(*id)
		v := tx.Bucket(bktReviewComments).Get(k)

		return decodeVal(v, &c)
	})

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (b *boltSource) UpdateReviewComment(c *github.PullRequestComment) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		k := encodeInt(*c.ID)

		v, err := encodeVal(c)

		if err != nil {
			return err
		}

		if err := tx.Bucket(bktReviewComments).Put(k, v); err != nil {
			return err
		}

		repo := repoURLRegexp.FindStringSubmatch(*c.URL)[1]

		l, err := b.LastUpdatedReviewComment(repo)

		if err != nil {
			return err
		}

		if l == nil || l.UpdatedAt.Before(*c.UpdatedAt) {
			err := tx.Bucket(bktReviewCommentsLastUpdated).Put([]byte(repo), k)

			if err != nil {
				return err
			}
		}

		return nil
	})
}
