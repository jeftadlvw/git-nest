package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/actions"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"github.com/jeftadlvw/git-nest/test_env"
	test_env_models "github.com/jeftadlvw/git-nest/test_env/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"testing"
)

func TestAddSubmoduleInContext(t *testing.T) {

	testRepoUrl, err := urls.HttpUrlFromString(test_env.RepoUrl)
	if err != nil {
		t.Fatal(err)
	}

	testFile := "testfileinrepo"
	testDirEmpty := "testdirinrepoempty"
	testDirFull := "testdirinrepofull"

	cases := []struct {
		url                      string
		cloneDir                 string
		ref                      string
		submodules               []models.Submodule
		simulateNoGitInstalled   bool
		simulateDuplicateOrigins bool
		err                      bool
	}{
		{test_env.RepoUrl, "", "", []models.Submodule{}, true, false, true},
		{test_env.RepoUrl, "", "", []models.Submodule{}, false, false, false},
		{test_env.RepoUrl, testFile, "", []models.Submodule{}, false, false, true},
		{test_env.RepoUrl, testFile + "/", "", []models.Submodule{}, false, false, true},
		{test_env.RepoUrl, testDirEmpty, "", []models.Submodule{}, false, false, false},
		{test_env.RepoUrl, testDirEmpty + "/", "", []models.Submodule{}, false, false, false},
		{test_env.RepoUrl, testDirFull, "", []models.Submodule{}, false, false, true},
		{test_env.RepoUrl, testDirFull + "/", "", []models.Submodule{}, false, false, false},
		{test_env.RepoUrl, "../foo", "", []models.Submodule{}, false, false, true},
		{test_env.RepoUrl, "foo/bar/../", "", []models.Submodule{}, false, false, false},
		{test_env.RepoUrl, "foo/bar/..", "", []models.Submodule{}, false, false, false},
		{test_env.RepoUrl, "foo/bar/../..", "", []models.Submodule{}, false, false, true},
		{test_env.RepoUrl, "foo/bar/../../", "", []models.Submodule{}, false, false, false},
		{test_env.RepoUrl, "foo/bar/../../foo", "", []models.Submodule{}, false, false, false},
		{test_env.RepoUrl, "foo", "", []models.Submodule{}, false, false, false},
		{test_env.RepoUrl, "f!oo", "", []models.Submodule{}, false, false, true},
		{test_env.RepoUrl, "f!oo", "", []models.Submodule{}, false, false, true},
		{test_env.RepoUrl, "fo*o", "", []models.Submodule{}, false, false, true},
		{test_env.RepoUrl, "/foo", "", []models.Submodule{}, false, false, true},
		{test_env.RepoUrl, "/../foo", "", []models.Submodule{}, false, false, true},
		{test_env.RepoUrl, "foo", test_env.RepoBranch1, []models.Submodule{}, false, false, false},
		{test_env.RepoUrl, "", "", []models.Submodule{{"example-repository", testRepoUrl, ""}}, false, false, true},
		{test_env.RepoUrl, "", "", []models.Submodule{{"example-repository", testRepoUrl, ""}}, false, true, true},
		{test_env.RepoUrl, "", test_env.RepoBranch1, []models.Submodule{{"foo", testRepoUrl, ""}}, false, false, true},
		{test_env.RepoUrl, "", test_env.RepoBranch1, []models.Submodule{{"foo", testRepoUrl, ""}}, false, true, false},
		{test_env.RepoUrl, "", test_env.RepoBranch1, []models.Submodule{{"foo", testRepoUrl, test_env.RepoBranch1}}, false, true, false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestAddSubmoduleInContext-%d", index+1), func(t *testing.T) {
			t.Parallel()

			testEnv, err := test_env.CreateTestEnvironment(test_env_models.EnvSettings{EmptyGit: true})
			if err != nil {
				t.Fatalf("error creating test environment for case %d: %s", index+1, err)
				return
			}
			defer testEnv.Destroy()

			// create test file and directory
			err = utils.WriteStrToFile(testEnv.Dir.SJoin(testFile), "")
			if err != nil {
				t.Fatalf("error writing test file: %s", err)
			}

			localTestDirEmpty := testEnv.Dir.SJoin(testDirEmpty)
			localTestDirFull := testEnv.Dir.SJoin(testDirFull)
			err = os.Mkdir(localTestDirEmpty.String(), os.ModePerm)
			err = os.Mkdir(localTestDirFull.String(), os.ModePerm)

			err = utils.WriteStrToFile(localTestDirFull.SJoin(testFile), "")
			if err != nil {
				t.Fatalf("error writing test file: %s", err)
			}

			// create context
			context, err := internal.CreateContext(testEnv.Dir)
			if err != nil {
				t.Fatalf("error creating context for case %d: %s", index+1, err)
			}

			context.Config.Submodules = tc.submodules

			if tc.simulateNoGitInstalled {
				context.IsGitInstalled = false
			}

			if tc.simulateDuplicateOrigins {
				context.Config.Config.AllowDuplicateOrigins = true
			}

			// add submodule
			url, err := urls.HttpUrlFromString(tc.url)
			if err != nil {
				t.Errorf("unable to convert url for case %d: %s", index+1, err)
			}

			cloneDir := tc.cloneDir
			err = actions.AddSubmoduleInContext(&context, url, tc.ref, models.Path(cloneDir))

			// check things
			if tc.err && err == nil {
				t.Errorf("AddSubmoduleInContext() for case %d returned no error but expected one", index+1)
			}
			if !tc.err && err != nil {
				t.Errorf("AddSubmoduleInContext() for case %d returned error, but should've not -> %s", index+1, err)
			}
		})
	}
}
