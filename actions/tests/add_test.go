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
		{test_env.RepoUrl, "", "", []models.Submodule{{"example-repository", &testRepoUrl, ""}}, false, false, nil, true},
		{test_env.RepoUrl, "", "", []models.Submodule{{"example-repository", &testRepoUrl, ""}}, false, true, nil, true},
		{test_env.RepoUrl, "", test_env.RepoBranch1, []models.Submodule{{"foo", &testRepoUrl, ""}}, false, false, nil, true},
		{test_env.RepoUrl, "", test_env.RepoBranch1, []models.Submodule{{"foo", &testRepoUrl, ""}}, false, true, expectedMigrationsRef, false},
		{test_env.RepoUrl, "", test_env.RepoBranch1, []models.Submodule{{"foo", &testRepoUrl, test_env.RepoBranch1}}, false, true, expectedMigrationsRef, false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestAddSubmoduleInContext-%d", index+1), func(t *testing.T) {
			t.Parallel()

			tempDir := models.Path(t.TempDir())
			err := test_env.CreateTestEnvironment(tempDir, test_env_models.EnvSettings{EmptyGit: true})
			if err != nil {
				t.Fatalf("error creating test environment: %s", err)
				return
			}

			// create test file and directory
			err = utils.WriteStrToFile(tempDir.SJoin(testFile), "")
			if err != nil {
				t.Fatalf("error writing test file: %s", err)
			}

			localTestDirEmpty := tempDir.SJoin(testDirEmpty)
			localTestDirFull := tempDir.SJoin(testDirFull)
			err = os.Mkdir(localTestDirEmpty.String(), os.ModePerm)
			err = os.Mkdir(localTestDirFull.String(), os.ModePerm)

			err = utils.WriteStrToFile(localTestDirFull.SJoin(testFile), "")
			if err != nil {
				t.Fatalf("error writing test file: %s", err)
			}

			// create context
			context, err := internal.CreateContext(tempDir)
			if err != nil {
				t.Fatalf("error creating context: %s", err)
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
				t.Fatalf("unable to convert url: %s", err)
			}

			cloneDir := tc.cloneDir
			migrationArr, err := actions.AddSubmoduleInContext(&context, url, tc.ref, models.Path(cloneDir))

			// check migration array
			if !tc.err && len(tc.expectedMigrations) != len(migrationArr) {
				t.Fatalf("got unequal amounts of migrations: expected %d, got %d", len(tc.expectedMigrations), len(migrationArr))
			}
			if !tc.err {
				for mindex, migration := range migrationArr {
					if reflect.TypeOf(migration) != reflect.TypeOf(tc.expectedMigrations[mindex]) {
						t.Fatalf("unexpected migration at index %d: %T != %T", mindex, migration, tc.expectedMigrations[mindex])
					}
				}
			}

			// run migrations if action function call did not return an error
			if err == nil {
				err = migrations.RunMigrations(migrationArr...)
			}

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
