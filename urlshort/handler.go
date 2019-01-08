package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(writer, request, url, http.StatusFound)
		}
		fallback.ServeHTTP(writer, request)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(parsedYaml)
	return MapHandler(pathsToUrls, fallback), nil
}

// JSONHandler will parse the provided JSON and then return and http.HandlerFunc
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(jsonBytes)
	if err != nil {
		return nil, err
	}
	pathToUrls := buildMap(parsedJSON)
	return MapHandler(pathToUrls, fallback), nil
}

func parseJSON(jsonBytes []byte) ([]pathURL, error) {
	var pathUrls []pathURL
	return pathUrls, nil
}

func parseYAML(yamlBytes []byte) ([]pathURL, error) {
	var pathUrls []pathURL
	err := yaml.Unmarshal(yamlBytes, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func buildMap(pathUrls []pathURL) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pathURL := range pathUrls {
		pathsToUrls[pathURL.Path] = pathURL.URL
	}
	return pathsToUrls
}

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
