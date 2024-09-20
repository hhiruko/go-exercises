package handler

import (
	"encoding/json"
	"net/http"
)

type Json struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//		[
//	  	{
//	 		"path": "/some-path"
//	   		"url": "https://www.some-url.com/demo"
//	  	}
//		]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	pathMap := mapJSON(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(jsonByte []byte) ([]Json, error) {
	var j []Json
	err := json.Unmarshal(jsonByte, &j)
	return j, err
}

func mapJSON(parsedJson []Json) map[string]string {
	jsonMap := make(map[string]string)
	for _, y := range parsedJson {
		jsonMap[y.Path] = y.Url
	}
	return jsonMap
}
