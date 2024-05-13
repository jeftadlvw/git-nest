package migrations

import "github.com/jeftadlvw/git-nest/interfaces"

type MigrationChain struct {
	migrations []interfaces.Migration
}

func (mc *MigrationChain) Add(m interfaces.Migration) {
	mc.migrations = append(mc.migrations, m)
}

func (mc *MigrationChain) Migrations() []interfaces.Migration {
	return mc.migrations
}
