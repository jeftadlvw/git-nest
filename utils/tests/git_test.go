package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"path/filepath"
	"testing"
)

const (
	testRepoUrl                  = "https://github.com/jeftadlvw/example-repository.git"
	testRepoUrlNoSuffix          = "https://github.com/jeftadlvw/example-repository"
	testRepoBranchDefault        = "main"
	testRepoBranchDefaultRefLong = "0ab2d7ab4e49272a3f8955fbc79d34895d49bb31"
	testRepoBranch1              = "foo"
	testRepoBranch1RefLong       = "65a3ef29441587285eb2ceb42bee5f4a2a534110"
	testRepoCommit               = "d9c591c" // on branch 2
	testRepoCommitLong           = "d9c591ca90aa1cedda54d1d6ebb45be3b52e5d6e"
	nonExistingDir               = models.Path("foo124TestHelloWorld!")
	nonExistingRef               = "foo124TestHelloWorld!"
)

func TestCloneGitRepository(t *testing.T) {
	t.Parallel()

	cases := []struct {
		url          string
		path         models.Path
		cloneDirName string
		tempPath     bool
		err          bool
	}{
		{"", "", "", false, true},
		{"iDoNotExist", "", "", false, true},
		{testRepoUrl, "foobartest", "", false, true},
		{testRepoUrl, "", "./foo", true, true},
		{testRepoUrl, "", "foo/bar", true, true},
		{testRepoUrl, "", ".", true, false},
		{testRepoUrl, "", "foo", true, false},
		{testRepoUrl, "", "", true, false},
		{testRepoUrlNoSuffix, "", "", true, false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestCloneGitRepository-%d", index), func(t *testing.T) {
			t.Parallel()
			clonePath := tc.path
			if tc.tempPath {
				tempDir, err := utils.CreateTempDir()
				if err != nil {
					t.Fatalf("error creating temporary directory for %v: %s", tc, err)
					return
				}
				defer os.RemoveAll(clonePath.String())
				clonePath = tempDir
			}

			err := utils.CloneGitRepository(tc.url, clonePath, tc.cloneDirName)
			if tc.err && err == nil {
				t.Errorf("CloneGitRepository() for %v returned no error but expected one", tc)
			}
			if !tc.err && err != nil {
				t.Errorf("CloneGitRepository() for %v returned error, but should've not -> %s", tc, err)
			}
		})
	}

}

func TestChangeGitHead(t *testing.T) {
	t.Parallel()

	cases := []struct {
		ref            string
		dir            models.Path
		useExampleRepo bool
		err            bool
	}{
		{"", "", false, true},
		{"     ", "", false, true},
		{testRepoBranchDefault, nonExistingDir, false, true},
		{"", "", true, true},
		{"    ", "", true, true},
		{nonExistingRef, "", true, true},
		{testRepoBranchDefault, "", true, false},
		{"    \n\t" + testRepoBranchDefault + "   ", "", true, false},
		{testRepoBranch1, "", true, false},
		{testRepoCommit, "", true, false},
		{testRepoCommitLong, "", true, false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestChangeGitHead-%d", index+1), func(t *testing.T) {
			t.Parallel()
			repoDir := tc.dir

			if tc.useExampleRepo {
				tempDir, err := utils.CreateTempDir()
				if err != nil {
					t.Fatalf("error creating temporary directory for %v: %s", tc, err)
					return
				}
				defer os.RemoveAll(tempDir.String())

				err = utils.CloneGitRepository(testRepoUrl, tempDir, "temp")
				if err != nil {
					t.Fatalf("unable to clone git repository for further tests: %s", err)
				}

				repoDir = tempDir.SJoin("temp")

			}

			err := utils.ChangeGitHead(repoDir, tc.ref)
			if tc.err && err == nil {
				t.Errorf("CloneGitRepository() for %v returned no error but expected one", tc)
			}
			if !tc.err && err != nil {
				t.Errorf("CloneGitRepository() for %v returned error, but should've not -> %s", tc, err)
			}
		})
	}
}

func TestGetGitRootDirectory(t *testing.T) {
	t.Parallel()

	testTempDir, err := utils.CreateTempDir()
	if err != nil {
		t.Fatalf("error creating temporary directory for %s", err)
		return
	}
	defer os.RemoveAll(testTempDir.String())

	cases := []struct {
		dir        models.Path
		useTempDir bool
		cloneDir   string
		expected   string
		err        bool
	}{
		{"", false, "", "", true},
		{"   ", false, "", "", true},
		{nonExistingDir, false, "", "", true},
		{testTempDir, false, "", "", true},
		{"", true, "", "example-repository", false},
		{"", true, ".", ".", false},
		{"", true, "temprepo", "temprepo", false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestGetGitRootDirectory-%d", index+1), func(t *testing.T) {
			t.Parallel()

			repoDir := tc.dir
			expectDir := tc.expected

			if tc.useTempDir {
				tempDir, err := utils.CreateTempDir()
				if err != nil {
					t.Fatalf("error creating temporary directory for %v: %s", tc, err)
					return
				}
				defer os.RemoveAll(tempDir.String())

				err = utils.CloneGitRepository(testRepoUrl, tempDir, tc.cloneDir)
				if err != nil {
					t.Fatalf("unable to clone git repository for further tests: %s", err)
					return
				}

				if tc.cloneDir == "" {
					repoDir = tempDir.SJoin("example-repository")
				} else {
					repoDir = tempDir.SJoin(tc.cloneDir)
					repoDir = repoDir.Clean()
				}

				expectDirTemp := tempDir.SJoin(tc.expected)
				expectDir = expectDirTemp.String()
				expectDir, err = filepath.EvalSymlinks(expectDir)
				if err != nil {
					t.Fatalf("unable to resolve expected directory for further tests: %s", err)
					return
				}
			}

			remoteUrl, err := utils.GetGitRootDirectory(repoDir)
			if tc.err && err == nil {
				t.Errorf("GetGitRootDirectory() for %v returned no error but expected one", tc)
				return
			}
			if !tc.err && err != nil {
				t.Errorf("GetGitRootDirectory() for %v returned error, but should've not -> %s", tc, err)
				return
			}
			if remoteUrl != expectDir {
				t.Errorf("GetGitRootDirectory() for %v returned unexpected remote: >%s<, expected >%s<,", tc, remoteUrl, expectDir)
			}
		})
	}
}

