package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
)

var (
	keyLUIssue = []byte("last-updated-issue")
)

func (b *boltSource) LastUpdatedIssue() (*github.Issue, error) {
	var id *int

	err := b.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(bktMisc).Get(keyLUIssue)

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

	var iss github.Issue

	err = b.db.View(func(tx *bolt.Tx) error {
		k := encodeInt(*id)
		v := tx.Bucket(bktIssues).Get(k)

		return decodeVal(v, &iss)
	})

	if err != nil {
		return nil, err
	}

	return &iss, nil
}

func (b *boltSource) UpdateIssue(iss *github.Issue) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		k := encodeInt(*iss.ID)

		v, err := encodeVal(iss)

		if err != nil {
			return err
		}

		if err := tx.Bucket(bktIssues).Put(k, v); err != nil {
			return err
		}

		l, err := b.LastUpdatedIssue()

		if err != nil {
			return err
		}

		if l == nil || l.UpdatedAt.Before(*iss.UpdatedAt) {
			err = tx.Bucket(bktMisc).Put(keyLUIssue, k)

			if err != nil {
				return err
			}
		}

		return nil
	})
}
