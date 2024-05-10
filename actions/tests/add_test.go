package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/actions"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/migrations"
	mcontext "github.com/jeftadlvw/git-nest/migrations/context"
	"github.com/jeftadlvw/git-nest/migrations/git"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"github.com/jeftadlvw/git-nest/test_env"
	test_env_models "github.com/jeftadlvw/git-nest/test_env/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"reflect"
	"testing"
)

func TestAddSubmoduleInContext(t *testing.T) {

	testRepoUrl, lerr := urls.HttpUrlFromString(test_env.RepoUrl)
	if lerr != nil {
		t.Fatal(lerr)
	}

	testFile := "testfileinrepo"
	testDirEmpty := "testdirinrepoempty"
	testDirFull := "testdirinrepofull"

	expectedMigrations := []interfaces.Migration{mcontext.AppendSubmodule{}, git.Clone{}}
	expectedMigrationsRef := []interfaces.Migration{mcontext.AppendSubmodule{}, git.Clone{}, git.Checkout{}}

	cases := []struct {
		url                      string
		cloneDir                 string
		ref                      string
		submodules               []models.Submodule
		simulateNoGitInstalled   bool
		simulateDuplicateOrigins bool
		expectedMigrations       []interfaces.Migration
		err                      bool
	}{
		{test_env.RepoUrl, "", "", []models.Submodule{}, true, false, nil, true},
		{test_env.RepoUrl, "", "", []models.Submodule{}, false, false, expectedMigrations, false},
		{test_env.RepoUrl, testFile, "", []models.Submodule{}, false, false, nil, true},
		{test_env.RepoUrl, testFile + "/", "", []models.Submodule{}, false, false, nil, true},
		{test_env.RepoUrl, testDirEmpty, "", []models.Submodule{}, false, false, expectedMigrations, false},
		{test_env.RepoUrl, testDirEmpty + "/", "", []models.Submodule{}, false, false, expectedMigrations, false},
		{test_env.RepoUrl, testDirFull, "", []models.Submodule{}, false, false, nil, true},
		{test_env.RepoUrl, testDirFull + "/", "", []models.Submodule{}, false, false, expectedMigrations, false},
		{test_env.RepoUrl, "../foo", "", []models.Submodule{}, false, false, nil, true},
		{test_env.RepoUrl, "foo/bar/../", "", []models.Submodule{}, false, false, expectedMigrations, false},
		{test_env.RepoUrl, "foo/bar/..", "", []models.Submodule{}, false, false, expectedMigrations, false},
		{test_env.RepoUrl, "foo/bar/../..", "", []models.Submodule{}, false, false, nil, true},
		{test_env.RepoUrl, "foo/bar/../../", "", []models.Submodule{}, false, false, expectedMigrations, false},
		{test_env.RepoUrl, "foo/bar/../../foo", "", []models.Submodule{}, false, false, expectedMigrations, false},
		{test_env.RepoUrl, "foo", "", []models.Submodule{}, false, false, expectedMigrations, false},
		{test_env.RepoUrl, "f!oo", "", []models.Submodule{}, false, false, nil, true},
		{test_env.RepoUrl, "f!oo", "", []models.Submodule{}, false, false, nil, true},
		{test_env.RepoUrl, "fo*o", "", []models.Submodule{}, false, false, nil, true},
		{test_env.RepoUrl, "/foo", "", []models.Submodule{}, false, false, nil, true},
		{test_env.RepoUrl, "/../foo", "", []models.Submodule{}, false, false, nil, true},
		{test_env.RepoUrl, "foo", test_env.RepoBranch1, []models.Submodule{}, false, false, expectedMigrationsRef, false},
		{test_env.RepoUrl, "", "", []models.Submodule{{"example-repository", testRepoUrl, ""}}, false, false, nil, true},
		{test_env.RepoUrl, "", "", []models.Submodule{{"example-repository", testRepoUrl, ""}}, false, true, nil, true},
		{test_env.RepoUrl, "", test_env.RepoBranch1, []models.Submodule{{"foo", testRepoUrl, ""}}, false, false, nil, true},
		{test_env.RepoUrl, "", test_env.RepoBranch1, []models.Submodule{{"foo", testRepoUrl, ""}}, false, true, expectedMigrationsRef, false},
		{test_env.RepoUrl, "", test_env.RepoBranch1, []models.Submodule{{"foo", testRepoUrl, test_env.RepoBranch1}}, false, true, expectedMigrationsRef, false},
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
			migrationArr, err := actions.AddSubmoduleInContext(&context, url, tc.ref, models.Path(cloneDir))

			// check migration array
			if !tc.err && len(tc.expectedMigrations) != len(migrationArr) {
				t.Fatalf("AddSubmoduleInContext() for case %d returned unequal amounts of migrations: expected %d, got %d", index+1, len(tc.expectedMigrations), len(migrationArr))
			}
			if !tc.err {
				for mindex, migration := range migrationArr {
					if reflect.TypeOf(migration) != reflect.TypeOf(tc.expectedMigrations[mindex]) {
						t.Fatalf("AddSubmoduleInContext() for case %d had unexpected migration at index %d: %s != %s", index+1, mindex, reflect.TypeOf(migration), reflect.TypeOf(tc.expectedMigrations[mindex]))
					}
				}
			}

			// run migrations if action function call did not return an error
			if err == nil {
				err = migrations.RunMigrations(migrationArr...)
			}

			// test for errors
			if tc.err && err == nil {
				t.Fatalf("AddSubmoduleInContext() for case %d returned no error but expected one", index+1)
			}
			if !tc.err && err != nil {
				t.Fatalf("AddSubmoduleInContext() for case %d returned error, but should've not -> %s", index+1, err)
			}

		})
	}
}
