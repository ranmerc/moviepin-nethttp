package utils

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidPath = errors.New("invalid path")
)

func GetIDFromPath(path string) (string, error) {
	matches := regexp.MustCompile(`/movies/([^/]+)/?$`).FindStringSubmatch(path)

	if len(matches) != 2 {
		return "", ErrInvalidPath
	}

	return matches[1], nil
}
