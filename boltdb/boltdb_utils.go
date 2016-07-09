package boltdb

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

//--------------------------------------------------------------------------------------------------

// NewDB xxx
func NewDB() *bolt.DB {
	log.Printf("Creating BoltDB...")
	bolt, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return bolt
}

//--------------------------------------------------------------------------------------------------

// CreateBucket creates a new bucket if it does not yet exist.
func CreateBucket(db *bolt.DB, name []byte) {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

// SetKeyValue adds a keys to the specified bucket
func SetKeyValue(db *bolt.DB, bucket []byte, key []byte, value []byte) {
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		err = bucket.Put(key, value)
		return err
	})
}

// DeleteKeyValue deletes key from the specified bucket
func DeleteKeyValue(db *bolt.DB, bucket []byte, key []byte) error {
	var err error
	db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(bucket)
		if bucket != nil {
			err = bucket.Delete(key)
		}
		return nil
	})
	return err
}

// GetValue get the value of the specified bucket and key
func GetValue(db *bolt.DB, bucket []byte, key []byte) (value []byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %v not found!", bucket)
		}
		var buffer bytes.Buffer
		buffer.Write(bucket.Get(key))
		value = buffer.Bytes()
		return nil
	})
	return
}
