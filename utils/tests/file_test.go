package tests

import (
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"testing"
)

func TestReadFileIO(t *testing.T) {
	file, err := utils.CreateTempFile("")
	if err != nil {
		t.Fatalf("Error while creating temporary file: %s", err)
		return
	}

	dir, err := utils.CreateTempDir()
	if err != nil {
		t.Fatalf("Error while creating temporary directory: %s", err)
		return
	}
	defer os.Remove(file.String())

	const fakeFile = models.Path("fooFileName")
	const fileContents = "Hello World!"

	err = utils.WriteStrToFile(dir, fileContents)
	if err == nil {
		t.Fatalf("WriteStrToFile() to directory should have returned an error, but did not")
		return
	}

	err = utils.WriteStrToFile(file, fileContents)
	if err != nil {
		t.Fatalf("Error while writing to temporary file: %s", err)
		return
	}

	_, err = utils.ReadFileToStr(fakeFile)
	if err == nil {
		t.Fatalf("ReadFileToStr() for non-existing file should error, but did not")
	}

	_, err = utils.ReadFileToStr(dir)
	if err == nil {
		t.Fatalf("ReadFileToStr() for existing directory should error, but did not")
	}

	readFileContents, err := utils.ReadFileToStr(file)
	if err != nil {
		t.Fatalf("ReadFileToStr() returned error, but should've not: %s", err)
		return
	}

	if readFileContents != fileContents {
		t.Fatalf("Writing and reading file caused different results, expected: %s, got: %s", fileContents, readFileContents)
	}
}
