package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
)

func (b *boltSource) UpdateTeam(t *github.Team) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		k := encodeInt(*t.ID)

		v, err := encodeVal(t)

		if err != nil {
			return err
		}

		return tx.Bucket(bktTeams).Put(k, v)
	})
}
