package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/actions"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/migrations"
	"github.com/jeftadlvw/git-nest/migrations/git"
	"github.com/jeftadlvw/git-nest/migrations/submodules"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"github.com/jeftadlvw/git-nest/test_env"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSynchronizeSubmodule(t *testing.T) {

	testFile := "testfileinrepo"
	testDir := "testdir"
	exampleSubmodule := models.Submodule{
		Path: "nested_module-1",
		Url:  &urls.HttpUrl{HostnameS: "github.com", Port: 443, PathS: "/jeftadlvw/example-repository", Secure: true},
		Ref:  "",
	}
	exampleSubmoduleRef := models.Submodule{
		Path: "nested_module-1",
		Url:  &urls.HttpUrl{HostnameS: "github.com", Port: 443, PathS: "/jeftadlvw/example-repository", Secure: true},
		Ref:  test_env.RepoBranch1,
	}
	exampleSubmoduleRefDefault := models.Submodule{
		Path: "nested_module-1",
		Url:  &urls.HttpUrl{HostnameS: "github.com", Port: 443, PathS: "/jeftadlvw/example-repository", Secure: true},
		Ref:  test_env.RepoBranchDefault,
	}

	createMigration := []interfaces.Migration{git.Clone{}}
	createAndCheckoutMigration := []interfaces.Migration{git.Clone{}, git.Checkout{}}
	updateUrlMigration := []interfaces.Migration{submodules.UpdateUrl{}}
	updateRefMigration := []interfaces.Migration{submodules.UpdateRef{}}
	updateUrlAndRefMigration := []interfaces.Migration{submodules.UpdateUrl{}, submodules.UpdateRef{}}

	cases := []struct {
		submodule          models.Submodule
		create             bool
		repoOriginOverride string
		repoRefOverride    string
		expectedMigrations []interfaces.Migration
		err                bool
	}{
		{models.Submodule{}, false, "", "", nil, true},
		{exampleSubmoduleRefDefault, false, "", "", createAndCheckoutMigration, false},
		{exampleSubmodule, false, "", "", createMigration, false},
		{exampleSubmoduleRefDefault, true, "", "", nil, false},
		{exampleSubmodule, true, "", test_env.RepoBranch1, updateRefMigration, false},
		{exampleSubmoduleRef, false, "", "", createAndCheckoutMigration, false},
		{exampleSubmodule, true, "", "", updateRefMigration, false},
		{exampleSubmoduleRef, true, "http://example.com/foo", test_env.RepoBranchDefault, updateUrlAndRefMigration, false},
		{exampleSubmoduleRefDefault, true, "http://example.com/foo", "", updateUrlMigration, false},
		{exampleSubmoduleRef, true, "", test_env.RepoCommit, updateRefMigration, false},
		{exampleSubmodule, true, "", test_env.RepoCommit, updateRefMigration, false},
		{exampleSubmoduleRef, true, "http://example.com/foo", test_env.RepoCommit, updateUrlAndRefMigration, false},
		{exampleSubmodule, true, "http://example.com/foo", test_env.RepoCommit, updateUrlAndRefMigration, false},
	}

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
				err = git.Clone{
					Url:          tc.submodule.Url,
					Path:         submodulePath.Parent(),
					CloneDirName: tc.submodule.Path.Base(),
				}.Migrate()

				if tc.repoRefOverride == "" && tc.submodule.Ref != "" {
					err = git.Checkout{
						Path: submodulePath,
						Ref:  tc.submodule.Ref,
					}.Migrate()
					if err != nil {
						t.Fatalf("error changing ref: %s", err)
					}
				}

				if err != nil {
					t.Fatalf("error pre-creating submodule: %s", err)
				}

				if tc.repoOriginOverride != "" {
					_, err = utils.RunCommandCombinedOutput(submodulePath, "git", "remote", "set-url", "origin", tc.repoOriginOverride)
					if err != nil {
						t.Fatalf("error setting remote url: %s", err)
					}
				}

				if tc.repoRefOverride != "" {
					err = git.Checkout{
						Path: submodulePath,
						Ref:  tc.repoRefOverride,
					}.Migrate()
					if err != nil {
						t.Fatalf("error changing ref: %s", err)
					}
				}
			}

			// sync submodule
			migrationArr, err := actions.SynchronizeSubmodule(&tc.submodule, context.ProjectRoot)

			if tc.err && err == nil {
				t.Fatalf("TestSynchronizeSubmodule-%d returned no error but expected one", index+1)
			}
			if !tc.err && err != nil {
				t.Fatalf("TestSynchronizeSubmodule-%d returned error, but should've not -> %s", index+1, err)
			}

			// check migration array
			if !tc.err && len(tc.expectedMigrations) != len(migrationArr) {
				typeConcatStr := " "
				for _, m := range migrationArr {
					typeConcatStr += fmt.Sprintf("%T %v ", m, m)
				}
				t.Fatalf("TestSynchronizeSubmodule() for case %d returned unequal amounts of migrations: expected %d, got %d (%s)", index+1, len(tc.expectedMigrations), len(migrationArr), typeConcatStr)
			}
			if !tc.err {
				for mindex, migration := range migrationArr {
					if reflect.TypeOf(migration) != reflect.TypeOf(tc.expectedMigrations[mindex]) {
						t.Fatalf("TestSynchronizeSubmodule() for case %d had unexpected migration at index %d: %T != %T", index+1, mindex, migration, tc.expectedMigrations[mindex])
					}
				}
			}

			// run migrations if action function call did not return an error
			if err == nil {
				err = migrations.RunMigrations(migrationArr...)
			}

			// test if sync migrations were successful
			if tc.repoOriginOverride != "" {
				url, err := urls.HttpUrlFromString(tc.repoOriginOverride)
				if err != nil {
					t.Fatalf("error creating override url: %s", err)
				}

				if tc.submodule.Url.String() != url.String() {
					t.Fatalf("submodule url was not overridden: expected >%s<, got >%s<", url.String(), tc.submodule.Url.String())
				}
			}

			if tc.repoRefOverride != "" && !strings.HasPrefix(tc.submodule.Ref, tc.repoRefOverride) {
				t.Fatalf("submodule ref was not overridden: expected >%s<, got >%s<", tc.repoRefOverride, tc.submodule.Ref)
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
