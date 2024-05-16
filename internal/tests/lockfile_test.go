package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"testing"
)

func TestAcquireLockFile(t *testing.T) {

	cases := []struct {
		createFile   bool
		createDir    bool
		expectedFile string
		err          bool
	}{
		{false, false, constants.LockFileName, false},
		{true, false, constants.LockFileName, true},
		{true, true, constants.LockFileName, true},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestAcquireLockFile-%d", index+1), func(t *testing.T) {
			var err error
			tempDir := models.Path(t.TempDir())

			if tc.createDir {
				err = os.MkdirAll(tempDir.String()+"/"+tc.expectedFile, os.ModePerm)
				if err != nil {
					t.Fatalf("could not create mock directory: %s", err)
				}
			}

			if !tc.createDir && tc.createFile {
				err = utils.WriteStrToFile(tempDir.SJoin(tc.expectedFile), "")
				if err != nil {
					t.Fatalf("could not create mock lock file: %s", err)
				}
			}

			lockFile, err := internal.AcquireLockFile(tempDir)
			defer internal.ReleaseLockFile(lockFile)

			// test for errors
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !tc.err {
				expectedFile := tempDir.SJoin(tc.expectedFile)
				if !expectedFile.IsFile() {
					t.Fatalf("lockfile was not created")
				}

				lockFileTest, err := internal.AcquireLockFile(tempDir)
				if err == nil {
					t.Fatalf("second acquiring was successful")
					defer internal.ReleaseLockFile(lockFileTest)
				}
			}
		})
	}
}

func TestReleaseLockFile(t *testing.T) {

	tempDir := models.Path(t.TempDir())

	lockFile, err := internal.AcquireLockFile(tempDir)
	if err != nil {
		t.Fatalf("acquiring lockfile failed: %s", err)
	}

	err = internal.ReleaseLockFile(lockFile)
	// test for errors
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	errSecondRelease := internal.ReleaseLockFile(lockFile)
	// test for errors
	if errSecondRelease == nil {
		t.Fatalf("releasing a lockfile a second time should cause error")
	}

}
