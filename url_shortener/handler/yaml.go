package handler

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

type Yaml struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := mapYAML(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yamlByte []byte) ([]Yaml, error) {
	var y []Yaml
	err := yaml.Unmarshal(yamlByte, &y)
	return y, err
}

func mapYAML(parsedYaml []Yaml) map[string]string {
	yamlMap := make(map[string]string)
	for _, y := range parsedYaml {
		yamlMap[y.Path] = y.Url
	}
	return yamlMap
}
