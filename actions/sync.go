package actions

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/migrations"
	"github.com/jeftadlvw/git-nest/migrations/git"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"github.com/jeftadlvw/git-nest/utils"
)

/*
SynchronizeConfigAndModules is a high level wrapper for synchronizing all changes between nested modules
and a changed configuration. It also updates the configuration in case the state in nested modules changes.
*/
func SynchronizeConfigAndModules(context *models.NestContext) ([]interfaces.Migration, error) {
	migrationChain := migrations.MigrationChain{}

	if !context.IsGitInstalled {
		return nil, errors.New("unable to synchronize if git is not installed")
	}

	for index := range len(context.Config.Submodules) {
		migrationArr, serr := SynchronizeSubmodule(&context.Config.Submodules[index], context.ProjectRoot)

		if serr != nil {
			fmt.Printf("could not synchronize nested module at position %d: %s\n", index+1, serr)
		}

		for _, migration := range migrationArr {
			migrationChain.Add(migration)
		}
	}

	return migrationChain.Migrations(), nil
}

/*
SynchronizeSubmodule is a high level wrapper to synchronize changes between one nested module
and an updated configuration.
*/
func SynchronizeSubmodule(s *models.Submodule, projectRoot models.Path) ([]interfaces.Migration, error) {
	migrationChain := migrations.MigrationChain{}

	absolutePath := projectRoot.Join(s.Path)

	// if s.Path is an already existing file, return error
	if absolutePath.IsFile() {
		return nil, fmt.Errorf("%s is a file", s.Path)
	}

	// if s.Path does not exist, then clone
	if !absolutePath.Exists() {
		migrationChain.Add(git.Clone{
			Url:          &s.Url,
			Path:         absolutePath.Parent(),
			CloneDirName: s.Path.Base(),
		})

		// and if s.Ref set, also perform checkout
		if s.Ref != "" {
			migrationChain.Add(git.Checkout{
				Path: absolutePath,
				Ref:  s.Ref,
			})
		}
	} else {
		// if s.Path already exists check the repository's origin url
		repositoryRemoteUrlStr, err := utils.GetGitRemoteUrl(absolutePath)
		if err != nil {
			return nil, fmt.Errorf("could not get remote url: %w", err)
		}

		repositoryRemoteUrl, err := urls.HttpUrlFromString(repositoryRemoteUrlStr)
		if err != nil {
			return nil, fmt.Errorf("internal error: could not parse url %s: %w", repositoryRemoteUrlStr, err)
		}

		// if origin url's do not match, choose repository as truth
		if repositoryRemoteUrl.String() != s.Url.String() {
			s.Url = repositoryRemoteUrl
		}

		// check the repository's head
		repositoryHeadLong, repositoryHeadAbbrev, err := utils.GetGitFetchHead(absolutePath)
		if err != nil {
			return nil, fmt.Errorf("could not get head: %w", err)
		}

		repositoryHead := repositoryHeadAbbrev
		if repositoryHead == "" {
			repositoryHead = repositoryHeadLong
		}

		// if the heads do not match, choose repository head as truth (== set submodule ref)
		if repositoryHead != s.Ref {
			s.Ref = repositoryHead
		}
	}

	return migrationChain.Migrations(), nil
}
