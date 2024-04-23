package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/test_env"
	test_env_models "github.com/jeftadlvw/git-nest/test_env/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"path/filepath"
	"testing"
)

const (
	nonExistingDir = models.Path("foo124TestHelloWorld!")
	nonExistingRef = "foo124TestHelloWorld!"
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
		{test_env.RepoUrl, "foobartest", "", false, true},
		{test_env.RepoUrl, "", "./foo", true, true},
		{test_env.RepoUrl, "", "foo/bar", true, true},
		{test_env.RepoUrl, "", "..", true, true},
		{test_env.RepoUrl, "", ".", true, false},
		{test_env.RepoUrl, "", "foo", true, false},
		{test_env.RepoUrl, "", "", true, false},
		{test_env.RepoUrlNoSuffix, "", "", true, false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestCloneGitRepository-%d", index), func(t *testing.T) {
			t.Parallel()
			clonePath := tc.path
			if tc.tempPath {
				tempDir, err := utils.CreateTempDir()
				if err != nil {
					t.Fatalf("error creating temporary directory for case %d: %s", index+1, err)
					return
				}
				defer os.RemoveAll(clonePath.String())
				clonePath = tempDir
			}

			err := utils.CloneGitRepository(tc.url, clonePath, tc.cloneDirName)
			if tc.err && err == nil {
				t.Errorf("CloneGitRepository() for case %d returned no error but expected one", index+1)
			}
			if !tc.err && err != nil {
				t.Errorf("CloneGitRepository() for case %d returned error, but should've not -> %s", index+1, err)
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
		{test_env.RepoBranchDefault, "", false, true},
		{test_env.RepoBranchDefault, nonExistingDir, false, true},
		{"", "", true, true},
		{"    ", "", true, true},
		{nonExistingRef, "", true, true},
		{test_env.RepoBranchDefault, "", true, false},
		{"    \n\t" + test_env.RepoBranchDefault + "   ", "", true, false},
		{test_env.RepoBranch1, "", true, false},
		{test_env.RepoCommit, "", true, false},
		{test_env.RepoCommitLong, "", true, false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestChangeGitHead-%d", index+1), func(t *testing.T) {
			t.Parallel()
			repoDir := tc.dir

			if tc.useExampleRepo {

				testEnv, err := test_env.CreateTestEnvironment(test_env_models.EnvSettings{Origin: test_env.RepoUrl, CloneDir: "temp"})
				if err != nil {
					t.Fatalf("error creating test environment for case %d: %s", index+1, err)
					return
				}
				defer testEnv.Destroy()

				repoDir = testEnv.Dir.SJoin("temp")
			}

			err := utils.ChangeGitHead(repoDir, tc.ref)
			if tc.err && err == nil {
				t.Errorf("CloneGitRepository() for case %d returned no error but expected one", index+1)
			}
			if !tc.err && err != nil {
				t.Errorf("CloneGitRepository() for case %d returned error, but should've not -> %s", index+1, err)
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
				var tempDir models.Path
				testEnv, err := test_env.CreateTestEnvironment(test_env_models.EnvSettings{Origin: test_env.RepoUrl, CloneDir: tc.cloneDir})
				if err != nil {
					t.Fatalf("error creating test environment for case %d: %s", index+1, err)
					return
				}
				defer testEnv.Destroy()

				tempDir = testEnv.Dir

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
				t.Errorf("GetGitRootDirectory() for case %d returned no error but expected one", index+1)
				return
			}
			if !tc.err && err != nil {
				t.Errorf("GetGitRootDirectory() for case %d returned error, but should've not -> %s", index+1, err)
				return
			}
			if remoteUrl != expectDir {
				t.Errorf("GetGitRootDirectory() for case %d returned unexpected remote: >%s<, expected >%s<,", index+1, remoteUrl, expectDir)
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
		{"", true, test_env.RepoUrl, test_env.RepoUrl, false},
		{"", true, test_env.RepoUrlNoSuffix, test_env.RepoUrlNoSuffix, false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestGetGitRemoteUrl-%d", index+1), func(t *testing.T) {
			t.Parallel()
			repoDir := tc.dir

			if tc.useTempDir {
				testEnv, err := test_env.CreateTestEnvironment(test_env_models.EnvSettings{Origin: tc.cloneTempRepo, CloneDir: "temp"})
				if err != nil {
					t.Fatalf("error creating test environment for case %d: %s", index+1, err)
					return
				}
				defer testEnv.Destroy()

				repoDir = testEnv.Dir.SJoin("temp")
			}

			remoteUrl, err := utils.GetGitRemoteUrl(repoDir)
			if tc.err && err == nil {
				t.Errorf("GetGitRemoteUrl() for case %d returned no error but expected one", index+1)
				return
			}
			if !tc.err && err != nil {
				t.Errorf("GetGitRemoteUrl() for case %d returned error, but should've not -> %s", index+1, err)
				return
			}
			if remoteUrl != tc.expected {
				t.Errorf("GetGitRemoteUrl() for case %d returned unexpected remote: >%s<, expected >%s<,", index+1, remoteUrl, tc.expected)
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
		{"", true, "", test_env.RepoBranchDefaultRefLong, test_env.RepoBranchDefault, false},
		{"", true, test_env.RepoBranch1, test_env.RepoBranch1RefLong, test_env.RepoBranch1, false},
		{"", true, test_env.RepoCommit, test_env.RepoCommitLong, "", false},
		{"", true, test_env.RepoCommitLong, test_env.RepoCommitLong, "", false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("GetGitFetchHead-%d", index+1), func(t *testing.T) {
			t.Parallel()
			repoDir := tc.dir

			if tc.useExampleRepo {
				envSettings := test_env_models.EnvSettings{Origin: test_env.RepoUrl, CloneDir: "temp"}
				if tc.checkoutBeforeTest != "" {
					envSettings.Ref = tc.checkoutBeforeTest
				}

				testEnv, err := test_env.CreateTestEnvironment(envSettings)
				if err != nil {
					t.Fatalf("error creating test environment for case %d: %s", index+1, err)
					return
				}
				defer testEnv.Destroy()

				repoDir = testEnv.Dir.SJoin("temp")
			}

			headRef, headRefAbbrev, err := utils.GetGitFetchHead(repoDir)
			if tc.err && err == nil {
				t.Errorf("GetGitFetchHead() for case %d returned no error but expected one", index+1)
				return
			}
			if !tc.err && err != nil {
				t.Errorf("GetGitFetchHead() for case %d returned error, but should've not -> %s", index+1, err)
				return
			}
			if headRef != tc.expectedRef {
				t.Errorf("GetGitFetchHead() for case %d returned unexpected ref: >%s<, expected >%s<,", index+1, headRef, tc.expectedRef)
			}
			if headRefAbbrev != tc.expectedAbbrev {
				t.Errorf("GetGitFetchHead() for case %d returned unexpected abbreviation: >%s<, expected >%s<,", index+1, headRefAbbrev, tc.expectedAbbrev)
			}
		})
	}
}
