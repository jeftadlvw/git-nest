package utils

import (
	"github.com/jeftadlvw/git-nest/models"
	"os"
)

/*
CreateTempDir creates a temporary directory using os.MkdirTemp.
*/
func CreateTempDir() (models.Path, error) {
	dir, err := os.MkdirTemp("", "")
	return models.Path(dir), err
}

/*
CreateTempFile creates a temporary file using os.CreateTemp.
*/
func CreateTempFile(dir models.Path) (models.Path, error) {
	file, err := os.CreateTemp(dir.String(), "")
	if err != nil {
		return "", err
	}
	return models.Path(file.Name()), nil
}
