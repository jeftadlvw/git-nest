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
					t.Fatalf("error creating temporary directory: %s", err)
					return
				}
				defer os.RemoveAll(clonePath.String())
				clonePath = tempDir
			}

			err := utils.CloneGitRepository(tc.url, clonePath, tc.cloneDirName, nil)
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}

}

func TestGitCheckout(t *testing.T) {
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
		t.Run(fmt.Sprintf("TestGitCheckout-%d", index+1), func(t *testing.T) {
			t.Parallel()
			repoDir := tc.dir

			if tc.useExampleRepo {
				tempDir := models.Path(t.TempDir())
				err := test_env.CreateTestEnvironment(tempDir, test_env_models.EnvSettings{Origin: test_env.RepoUrl, CloneDir: "temp"})
				if err != nil {
					t.Fatalf("error creating test environment: %s", err)
					return
				}

				repoDir = tempDir.SJoin("temp")
			}

			err := utils.GitCheckout(repoDir, tc.ref)
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}

func TestGetGitRootDirectory(t *testing.T) {
	t.Parallel()

	testTempDir, err := utils.CreateTempDir()
	if err != nil {
		t.Fatalf("error creating temporary directory: %s", err)
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
				tempDir := models.Path(t.TempDir())
				err := test_env.CreateTestEnvironment(tempDir, test_env_models.EnvSettings{Origin: test_env.RepoUrl, CloneDir: tc.cloneDir})
				if err != nil {
					t.Fatalf("error creating test environment: %s", err)
					return
				}

				tempDir = tempDir

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

			rootDir, err := utils.GetGitRootDirectory(repoDir)
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
				return
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
				return
			}
			if rootDir != expectDir {
				t.Fatalf("unexpected remote: >%s<, expected >%s<,", rootDir, expectDir)
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
				tempDir := models.Path(t.TempDir())
				err := test_env.CreateTestEnvironment(tempDir, test_env_models.EnvSettings{Origin: tc.cloneTempRepo, CloneDir: "temp"})
				if err != nil {
					t.Fatalf("error creating test environment: %s", err)
				}

				repoDir = tempDir.SJoin("temp")
			}

			remoteUrl, err := utils.GetGitRemoteUrl(repoDir)
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !tc.err && remoteUrl != tc.expected {
				t.Fatalf("unexpected remote: >%s<, expected >%s<,", remoteUrl, tc.expected)
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

				tempDir := models.Path(t.TempDir())
				err := test_env.CreateTestEnvironment(tempDir, envSettings)
				if err != nil {
					t.Fatalf("error creating test environment: %s", err)
				}

				repoDir = tempDir.SJoin("temp")
			}

			headRef, headRefAbbrev, err := utils.GetGitFetchHead(repoDir)
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if headRef != tc.expectedRef {
				t.Fatalf("unexpected ref: >%s<, expected >%s<,", headRef, tc.expectedRef)
			}
			if headRefAbbrev != tc.expectedAbbrev {
				t.Fatalf("unexpected abbreviation: >%s<, expected >%s<,", headRefAbbrev, tc.expectedAbbrev)
			}
		})
	}
}

func TestGetGitHasUntrackedChanges(t *testing.T) {
	t.Parallel()

	const testFileName = "test.txt"

	cases := []struct {
		dir            models.Path
		useExampleRepo bool
		createFile     bool
		commitFile     bool
		expected       bool
		err            bool
	}{
		{"", false, false, false, true, true},
		{"   ", false, false, false, true, true},
		{nonExistingDir, false, false, false, true, true},
		{"", true, false, false, false, false},
		{"", true, true, false, true, false},
		{"", true, true, true, false, false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("GetGitHasUntrackedChanges-%d", index+1), func(t *testing.T) {
			t.Parallel()
			repoDir := tc.dir

			if tc.useExampleRepo {
				tempDir := models.Path(t.TempDir())
				envSettings := test_env_models.EnvSettings{Origin: test_env.RepoUrl, CloneDir: "temp"}
				err := test_env.CreateTestEnvironment(tempDir, envSettings)
				if err != nil {
					t.Fatalf("error creating test environment: %s", err)
				}
				repoDir = tempDir.SJoin("temp")

				if tc.createFile {
					err = utils.WriteStrToFile(repoDir.SJoin(testFileName), "")
					if err != nil {
						t.Fatalf("error writing test file: %s", err)
					}

					if tc.commitFile {
						out, err := utils.RunCommandCombinedOutput(repoDir, "git", "add", testFileName)
						if err != nil {
							t.Fatalf("error git-adding test file: %s; %s", err, out)
						}

						out, err = utils.RunCommandCombinedOutput(repoDir, "git", "commit", "-m", "\"added test.txt\"")
						if err != nil {
							t.Fatalf("error git-commiting test file: %s; %s", err, out)
						}
					}
				}
			}

			hasUntrackedChanges, err := utils.GetGitHasUntrackedChanges(repoDir)
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if hasUntrackedChanges != tc.expected {
				t.Fatalf("GetGitHasUntrackedChanges() returned %t (should be %t)", hasUntrackedChanges, tc.expected)
			}
		})
	}
}

func TestGetGitHasUnpublishedChanges(t *testing.T) {
	t.Parallel()

	const testFileName = "test.txt"

	cases := []struct {
		dir            models.Path
		useExampleRepo bool
		createFile     bool
		commitFile     bool
		expected       bool
		err            bool
	}{
		{"", false, false, false, true, true},
		{"   ", false, false, false, true, true},
		{nonExistingDir, false, false, false, true, true},
		{"", true, false, false, false, false},
		{"", true, true, false, false, false},
		{"", true, true, true, true, false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("GetGitHasUnpublishedChanges-%d", index+1), func(t *testing.T) {
			t.Parallel()
			repoDir := tc.dir

			if tc.useExampleRepo {
				tempDir := models.Path(t.TempDir())
				envSettings := test_env_models.EnvSettings{Origin: test_env.RepoUrl, CloneDir: "temp"}
				err := test_env.CreateTestEnvironment(tempDir, envSettings)
				if err != nil {
					t.Fatalf("error creating test environment: %s", err)
				}
				repoDir = tempDir.SJoin("temp")

				if tc.createFile {
					err = utils.WriteStrToFile(repoDir.SJoin(testFileName), "")
					if err != nil {
						t.Fatalf("error writing test file: %s", err)
					}

					if tc.commitFile {
						out, err := utils.RunCommandCombinedOutput(repoDir, "git", "add", testFileName)
						if err != nil {
							t.Fatalf("error git-adding test file: %s; %s", err, out)
						}

						out, err = utils.RunCommandCombinedOutput(repoDir, "git", "commit", "-m", "\"added test.txt\"")
						if err != nil {
							t.Fatalf("error git-commiting test file: %s; %s", err, out)
						}
					}
				}
			}

			hasUnpublishedChanges, err := utils.GetGitHasUnpublishedChanges(repoDir)
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if hasUnpublishedChanges != tc.expected {
				t.Fatalf("returned %t, but should be %t", hasUnpublishedChanges, tc.expected)
			}
		})
	}
}
