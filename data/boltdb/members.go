package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
)

func (b *boltSource) UpdateUser(u *github.User) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		k := encodeInt(*u.ID)

		v, err := encodeVal(u)

		if err != nil {
			return err
		}

		return tx.Bucket(bktMembers).Put(k, v)
	})
}
