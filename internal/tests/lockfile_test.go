package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"testing"
)

const lockFileName = "hello_lockfile"

func TestAcquireLockFile(t *testing.T) {
	cases := []struct {
		createDir  bool
		createFile bool
		err        bool
	}{
		{false, false, false},
		{true, false, true},
		{true, true, true},
		{false, true, true},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestAcquireLockFile-%d", index+1), func(t *testing.T) {
			var err error
			tempDir := models.Path(t.TempDir())
			lockFilePath := tempDir.SJoin(lockFileName)

			if tc.createDir {
				err = os.MkdirAll(lockFilePath.String(), os.ModePerm)
				if err != nil {
					t.Fatalf("could not create mock directory: %s", err)
				}
			}

			if !tc.createDir && tc.createFile {
				err = utils.WriteStrToFile(lockFilePath, "")
				if err != nil {
					t.Fatalf("could not create mock lock file: %s", err)
				}
			}

			lockFile, err := internal.CreateLockFile(lockFilePath)
			defer lockFile.Release()

			// test for errors
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !tc.err {
				if !lockFilePath.IsFile() {
					t.Fatalf("lockfile was not created")
				}

				lockFileTest, err := internal.CreateLockFile(lockFilePath)
				if err == nil {
					t.Fatalf("second acquiring was successful")
					defer lockFileTest.Release()
				}
			}
		})
	}
}

func TestReleaseLockFile(t *testing.T) {
	tempDir := models.Path(t.TempDir())
	lockFilePath := tempDir.SJoin(lockFileName)

	lockFile, err := internal.CreateLockFile(lockFilePath)
	if err != nil {
		t.Fatalf("acquiring lockfile failed: %s", err)
	}

	err = lockFile.Release()
	// test for errors
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	errSecondRelease := lockFile.Release()
	// test for errors
	if errSecondRelease == nil {
		t.Fatalf("releasing a lockfile a second time should cause error")
	}
}
