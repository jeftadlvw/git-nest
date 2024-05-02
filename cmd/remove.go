package cmd

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/actions"
	cmdInternal "github.com/jeftadlvw/git-nest/cmd/internal"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/spf13/cobra"
)

func createRemoveCommand() *cobra.Command {
	var rmCmd = &cobra.Command{
		Use:     "remove [path]",
		Aliases: []string{"rm"},
		Short:   "Remove a submodule from this project",
		Run:     cmdInternal.RunWrapper(wrapRemoveSubmodule, cmdInternal.ArgExactN(1)),
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

	// read context
	context, err := internal.EvaluateContext()
	if err != nil {
		return fmt.Errorf("internal context error: %w.\nPlease fix any configuration errors to proceed", err)
	}

	err = actions.RemoveSubmoduleFromContext(&context, p, deleteDirectory, forceDelete)
	if err != nil {
		return err
	}

	// write configuration files
	_, _, err1, err2 := internal.WriteProjectConfigFiles(context)
	if err1 != nil {
		return fmt.Errorf("error writing configuration files: %w", err1)
	} else if err2 != nil {
		return fmt.Errorf("error writing configuration files: %w", err2)
	}

	return nil
}
