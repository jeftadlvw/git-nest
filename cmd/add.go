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
		Use:   "add url [ref] [location]",
		Short: "Add and clone a remote submodule into this project",
		Run:   internal.RunWrapper(wrapAddSubmodule, internal.ArgMinN(1), internal.ArgMaxN(3)),
	}

	return addCmd
}

func wrapAddSubmodule(cmd *cobra.Command, args []string) {
	url, ref, cloneDir, err := argsToParamsAddSubmodule(args)
	if err != nil {
		fmt.Printf("fatal: argument validation: %s\n", err)
		return
	}

	err = addSubmodule(url, ref, cloneDir)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func argsToParamsAddSubmodule(args []string) (urls.HttpUrl, string, models.Path, error) {
	var (
		url      urls.HttpUrl
		ref      string
		cloneDir models.Path
	)

	var argLen = len(args)

	if argLen >= 1 {
		u, err := urls.HttpUrlFromString(args[0])
		if err != nil {
			return urls.HttpUrl{}, "", "", fmt.Errorf("invalid url")
		}
		url = u
	}
	if argLen >= 2 {
		ref = strings.TrimSpace(args[1])
	}
	if argLen >= 3 {
		cloneDir = models.Path(strings.TrimSpace(args[2]))
	}

	return url, ref, cloneDir, nil
}

func addSubmodule(url urls.HttpUrl, ref string, cloneDir models.Path) error {

	// read context
	context, err := internal.ErrorWrappedEvaluateContext()
	if err != nil {
		return err
	}

	actionMigrations, err := actions.AddSubmoduleInContext(&context, url, ref, cloneDir)
	if err != nil {
		return err
	}

	actionMigrations = append(actionMigrations, mcontext.WriteConfigFiles{Context: &context})
	err = migrations.RunMigrations(actionMigrations...)
	if err != nil {
		return err
	}

	return nil
}
