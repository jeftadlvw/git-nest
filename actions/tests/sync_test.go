package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/actions"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/migrations"
	"github.com/jeftadlvw/git-nest/migrations/git"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"reflect"
	"testing"
)

func TestSynchronizeSubmodule(t *testing.T) {

	testFile := "testfileinrepo"
	testDir := "testdir"

	cases := []struct {
		submodule          models.Submodule
		create             bool
		originOverride     string
		refOverride        string
		expectedMigrations []interfaces.Migration
		err                bool
	}{}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestSynchronizeSubmodule-%d", index+1), func(t *testing.T) {
			t.Parallel()

			testEnvDir := models.Path(t.TempDir())

			// create test directory
			localTestDir := testEnvDir.SJoin(testDir)
			err := os.Mkdir(localTestDir.String(), os.ModePerm)
			if err != nil {
				t.Fatalf("error creating test directory: %s", err)
			}

			// create test file
			err = utils.WriteStrToFile(testEnvDir.SJoin(testFile), "")
			if err != nil {
				t.Fatalf("error writing test file: %s", err)
			}

			// create context
			context, err := internal.CreateContext(testEnvDir)
			if err != nil {
				t.Fatalf("error creating context for case %d: %s", index+1, err)
			}

			// create submodule
			if tc.create {
				submodulePath := testEnvDir.Join(tc.submodule.Path)
				migration := git.Clone{
					Url:          &tc.submodule.Url,
					Path:         submodulePath.Parent(),
					CloneDirName: tc.submodule.Path.Base(),
				}

				err = migration.Migrate()
				if err != nil {
					t.Fatalf("error pre-creating submodule: %s", err)
				}

				if tc.originOverride != "" {
					_, err = utils.RunCommandCombinedOutput("git", "remote", "set-url", "origin", tc.refOverride)
					if err != nil {
						t.Fatalf("error setting remote url: %s", err)
					}
				}

				if tc.refOverride != "" {
					err = git.Checkout{
						Path: submodulePath,
						Ref:  tc.refOverride,
					}.Migrate()
					if err != nil {
						t.Fatalf("error changing ref: %s", err)
					}
				}
			}

			// sync submodule
			migrationArr, err := actions.SynchronizeSubmodule(&tc.submodule, context.ProjectRoot)

			// check migration array
			if !tc.err && len(tc.expectedMigrations) != len(migrationArr) {
				t.Fatalf("TestSynchronizeSubmodule() for case %d returned unequal amounts of migrations: expected %d, got %d", index+1, len(tc.expectedMigrations), len(migrationArr))
			}
			if !tc.err {
				for mindex, migration := range migrationArr {
					if reflect.TypeOf(migration) != reflect.TypeOf(tc.expectedMigrations[mindex]) {
						t.Fatalf("TestSynchronizeSubmodule() for case %d had unexpected migration at index %d: %s != %s", index+1, mindex, reflect.TypeOf(migration), reflect.TypeOf(tc.expectedMigrations[mindex]))
					}
				}
			}

			if tc.originOverride != "" {
				url, err := urls.HttpUrlFromString(tc.originOverride)
				if err != nil {
					t.Fatalf("error creating override url: %s", err)
				}

				if tc.submodule.Url.String() != url.String() {
					t.Fatalf("submodule url was not overridden: expected %s, got %s", url.String(), tc.submodule.Url.String())
				}
			}

			if tc.refOverride != "" && tc.submodule.Ref != tc.refOverride {
				t.Fatalf("submodule ref was overridden: expected %s, got %s", tc.refOverride, tc.submodule.Ref)
			}

			// run migrations if action function call did not return an error
			if err == nil {
				err = migrations.RunMigrations(migrationArr...)
			}

			// test for errors
			if tc.err && err == nil {
				t.Fatalf("TestSynchronizeSubmodule-%d returned no error but expected one", index+1)
			}
			if !tc.err && err != nil {
				t.Fatalf("TestSynchronizeSubmodule-%d returned error, but should've not -> %s", index+1, err)
			}
		})
	}
}
