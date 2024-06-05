package cmd

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/actions"
	"github.com/jeftadlvw/git-nest/cmd/internal"
	"github.com/jeftadlvw/git-nest/migrations"
	mcontext "github.com/jeftadlvw/git-nest/migrations/context"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: fmt.Sprintf("Update and apply state changes"),
	RunE:  internal.RunWrapper(wrapSync),
}

func wrapSync(cmd *cobra.Command, args []string) error {
	return sync()
}

func sync() error {
	// read context
	context, err := internal.ErrorWrappedEvaluateContext()
	if err != nil {
		return err
	}

	actionMigrations, err := actions.SynchronizeConfigAndModules(&context)
	if err != nil {
		return err
	}

	actionMigrations = append(actionMigrations, mcontext.WriteConfigFiles{Context: &context})
	migrationError := migrations.RunMigrations(actionMigrations...)
	if migrationError != nil {
		return migrationError
	}

	return nil
}
