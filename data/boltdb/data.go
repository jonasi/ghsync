package boltdb

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"github.com/jonasi/ghsync"
)

var (
	bktRepos                     = []byte("repos")
	bktMembers                   = []byte("members")
	bktTeams                     = []byte("teams")
	bktIssues                    = []byte("issues")
	bktIssueComments             = []byte("issue-comments")
	bktIssueCommentsLastUpdated  = []byte("issue-comments-last-updated")
	bktReviewComments            = []byte("review-comments")
	bktReviewCommentsLastUpdated = []byte("review-comments-last-updated")
	bktCommitComments            = []byte("commit-comments")
	bktMisc                      = []byte("misc")
	allBuckets                   = [][]byte{bktRepos, bktMembers, bktTeams, bktIssues, bktMisc, bktIssueComments, bktIssueCommentsLastUpdated, bktReviewComments, bktReviewCommentsLastUpdated, bktCommitComments}
)

func New(path string) (ghsync.DataStore, error) {
	db, err := bolt.Open(path, 0644, nil)

	if err != nil {
		return nil, err
	}

	b := &boltSource{db}

	if err := b.init(); err != nil {
		return nil, err
	}

	return b, nil
}

type boltSource struct {
	db *bolt.DB
}

func (b *boltSource) init() error {
	return b.db.Update(func(tx *bolt.Tx) error {
		for _, bk := range allBuckets {
			_, err := tx.CreateBucketIfNotExists(bk)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func encodeInt(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func decodeInt(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func encodeVal(v interface{}) ([]byte, error) {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(v); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func decodeVal(b []byte, dest interface{}) error {
	return gob.NewDecoder(bytes.NewReader(b)).Decode(dest)
}
