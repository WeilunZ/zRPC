package utils

import (
	"errors"
	"strings"
)

func ParseServicePath(path string) (string, string, error) {
	idx := strings.LastIndex(path, "/")
	if idx == 0 || idx == -1 || !strings.HasPrefix(path, "/") {
		return "", "", errors.New("invalid path")
	}
	return path[1:idx], path[idx+1:], nil
}
