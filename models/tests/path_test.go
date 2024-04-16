package tests

import (
	"github.com/jeftadlvw/git-nest/models"
	"os"
	"testing"
)

func TestPathExists(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Errorf("Error creating temporary directory: %s", err)
		return
	}
	tempDir := models.Path(dir)
	defer os.RemoveAll(tempDir.String())

	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("Error creating temporary file: %s", err)
		return
	}
	tempFile := models.Path(file.Name())
	defer os.Remove(tempFile.String())

	if !tempDir.Exists() {
		t.Errorf("Directory does not exist: %s", tempFile.String())
	}

	if !tempDir.IsDir() {
		t.Errorf("Directory is not a directory: %s", tempDir.String())
	}

	if tempDir.IsFile() {
		t.Errorf("Path is file but should be dir: %s", tempDir.String())
	}

	if !tempFile.Exists() {
		t.Errorf("File does not exist: %s", tempFile.String())
	}

	if !tempFile.IsFile() {
		t.Errorf("File is not a file: %s", tempFile.String())
	}

	if tempFile.IsDir() {
		t.Errorf("Path is file but should be dir: %s", tempFile.String())
	}
}
