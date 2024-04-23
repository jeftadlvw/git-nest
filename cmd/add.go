package cmd

import (
	"fmt"
	cmdInternal "github.com/jeftadlvw/git-nest/cmd/internal"
	internal "github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"github.com/jeftadlvw/git-nest/utils"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

func createAddCmd() *cobra.Command {
	var infoCmd = &cobra.Command{
		Use:   "add url [path] [ref]",
		Short: "Adds and clones a remote submodule into this project",
		Run:   cmdInternal.RunWrapper(wrap, cmdInternal.ArgMinN(1), cmdInternal.ArgMaxN(3)),
	}

	return infoCmd
}

func wrap(args []string) {
	url, path, ref, err := argsToParams(args)
	if err != nil {
		fmt.Printf("fatal: argument validation: %s\n", err)
		return
	}

	err = addSubmodule(url, path, ref)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func argsToParams(args []string) (urls.HttpUrl, models.Path, string, error) {
	var (
		url  urls.HttpUrl
		path models.Path
		ref  string
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
		path = models.Path(args[1])
		path = path.Clean()
	}
	if argLen >= 3 {
		ref = strings.TrimSpace(args[2])
	}

	return url, path, ref, nil
}

func addSubmodule(url urls.HttpUrl, path models.Path, ref string) error {

	fmt.Printf("Hello subcommand! %v %s %s\n", url, path, ref)

	// read context
	// current configuration might have some errors which will be needed to fix first
	// if no error, defer writing context
	context, err := internal.EvaluateContext()
	if err != nil {
		return fmt.Errorf("internal context error: %w.\nPlease fix any configuration errors to proceed", err)
	}

	// if path is empty, set it to the url's last part
	if path.Empty() {
		path = models.Path(filepath.Base(url.Path))
	}

	// check if path escapes project directory
	absolutePath := context.WorkingDirectory.Join(path)
	relativeToProjectRoot := context.ProjectRoot.Join(absolutePath)
	if strings.Contains(relativeToProjectRoot.String(), "..") {
		return fmt.Errorf("validation error: %s escapes the project root", path)
	}

	// check if path is existing file
	if absolutePath.IsFile() {
		return fmt.Errorf("validation error: %s is a file", path)
	}

	// if path is a directory, ensure it's empty
	if absolutePath.BContains("*") {
		return fmt.Errorf("validation error: %s is not an empty directory", path)
	}

	newSubmodule := models.Submodule{
		Path: relativeToProjectRoot,
		Url:  url,
		Ref:  ref,
	}
	appendedSubmoduleSlice := append(context.Config.Submodules, newSubmodule)

	// check if url already exists in config
	// based on the settings, check if ref also already exists
	err = models.CheckForDuplicateSubmodules(context.Config.Config.AllowDuplicateOrigins, appendedSubmoduleSlice...)
	if err != nil {
		return fmt.Errorf("validation error for new submodule: %s", err)
	}

	// check if git is installed
	if !context.IsGitInstalled {
		return fmt.Errorf("please install git in order to add a submodule")
	}

	// clone the repository
	err = utils.CloneGitRepository(newSubmodule.Url.String(), newSubmodule.Path.Parent(), newSubmodule.Path.Base())
	if err != nil {
		return fmt.Errorf("error while cloning %s: %s", newSubmodule.Url.String(), err)
	}

	// change ref
	if newSubmodule.Ref != "" {
		err = utils.ChangeGitHead(newSubmodule.Path, newSubmodule.Ref)
		if err != nil {
			return fmt.Errorf("error while changing ref: %s", err)
		}
	}

	// add submodule to config
	context.Config.Submodules = appendedSubmoduleSlice

	// write configuration files
	_, _, err1, err2 := internal.WriteProjectConfigFiles(context)
	if err1 != nil {
		return fmt.Errorf("error writing configuration files: %w", err1)
	} else if err2 != nil {
		return fmt.Errorf("error writing configuration files: %w", err2)
	}

	return nil
}
