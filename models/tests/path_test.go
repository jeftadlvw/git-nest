package tests

import (
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"testing"
)

func TestPathExists(t *testing.T) {
	tempDir, err := utils.CreateTempDir()
	if err != nil {
		t.Errorf("Error creating temporary directory: %s", err)
		return
	}

	tempFile, err := utils.CreateTempFile(tempDir)
	if err != nil {
		t.Errorf("Error creating temporary file: %s", err)
		return
	}

	if !tempDir.Exists() {
		t.Errorf("tempDirectory does not exist: %s", tempFile.String())
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

	os.RemoveAll(tempDir.String())

	if tempDir.Exists() || tempDir.IsDir() || tempDir.IsFile() {
		t.Errorf("tempDirectory should not exist: %s", tempFile.String())
	}

	if tempFile.Exists() || tempFile.IsDir() || tempFile.IsFile() {
		t.Errorf("File should not exist: %s", tempFile.String())
	}
}
