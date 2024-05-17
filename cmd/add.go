package cmd

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/actions"
	"github.com/jeftadlvw/git-nest/cmd/internal"
	"github.com/jeftadlvw/git-nest/migrations"
	mcontext "github.com/jeftadlvw/git-nest/migrations/context"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"github.com/spf13/cobra"
	"strings"
)

func createAddCmd() *cobra.Command {
	var addCmd = &cobra.Command{
		Use:   "add [url]",
		Short: "Add and clone a remote submodule into this project",
		Run:   internal.RunWrapper(wrapAddSubmodule, internal.ArgExactN(1)),
	}

	addCmd.Flags().StringP("ref", "r", "", "repository reference")
	addCmd.Flags().StringP("path", "p", "", "custom module path to clone into")

	return addCmd
}

func wrapAddSubmodule(cmd *cobra.Command, args []string) {

	var (
		url      urls.HttpUrl
		ref      string
		cloneDir models.Path
	)

	// validate url
	u, err := urls.HttpUrlFromString(args[0])
	if err != nil {
		fmt.Printf("error: invalid url\n")
		return
	}
	url = u

	refRaw, _ := cmd.Flags().GetString("ref")
	ref = strings.TrimSpace(refRaw)
	if ref == "" && ref != refRaw {
		fmt.Printf("error: no value defined for flag 'ref' \n")
	}

	cloneDirRaw, _ := cmd.Flags().GetString("path")
	cloneDir = models.Path(strings.TrimSpace(cloneDirRaw))
	if cloneDir.Empty() && cloneDir.String() != cloneDirRaw {
		fmt.Printf("error: no value defined for flag 'path' \n")
	}

	err = addSubmodule(url, ref, cloneDir)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func addSubmodule(url urls.HttpUrl, ref string, cloneDir models.Path) error {

	// acquire lock
	lockFile, err := internal.ErrorWrappedLockAcquiringAtProjectRootFromCwd()
	defer func() {
		err := internal.ErrorWrappedLockReleasing(lockFile)
		if err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		return err
	}

	// read context
	context, err := internal.ErrorWrappedEvaluateContext()
	if err != nil {
		return err
	}

	// run subcommand action
	actionMigrations, err := actions.AddSubmoduleInContext(&context, url, ref, cloneDir)
	if err != nil {
		return err
	}

	// run migrations
	actionMigrations = append(actionMigrations, mcontext.WriteConfigFiles{Context: &context})
	migrationError := migrations.RunMigrations(actionMigrations...)
	if migrationError != nil {
		return migrationError
	}

	return nil
}
