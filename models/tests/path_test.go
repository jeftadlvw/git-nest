package tests

import (
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"testing"
)

func TestPathExists(t *testing.T) {
	tempDir, err := utils.CreateTempDir()
	if err != nil {
		t.Fatalf("error creating temporary directory: %s", err)
		return
	}

	tempFile, err := utils.CreateTempFile(tempDir)
	if err != nil {
		t.Fatalf("Error creating temporary file: %s", err)
		return
	}

	if !tempDir.Exists() {
		t.Fatalf("tempDirectory does not exist: %s", tempFile.String())
	}

	if !tempDir.IsDir() {
		t.Fatalf("directory is not a directory: %s", tempDir.String())
	}

	if tempDir.IsFile() {
		t.Fatalf("PathS is file but should be dir: %s", tempDir.String())
	}

	if !tempFile.Exists() {
		t.Fatalf("file does not exist: %s", tempFile.String())
	}

	if !tempFile.IsFile() {
		t.Fatalf("file is not a file: %s", tempFile.String())
	}

	if tempFile.IsDir() {
		t.Fatalf("PathS is file but should be dir: %s", tempFile.String())
	}

	_ = os.RemoveAll(tempDir.String())

	if tempDir.Exists() || tempDir.IsDir() || tempDir.IsFile() {
		t.Fatalf("tempDirectory should not exist: %s", tempFile.String())
	}

	if tempFile.Exists() || tempFile.IsDir() || tempFile.IsFile() {
		t.Fatalf("file should not exist: %s", tempFile.String())
	}
}
