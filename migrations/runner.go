package migrations

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/interfaces"
)

func RunMigrations(migrations ...interfaces.Migration) error {
	for index, migration := range migrations {
		if err := migration.Migrate(); err != nil {
			return fmt.Errorf("%w (migration #%d)", err, index+1)
		}
	}

	return nil
}
