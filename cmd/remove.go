package cmd

import (
	"github.com/jeftadlvw/git-nest/actions"
	"github.com/jeftadlvw/git-nest/cmd/internal"
	"github.com/jeftadlvw/git-nest/migrations"
	mcontext "github.com/jeftadlvw/git-nest/migrations/context"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/spf13/cobra"
)

func createRemoveCommand() *cobra.Command {
	var rmCmd = &cobra.Command{
		Use:     "remove [path]",
		Aliases: []string{"rm"},
		Short:   "Remove a submodule from this project",
		RunE:    internal.RunWrapper(wrapRemoveSubmodule, internal.ArgExactN(1)),
	}

	rmCmd.Flags().BoolP("delete", "d", false, "delete existing directory")
	rmCmd.Flags().BoolP("force", "f", false, "force delete existing directory")

	return rmCmd
}

func wrapRemoveSubmodule(cmd *cobra.Command, args []string) error {
	deleteDirectory, _ := cmd.Flags().GetBool("delete")
	forceDelete, _ := cmd.Flags().GetBool("force")
	return removeSubmodule(models.Path(args[0]), deleteDirectory, forceDelete)
}

func removeSubmodule(p models.Path, deleteDirectory bool, forceDelete bool) error {
	// read context
	context, err := internal.ErrorWrappedEvaluateContext()
	if err != nil {
		return err
	}

	actionMigrations, err := actions.RemoveSubmoduleFromContext(&context, p, deleteDirectory, forceDelete)
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
