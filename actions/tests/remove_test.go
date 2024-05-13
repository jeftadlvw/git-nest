package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/actions"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/migrations"
	mcontext "github.com/jeftadlvw/git-nest/migrations/context"
	"github.com/jeftadlvw/git-nest/migrations/fs"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/test_env"
	test_env_models "github.com/jeftadlvw/git-nest/test_env/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"reflect"
	"testing"
)

func TestRemoveSubmoduleFromContext(t *testing.T) {

	repoDir := "example-repository"
	testFile := "testfileinrepo"
	testDirEmpty := "testdirinrepoempty"
	testDirFull := "testdirinrepofull"

	expectedMigrations := []interfaces.Migration{mcontext.RemoveSubmodule{}}
	expectedMigrationsDeleteDir := []interfaces.Migration{fs.DeleteDirectory{}, mcontext.RemoveSubmodule{}}

	cases := []struct {
		path                   string
		joinWithRoot           bool
		removeNonEmptyDir      bool
		addTempFile            bool
		commitTempFile         bool
		forceDelete            bool
		addPathToSubmodules    bool
		simulateNoGitInstalled bool
		expectedMigrations     []interfaces.Migration
		err                    bool
	}{
		// test path conditions
		{"/invalid/root", false, false, false, false, false, false, false, nil, true},
		{"../invalid/dir", false, false, false, false, false, false, false, nil, true},
		{"", true, false, false, false, false, false, false, nil, true},
		{".", true, false, false, false, false, false, false, nil, true},

		// find submodule
		{"foo", true, false, false, false, false, false, false, nil, true},
		{"foo", true, false, false, false, false, true, false, expectedMigrations, false},
		{repoDir, true, false, false, false, false, true, false, nil, true},
		{repoDir, true, true, false, false, false, true, false, expectedMigrationsDeleteDir, false},
		{repoDir, true, true, true, false, false, true, false, nil, true},
		{repoDir, true, true, true, false, true, true, false, expectedMigrationsDeleteDir, false},
		{repoDir, true, true, true, true, false, true, false, nil, true},
		{repoDir, true, true, true, true, true, true, false, expectedMigrationsDeleteDir, false},
		{repoDir, true, true, true, true, true, true, true, expectedMigrationsDeleteDir, false},
		{repoDir, true, true, false, false, false, true, true, expectedMigrationsDeleteDir, false},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestAddSubmoduleInContext-%d", index+1), func(t *testing.T) {
			t.Parallel()

			testEnv, err := test_env.CreateTestEnvironment(test_env_models.EnvSettings{Origin: test_env.RepoUrl, CloneDir: repoDir})
			if err != nil {
				t.Fatalf("error creating test environment for case %d: %s", index+1, err)
				return
			}
			defer testEnv.Destroy()

			if tc.addTempFile {
				// create test file and directory
				err = utils.WriteStrToFile(testEnv.Dir.SJoin(repoDir+"/"+testFile), "")
				if err != nil {
					t.Fatalf("error writing test file: %s", err)
				}

				absoluteRepoDirPath := testEnv.Dir.SJoin(repoDir)

				if tc.commitTempFile {
					out, err := utils.RunCommandCombinedOutput(absoluteRepoDirPath, "git", "add", ".")
					if err != nil {
						t.Fatalf("error git-adding test file for case %d: %s; %s", index+1, err, out)
						return
					}

					out, err = utils.RunCommandCombinedOutput(absoluteRepoDirPath, "git", "commit", "-m", "\"test commit\"")
					if err != nil {
						t.Fatalf("error git-commiting test file for case %d: %s; %s", index+1, err, out)
						return
					}
				}
			}

			localTestDirEmpty := testEnv.Dir.SJoin(testDirEmpty)
			localTestDirFull := testEnv.Dir.SJoin(testDirFull)

			err = os.Mkdir(localTestDirEmpty.String(), os.ModePerm)
			if err != nil {
				t.Fatalf("error creating test directory: %s", err)
			}

			err = os.Mkdir(localTestDirFull.String(), os.ModePerm)
			if err != nil {
				t.Fatalf("error creating test directory: %s", err)
			}

			err = utils.WriteStrToFile(localTestDirFull.SJoin(testFile), "")
			if err != nil {
				t.Fatalf("error writing test file: %s", err)
			}

			// create context
			context, err := internal.CreateContext(testEnv.Dir)
			if err != nil {
				t.Fatalf("error creating context for case %d: %s", index+1, err)
			}

			var submodules []models.Submodule
			if tc.addPathToSubmodules {
				submodules = append(submodules, models.Submodule{Path: models.Path(tc.path)})
			}

			context.Config.Submodules = submodules

			if tc.simulateNoGitInstalled {
				context.IsGitInstalled = false
			}

			p := models.Path(tc.path)
			if tc.joinWithRoot {
				p = context.ProjectRoot.SJoin(tc.path)
			}

			// remove submodule
			migrationArr, err := actions.RemoveSubmoduleFromContext(&context, p, tc.removeNonEmptyDir, tc.forceDelete)

			// check migration array
			if !tc.err && len(tc.expectedMigrations) != len(migrationArr) {
				t.Fatalf("AddSubmoduleInContext() for case %d returned unequal amounts of migrations: expected %d, got %d", index+1, len(tc.expectedMigrations), len(migrationArr))
			}
			if !tc.err {
				for mindex, migration := range migrationArr {
					if reflect.TypeOf(migration) != reflect.TypeOf(tc.expectedMigrations[mindex]) {
						t.Fatalf("AddSubmoduleInContext() for case %d had unexpected migration at index %d: %T != %T", index+1, mindex, migration, tc.expectedMigrations[mindex])
					}
				}
			}

			// run migrations if action function call did not return an error
			if err == nil {
				err = migrations.RunMigrations(migrationArr...)
			}

			// test for errors
			if tc.err && err == nil {
				t.Fatalf("RemoveSubmoduleFromContext-%d returned no error but expected one", index+1)
			}
			if !tc.err && err != nil {
				t.Fatalf("RemoveSubmoduleFromContext-%d returned error, but should've not -> %s", index+1, err)
			}
		})
	}
}
