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
		createFile              bool
		createDir               bool
		isGitRepositoryOverride bool
		expectedFile            string
		success                 bool
		err                     bool
	}{
		{false, false, false, constants.LockFileName, true, false},
		{false, false, true, constants.LockFileNameGitRepo, true, false},
		{true, false, false, constants.LockFileName, false, false},
		{true, false, true, constants.LockFileNameGitRepo, false, false},
		{true, true, false, constants.LockFileName, false, true},
		{true, true, true, constants.LockFileNameGitRepo, false, true},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestAcquireLockFile-%d", index+1), func(t *testing.T) {
			tempDir := models.Path(t.TempDir())
			context, err := internal.CreateContext(tempDir)
			if err != nil {
				t.Fatalf("could not create temporary context: %s", err)
			}

			context.IsGitRepository = tc.isGitRepositoryOverride
			if tc.isGitRepositoryOverride {
				err = os.MkdirAll(context.ProjectRoot.String()+"/.git", os.ModePerm)
				if err != nil {
					t.Fatalf("could not create mock .git directory: %s", err)
				}
			}

			if tc.createDir {
				err = os.MkdirAll(context.ProjectRoot.String()+"/"+tc.expectedFile, os.ModePerm)
				if err != nil {
					t.Fatalf("could not create mock directory: %s", err)
				}
			}

			if !tc.createDir && tc.createFile {
				err = utils.WriteStrToFile(context.ProjectRoot.SJoin(tc.expectedFile), "")
				if err != nil {
					t.Fatalf("could not create mock lock file: %s", err)
				}
			}

			success, err := internal.AcquireLockFile(context)
			// test for errors
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if tc.success && !success {
				t.Fatalf("expected success, but wasn't")
			}
			if !tc.success && success {
				t.Fatalf("expected failure, but wasn't")
			}
		})
	}
}

func TestReleaseLockFile(t *testing.T) {

	cases := []struct {
		createFile              bool
		createDir               bool
		isGitRepositoryOverride bool
		expectedFile            string
		err                     bool
	}{
		{false, false, false, constants.LockFileName, false},
		{false, false, true, constants.LockFileNameGitRepo, false},
		{true, false, false, constants.LockFileName, false},
		{true, false, true, constants.LockFileNameGitRepo, false},
		{true, true, false, constants.LockFileName, true},
		{true, true, true, constants.LockFileNameGitRepo, true},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("ReleaseLockFile-%d", index+1), func(t *testing.T) {
			tempDir := models.Path(t.TempDir())
			context, err := internal.CreateContext(tempDir)
			if err != nil {
				t.Fatalf("could not create temporary context: %s", err)
			}

			context.IsGitRepository = tc.isGitRepositoryOverride
			if tc.isGitRepositoryOverride {
				err = os.MkdirAll(context.ProjectRoot.String()+"/.git", os.ModePerm)
				if err != nil {
					t.Fatalf("could not create mock .git directory: %s", err)
				}
			}

			if tc.createDir {
				err = os.MkdirAll(context.ProjectRoot.String()+"/"+tc.expectedFile, os.ModePerm)
				if err != nil {
					t.Fatalf("could not create mock directory: %s", err)
				}
			}

			if !tc.createDir && tc.createFile {
				err = utils.WriteStrToFile(context.ProjectRoot.SJoin(tc.expectedFile), "")
				if err != nil {
					t.Fatalf("could not create mock lock file: %s", err)
				}
			}

			err = internal.ReleaseLockFile(context)
			// test for errors
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}
