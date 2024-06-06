package cmd

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/actions"
	cmdInternal "github.com/jeftadlvw/git-nest/cmd/internal"
	"github.com/jeftadlvw/git-nest/migrations"
	"github.com/spf13/cobra"
)

func createPullCommand() *cobra.Command {
	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull new updates in all nested modules",
		RunE:  cmdInternal.RunWrapper(wrapGitPullModules, cmdInternal.ArgNone()),
	}

	return pullCmd
}

func wrapGitPullModules(cmd *cobra.Command, args []string) error {
	return gitPullModules()
}

func gitPullModules() error {
	// read context
	context, err := cmdInternal.ErrorWrappedEvaluateContext()
	if err != nil {
		return err
	}

	if len(context.Config.Submodules) == 0 {
		fmt.Println(cmdInternal.NoNestedModulesMsg)
		return nil
	}

	actionMigrations, err := actions.PullAllSubmodules(&context)
	if err != nil {
		return err
	}

	migrationError := migrations.RunMigrations(actionMigrations...)
	if migrationError != nil {
		return migrationError
	}

	return nil
}
