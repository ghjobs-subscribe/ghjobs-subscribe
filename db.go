package main

import (
	"fmt"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

type impl struct {
	DB *bolt.DB
}

func (i *impl) initDB() {
	gopath := os.Getenv("GOPATH")
	ghjsDBPath := path.Join(gopath, "src", "github.com", "ghjobs-subscribe", "ghjobs-subscribe", "subs.db")
	var err error
	i.DB, err = bolt.Open(ghjsDBPath, 0600, nil)
	if err != nil {
		logrus.Fatal(err)
	}
}

func (i *impl) createUserBucket(email string) bool {
	err := i.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(email))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return false
	}
	return true
}

func (i *impl) checkBucketExists(email string) bool {
	err := i.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(email))
		if b != nil {
			return fmt.Errorf("email exists")
		}
		return nil
	})
	if err != nil {
		return false
	}
	return true
}
