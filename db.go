package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"
)

type impl struct {
	DB *bolt.DB
}

func (i *impl) initDB() error {
	home := os.Getenv("HOME")
	ghjsDBPath := path.Join(home, "ghjs-data", "subs.db")
	var err error
	i.DB, err = bolt.Open(ghjsDBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("initDB: %v", err)
	}
	return nil
}

// TODO: parse date from time.Now().String()
func (i *impl) createUserProfile(email string) error {
	err := i.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(email))
		if err != nil {
			return fmt.Errorf("createUserProfile: %v", err)
		}

		err = b.Put([]byte("userEmail"), []byte(email))
		if err != nil {
			return err
		}
		err = b.Put([]byte("userFirstName"), []byte(""))
		if err != nil {
			return err
		}
		err = b.Put([]byte("userLastName"), []byte(""))
		if err != nil {
			return err
		}
		err = b.Put([]byte("userActive"), []byte("false"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("subTag"), []byte(""))
		if err != nil {
			return err
		}
		err = b.Put([]byte("subLocation"), []byte(""))
		if err != nil {
			return err
		}
		err = b.Put([]byte("emailLastSent"), []byte(""))
		if err != nil {
			return err
		}
		err = b.Put([]byte("userCreatedOn"), []byte(time.Now().Format("2006-01-02 15:04:05")))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (i *impl) checkUserExists(email string) bool {
	userExists := false
	i.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(email))
		if b != nil {
			userExists = true
		}
		return nil
	})
	return userExists
}

func (i *impl) checkUserSubscription(email string) bool {
	userActive := false
	i.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(email))
		ua := b.Get([]byte("userActive"))
		if string(ua) == "true" {
			userActive = true
		}
		return nil
	})
	return userActive
}

func (i *impl) changeUserSubscription(email, userActiveValue string) error {
	err := i.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(email))
		err := b.Put([]byte("userActive"), []byte(userActiveValue))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
