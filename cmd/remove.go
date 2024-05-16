package cmd

import (
	"fmt"
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
		Run:     internal.RunWrapper(wrapRemoveSubmodule, internal.ArgExactN(1)),
	}

	rmCmd.Flags().BoolP("delete", "d", false, "delete existing directory")
	rmCmd.Flags().BoolP("force", "f", false, "force delete existing directory")

	return rmCmd
}

func wrapRemoveSubmodule(cmd *cobra.Command, args []string) {
	deleteDirectory, _ := cmd.Flags().GetBool("delete")
	forceDelete, _ := cmd.Flags().GetBool("force")
	err := removeSubmodule(models.Path(args[0]), deleteDirectory, forceDelete)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func removeSubmodule(p models.Path, deleteDirectory bool, forceDelete bool) error {

	// acquire lock
	lockFile, err := internal.ErrorWrappedLockAcquiringAtProjectRootFromCwd()
	defer func() {
		err := internal.ErrorWrappedLockReleasing(lockFile)
		if err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		return err
	}

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
