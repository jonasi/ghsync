package main

import (
	"flag"
	"github.com/jonasi/ghsync"
	"github.com/jonasi/ghsync/data/boltdb"
	"log"
)

func main() {
	var (
		c      ghsync.Config
		dbPath string
	)

	flag.StringVar(&c.Token, "token", "", "")
	flag.StringVar(&c.Organization, "org", "", "")
	flag.StringVar(&c.ListenAddr, "listen", ":8080", "")
	flag.StringVar(&c.PublicAddr, "public", "", "")
	flag.StringVar(&dbPath, "db", "", "")

	flag.Parse()

	d, err := boltdb.New(dbPath)

	if err != nil {
		log.Fatalf("Db error: %s", err)
	}

	s := ghsync.New(c, d)

	if err := s.Start(); err != nil {
		log.Fatalf("Start err: %s", err)
	}

	if err := s.Wait(); err != nil {
		log.Fatalf("Serve err: %s", err)
	}
}