func TestGetGitRemoteUrl(t *testing.T) {
	t.Parallel()

	cases := []struct {
		dir           models.Path
		useTempDir    bool
		cloneTempRepo string
		expected      string
		err           bool
	}{
		{"", false, "", "", true},
		{"   ", false, "", "", true},
		{nonExistingDir, false, "", "", true},
		{"", true, testRepoUrl, testRepoUrl, false},
		{"", true, testRepoUrlNoSuffix, testRepoUrlNoSuffix, false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestGetGitRemoteUrl-%d", index+1), func(t *testing.T) {
			t.Parallel()
			repoDir := tc.dir

			if tc.useTempDir {
				tempDir, err := utils.CreateTempDir()
				if err != nil {
					t.Fatalf("error creating temporary directory for %v: %s", tc, err)
					return
				}
				defer os.RemoveAll(tempDir.String())

				err = utils.CloneGitRepository(tc.cloneTempRepo, tempDir, "temp")
				if err != nil {
					t.Fatalf("unable to clone git repository for further tests: %s", err)
					return
				}

				repoDir = tempDir.SJoin("temp")
			}

			remoteUrl, err := utils.GetGitRemoteUrl(repoDir)
			if tc.err && err == nil {
				t.Errorf("GetGitRemoteUrl() for %v returned no error but expected one", tc)
				return
			}
			if !tc.err && err != nil {
				t.Errorf("GetGitRemoteUrl() for %v returned error, but should've not -> %s", tc, err)
				return
			}
			if remoteUrl != tc.expected {
				t.Errorf("GetGitRemoteUrl() for %v returned unexpected remote: >%s<, expected >%s<,", tc, remoteUrl, tc.expected)
			}
		})
	}
}

func TestGetGitFetchHead(t *testing.T) {
	t.Parallel()

	cases := []struct {
		dir                models.Path
		useExampleRepo     bool
		checkoutBeforeTest string
		expectedRef        string
		expectedAbbrev     string
		err                bool
	}{
		{"", false, "", "", "", true},
		{"   ", false, "", "", "", true},
		{nonExistingDir, false, "", "", "", true},
		{"", true, "", testRepoBranchDefaultRefLong, testRepoBranchDefault, false},
		{"", true, testRepoBranch1, testRepoBranch1RefLong, testRepoBranch1, false},
		{"", true, testRepoCommit, testRepoCommitLong, "", false},
		{"", true, testRepoCommitLong, testRepoCommitLong, "", false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("GetGitFetchHead-%d", index+1), func(t *testing.T) {
			t.Parallel()
			repoDir := tc.dir

			if tc.useExampleRepo {
				tempDir, err := utils.CreateTempDir()
				if err != nil {
					t.Fatalf("error creating temporary directory for %v: %s", tc, err)
					return
				}
				defer os.RemoveAll(tempDir.String())

				err = utils.CloneGitRepository(testRepoUrl, tempDir, "temp")
				if err != nil {
					t.Fatalf("unable to clone git repository for further tests: %s", err)
					return
				}

				repoDir = tempDir.SJoin("temp")
			}

			if tc.checkoutBeforeTest != "" {
				err := utils.ChangeGitHead(repoDir, tc.checkoutBeforeTest)
				if err != nil {
					t.Fatalf("unable to checkout before test: %s", err)
					return
				}
			}

			headRef, headRefAbbrev, err := utils.GetGitFetchHead(repoDir)
			if tc.err && err == nil {
				t.Errorf("GetGitFetchHead() for %v returned no error but expected one", tc)
				return
			}
			if !tc.err && err != nil {
				t.Errorf("GetGitFetchHead() for %v returned error, but should've not -> %s", tc, err)
				return
			}
			if headRef != tc.expectedRef {
				t.Errorf("GetGitFetchHead() for %v returned unexpected ref: >%s<, expected >%s<,", tc, headRef, tc.expectedRef)
			}
			if headRefAbbrev != tc.expectedAbbrev {
				t.Errorf("GetGitFetchHead() for %v returned unexpected abbreviation: >%s<, expected >%s<,", tc, headRefAbbrev, tc.expectedAbbrev)
			}
		})
	}
}
