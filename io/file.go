package io

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"os"
)

/*
ReadFileToStr is a wrapper for os.ReadFile that checks for the file's existence before reading a file and returns the file's contents as a string.
*/
func ReadFileToStr(path models.Path) (string, error) {
	if !path.IsFile() {
		return "", fmt.Errorf("%s is not a file", path.String())
	}

	bytes, err := os.ReadFile(path.String())
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}