package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
)

func (b *boltSource) UpdateRepo(r *github.Repository) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		k := encodeInt(*r.ID)

		v, err := encodeVal(r)

		if err != nil {
			return err
		}

		return tx.Bucket(bktRepos).Put(k, v)
	})
}

func (b *boltSource) ForEachRepo(fn func(r *github.Repository) error) error {
	var repos []*github.Repository

	err := b.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bktRepos).ForEach(func(k, v []byte) error {
			var r github.Repository

			if err := decodeVal(v, &r); err != nil {
				return err
			}

			repos = append(repos, &r)
			return nil
		})
	})

	if err != nil {
		return err
	}

	for _, r := range repos {
		if err := fn(r); err != nil {
			return err
		}
	}

	return nil
}
