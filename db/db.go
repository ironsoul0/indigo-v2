package db

import (
	"log"

	"github.com/boltdb/bolt"
)

func New(path string) *bolt.DB {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal("could not connect to database")
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(USERS_BUCKET))
		if err != nil {
			return err
		}
		return nil
	})

	return db
}
