package cmd

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/actions"
	cmdInternal "github.com/jeftadlvw/git-nest/cmd/internal"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/models/urls"
	"github.com/spf13/cobra"
	"strings"
)

func createAddCmd() *cobra.Command {
	var infoCmd = &cobra.Command{
		Use:   "add url [ref] [location]",
		Short: "Adds and clones a remote submodule into this project",
		Run:   cmdInternal.RunWrapper(wrap, cmdInternal.ArgMinN(1), cmdInternal.ArgMaxN(3)),
	}

	return infoCmd
}

func wrap(args []string) {
	url, ref, cloneDir, err := argsToParams(args)
	if err != nil {
		fmt.Printf("fatal: argument validation: %s\n", err)
		return
	}

	err = addSubmodule(url, ref, cloneDir)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func argsToParams(args []string) (urls.HttpUrl, string, string, error) {
	var (
		url      urls.HttpUrl
		ref      string
		cloneDir string
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
		cloneDir = strings.TrimSpace(args[1])
	}
	if argLen >= 3 {
		ref = strings.TrimSpace(args[2])
	}

	return url, ref, cloneDir, nil
}

func addSubmodule(url urls.HttpUrl, ref string, cloneDir string) error {

	// read context
	// current configuration might have some errors which will be needed to fix first
	// if no error, defer writing context
	context, err := internal.EvaluateContext()
	if err != nil {
		return fmt.Errorf("internal context error: %w.\nPlease fix any configuration errors to proceed", err)
	}

	err = actions.AddSubmoduleInContext(&context, url, ref, cloneDir)
	if err != nil {
		return err
	}

	// write configuration files
	_, _, err1, err2 := internal.WriteProjectConfigFiles(context)
	if err1 != nil {
		return fmt.Errorf("error writing configuration files: %w", err1)
	} else if err2 != nil {
		return fmt.Errorf("error writing configuration files: %w", err2)
	}

	return nil
}
