package handler

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"net/http"
)

func BOLTHandler(fallback http.Handler) (http.HandlerFunc, error) {
	db, err := bolt.Open("url_shortener.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	propagate(db)
	return MapHandler(read(db), fallback), nil
}

func propagate(db *bolt.DB) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("path_urls"))
		if bucket == nil {
			var err error
			bucket, err = tx.CreateBucket([]byte("path_urls"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			pathUrls := map[string]string{
				"/bolt":   "https://github.com/boltdb/bolt",
				"/boltdb": "https://github.com/boltdb",
			}

			for path, url := range pathUrls {
				err := bucket.Put([]byte(path), []byte(url))
				if err != nil {
					return fmt.Errorf("failed to insert path %s: %v", path, err)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("could not propagate path_urls: %v", err)
	}
}

func read(db *bolt.DB) map[string]string {
	pathMap := map[string]string{}
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("path_urls"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			pathMap[string(key)] = string(value)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("could not read path_urls: %v", err)
	}

	return pathMap
}
