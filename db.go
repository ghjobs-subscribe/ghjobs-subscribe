package main

import (
	"fmt"
	"os"
	"path"
	"time"

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
	i.DB, err = bolt.Open(ghjsDBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logrus.Fatal(err)
	}
}

// frequence: weekly (1), fortnightly (2), monthly (3)
func (i *impl) createUserBucket(email string) bool {
	err := i.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(email))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		err = b.Put([]byte("email"), []byte(email))
		if err != nil {
			return err
		}
		err = b.Put([]byte("firstname"), []byte(""))
		if err != nil {
			return err
		}
		err = b.Put([]byte("lastname"), []byte(""))
		if err != nil {
			return err
		}
		err = b.Put([]byte("ghjobstag"), []byte(""))
		if err != nil {
			return err
		}
		err = b.Put([]byte("ghjobslocation"), []byte(""))
		if err != nil {
			return err
		}
		err = b.Put([]byte("frequence"), []byte("1"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("lastsent"), []byte(""))
		if err != nil {
			return err
		}
		err = b.Put([]byte("nextsend"), []byte(""))
		if err != nil {
			return err
		}
		err = b.Put([]byte("createdon"), []byte(time.Now().String()))
		if err != nil {
			return err
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
