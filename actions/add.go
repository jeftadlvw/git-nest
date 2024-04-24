package actions

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"github.com/jeftadlvw/git-nest/utils"
	"path/filepath"
	"strings"
)

/*
AddSubmoduleInContext is a high-level wrapper that adds a submodule into a context,
checking for duplicates before cloning the repository.
*/
func AddSubmoduleInContext(context *models.NestContext, url urls.HttpUrl, ref string, cloneDirName string) error {
	var err error

	// check if git is installed
	if !context.IsGitInstalled {
		return fmt.Errorf("please install git in order to add a submodule")
	}

	// check if cloneDirName escapes project directory
	absolutePath := context.WorkingDirectory.SJoin(cloneDirName)
	relativeToProjectRoot, err := context.ProjectRoot.Relative(absolutePath)
	if err != nil {
		return fmt.Errorf("internal error: could not find relative to project root: %w", err)
	}

	if strings.Contains(relativeToProjectRoot.String(), "..") {
		return fmt.Errorf("validation error: %s escapes the project root", cloneDirName)
	}

	// check if cloneDirName is existing file
	if absolutePath.IsFile() {
		return fmt.Errorf("validation error: %s is a file", cloneDirName)
	}

	newSubmodule := models.Submodule{
		Path: relativeToProjectRoot,
		Url:  url,
		Ref:  ref,
	}

	// validate submodule's content
	err = newSubmodule.Validate()
	if err != nil {
		return fmt.Errorf("validation error: %s", err)
	}

	appendedSubmoduleSlice := append(context.Config.Submodules, newSubmodule)

	// check if url already exists in config
	// based on the settings, check if ref also already exists
	err = models.CheckForDuplicateSubmodules(context.Config.Config.AllowDuplicateOrigins, appendedSubmoduleSlice...)
	if err != nil {
		return fmt.Errorf("validation error for new submodule: %s", err)
	}

	// clone the repository
	err = utils.CloneGitRepository(newSubmodule.Url.String(), context.ProjectRoot, newSubmodule.Path.String())
	if err != nil {
		return fmt.Errorf("error while cloning: %s", err)
	}

	// change ref
	if newSubmodule.Ref != "" {

		localSubmoduleClonePath := relativeToProjectRoot.String()
		if localSubmoduleClonePath == "" {
			localSubmoduleClonePath = strings.TrimSuffix(filepath.Base(url.String()), ".git")
		}

		err = utils.ChangeGitHead(context.ProjectRoot.SJoin(localSubmoduleClonePath), newSubmodule.Ref)
		if err != nil {
			return fmt.Errorf("error while changing ref: %s", err)
		}
	}

	// add submodule to config
	context.Config.Submodules = appendedSubmoduleSlice

	return nil
}
