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
	Run:   internal.RunWrapper(wrapSync),
}

func wrapSync(cmd *cobra.Command, args []string) {
	err := sync()
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
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
	err = migrations.RunMigrations(actionMigrations...)
	if err != nil {
		return err
	}

	return nil
}
