package actions

import (
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/migrations"
	"github.com/jeftadlvw/git-nest/migrations/git"
	"github.com/jeftadlvw/git-nest/models"
)

/*
PullAllSubmodules is a wrapper that adds a git.Pull migration for every nested module.
*/
func PullAllSubmodules(context *models.NestContext) ([]interfaces.Migration, error) {
	migrationChain := migrations.MigrationChain{}

	for _, submodule := range context.Config.Submodules {
		migrationChain.Add(git.Pull{Path: submodule.Path})
	}

	return migrationChain.Migrations(), nil
}
