package handler

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, ok := pathsToUrls[r.URL.Path]
		if ok {
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func Handler(filePath string, fallback http.Handler) (http.HandlerFunc, error) {
	ext := filepath.Ext(filePath)
	bytes := readFile(filePath)

	if ext == ".yaml" || ext == ".yml" {
		return YAMLHandler(bytes, fallback)
	}

	if ext == ".json" {
		return JSONHandler(bytes, fallback)
	}

	return nil, errors.New("unknown file extension")
}

func readFile(filePath string) []byte {
	f, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Can't read file "+filePath, err)
	}
	return f
}
