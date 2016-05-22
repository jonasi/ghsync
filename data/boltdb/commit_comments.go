package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
)

func (b *boltSource) UpdateCommitComment(c *github.RepositoryComment) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		k := encodeInt(*c.ID)

		v, err := encodeVal(c)

		if err != nil {
			return err
		}

		return tx.Bucket(bktCommitComments).Put(k, v)
	})
}
