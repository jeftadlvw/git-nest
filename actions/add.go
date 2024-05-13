package actions

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/migrations"
	mcontext "github.com/jeftadlvw/git-nest/migrations/context"
	"github.com/jeftadlvw/git-nest/migrations/git"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"os"
	"path/filepath"
	"strings"
)

/*
AddSubmoduleInContext is a high-level wrapper that adds a submodule into a context,
checking for duplicates before cloning the repository.
*/
func AddSubmoduleInContext(context *models.NestContext, url urls.HttpUrl, ref string, cloneDir models.Path) ([]interfaces.Migration, error) {

	var (
		err            error
		migrationChain = migrations.MigrationChain{}
	)

	// check if git is installed
	if !context.IsGitInstalled {
		return nil, fmt.Errorf("please install git in order to add a submodule")
	}

	var relativeToRoot = cloneDir
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

	relativeToRoot, err = internal.PathRelativeToRootWithJoinedOriginIfNotAbs(context.ProjectRoot, context.WorkingDirectory, cloneDir)
	if err != nil {
		return nil, fmt.Errorf("internal error: could not find relative to project root: %w", err)
	}

	// check if relative path escapes project root
	if internal.PathContainsUp(relativeToRoot) {
		return nil, fmt.Errorf("validation error: %s escapes the project root", cloneDir)
	}

	// join project root and absolute path, check if it's not an existing file and create that directory
	absolutePath = context.ProjectRoot.Join(relativeToRoot)

	if !absolutePath.Exists() {
		err = os.MkdirAll(absolutePath.String(), os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("internal error: could not create directory %s: %w", absolutePath, err)
		}
	} else {
		if absolutePath.IsFile() {
			return nil, fmt.Errorf("validation error: %s is a file", cloneDir)
		}
		if absolutePath.BContains("*") {
			return nil, fmt.Errorf("validation error: %s is not empty", cloneDir)
		}
	}

	newSubmodule := models.Submodule{
		Path: relativeToRoot,
		Url:  &url,
		Ref:  ref,
	}

	// append submodule and clone it
	migrationChain.Add(mcontext.AppendSubmodule{Context: context, Submodule: newSubmodule})
	migrationChain.Add(git.Clone{Url: newSubmodule.Url, Path: absolutePath.Parent(), CloneDirName: absolutePath.Base()})

	if newSubmodule.Ref != "" {
		localSubmoduleClonePath := relativeToRoot.String()
		if localSubmoduleClonePath == "" {
			localSubmoduleClonePath = strings.TrimSuffix(filepath.Base(url.String()), ".git")
		}

		migrationChain.Add(git.Checkout{Path: context.ProjectRoot.SJoin(localSubmoduleClonePath), Ref: newSubmodule.Ref})
	}

	return migrationChain.Migrations(), nil
}
