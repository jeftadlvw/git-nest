package actions

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/migrations"
	mcontext "github.com/jeftadlvw/git-nest/migrations/context"
	"github.com/jeftadlvw/git-nest/migrations/fs"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
)

/*
RemoveSubmoduleFromContext is a high-level wrapper that removes a submodule from the context.
It also removes any existing directories if there are no changes (unless forced).
*/
func RemoveSubmoduleFromContext(context *models.NestContext, p models.Path, removeDir bool, forceDelete bool) ([]interfaces.Migration, error) {

	var (
		err            error
		migrationChain = migrations.MigrationChain{}
	)

	relativeToRoot, err := internal.PathRelativeToRootWithJoinedOriginIfNotAbs(context.ProjectRoot, context.WorkingDirectory, p)
	if err != nil {
		return nil, fmt.Errorf("internal error: could not find relative to project root: %w", err)
	}

	// check if relative path escapes project root
	if internal.PathContainsUp(relativeToRoot) {
		return nil, fmt.Errorf("validation error: %s escapes the project root", p)
	}

	absolutePath := context.ProjectRoot.Join(relativeToRoot)
	if absolutePath.Equals(context.ProjectRoot) {
		return nil, fmt.Errorf("validation error: path cannot be project root")
	}

	// check if passed submodule exists in context
	var removeIndex = -1
	for i, submodule := range context.Config.Submodules {
		if submodule.Path.String() == relativeToRoot.String() {
			removeIndex = i
		}
	}

	// return error if no match found
	if removeIndex == -1 {
		return nil, fmt.Errorf("passed submodule does not exist: %s", relativeToRoot)
	}

	// check if directory exists
	if absolutePath.IsDir() {
		if absolutePath.BContains("*") && !removeDir {
			return nil, fmt.Errorf("submodule directry at %s is not empty.\nUse git-nest rm [path] -d to remove it.", p)
		}

		if context.IsGitInstalled {
			// check if repository has untracked changes
			hasUntrackedChanges, err := utils.GetGitHasUntrackedChanges(absolutePath)
			if err != nil {
				return nil, fmt.Errorf("internal error: could not check if uncommitted changes exist: %w", err)
			}

			if hasUntrackedChanges && !forceDelete {
				return nil, fmt.Errorf("submodule repository at %s contains uncommitted changes.\nCommit and push your changes or use git-nest rm [path] -df to forcefully remove it.", p)
			}

			// check if repository has unpublished changes
			hasUnpublishedChanges, err := utils.GetGitHasUnpublishedChanges(absolutePath)
			if err != nil {
				return nil, fmt.Errorf("internal error: could not check if unpushed commits exist: %w", err)
			}

			if hasUnpublishedChanges && !forceDelete {
				return nil, fmt.Errorf("submodule repository at %s has unpushed commits.\nPush your changes or use git-nest rm [path] -df to forcefully remove it.", p)
			}
		}

		migrationChain.Add(fs.DeleteDirectory{Path: absolutePath})
	}

	// remove submodule from context
	migrationChain.Add(mcontext.RemoveSubmodule{Context: context, SubmoduleIndex: removeIndex})

	return migrationChain.Migrations(), nil
}
