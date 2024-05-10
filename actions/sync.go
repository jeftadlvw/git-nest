package actions

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/migrations"
	"github.com/jeftadlvw/git-nest/models"
)

/*
SynchronizeConfigAndModules is a high level wrapper for synchronizing changes between nested modules
and a changed configuration. It also updates the configuration in case the state in nested modules changes.
*/
func SynchronizeConfigAndModules(context *models.NestContext) ([]interfaces.Migration, error) {
	var (
		err            error
		migrationChain = migrations.MigrationChain{}
	)

	fmt.Println("Hello sync!", err)

	return migrationChain.Migrations(), nil
}
