package actions

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"path/filepath"
	"strings"
)

/*
AddSubmoduleInContext is a high-level wrapper that adds a submodule into a context,
checking for duplicates before cloning the repository.
*/
func AddSubmoduleInContext(context *models.NestContext, url urls.HttpUrl, ref string, cloneDir models.Path) error {
	var err error

	// check if git is installed
	if !context.IsGitInstalled {
		return fmt.Errorf("please install git in order to add a submodule")
	}

	var relativeToRoot = cloneDir
	var calcRelativeTo models.Path
	var absolutePath models.Path
	var repositoryName = strings.TrimSuffix(filepath.Base(url.String()), ".git")

	// if cloneDir is empty, set it to repository name
	// if cloneDir has trailing separator, append repository name
	if strings.HasSuffix(string(cloneDir), string(filepath.Separator)) {
		fmt.Println("has suffix")
		cloneDir = cloneDir.SJoin(repositoryName)
	} else if cloneDir.Empty() {
		cloneDir = models.Path(repositoryName)
	}

	// check if cloneDir is absolute
	// if cloneDir is absolute, then calculate relative to project root
	// else, expect cloneDir to be relative to cwd, so join it with cwd and then calculate relative path to project root
	if filepath.IsAbs(cloneDir.String()) {
		calcRelativeTo = cloneDir
	} else {
		calcRelativeTo = context.WorkingDirectory.Join(cloneDir)
	}

	relativeToRoot, err = context.ProjectRoot.Relative(calcRelativeTo)
	if err != nil {
		return fmt.Errorf("internal error: could not find relative to project root: %w", err)
	}

	// check if relative path escapes project root
	if strings.Contains(relativeToRoot.String(), "..") {
		return fmt.Errorf("validation error: %s escapes the project root", cloneDir)
	}

	// join project root and absolute path, check if it's not an existing file and create that directory
	absolutePath = context.ProjectRoot.Join(relativeToRoot)

	if !absolutePath.Exists() {
		err = os.MkdirAll(absolutePath.String(), os.ModePerm)
		if err != nil {
			return fmt.Errorf("internal error: could not create directory %s: %w", absolutePath, err)
		}
	} else {
		if absolutePath.IsFile() {
			return fmt.Errorf("validation error: %s is a file", cloneDir)
		}
		if absolutePath.BContains("*") {
			return fmt.Errorf("validation error: %s is not empty", cloneDir)
		}
	}

	newSubmodule := models.Submodule{
		Path: relativeToRoot,
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
	err = utils.CloneGitRepository(newSubmodule.Url.String(), absolutePath.Parent(), absolutePath.Base())
	if err != nil {
		return fmt.Errorf("error while cloning into %s: %s", absolutePath, err)
	}

	// change ref
	if newSubmodule.Ref != "" {

		localSubmoduleClonePath := relativeToRoot.String()
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
