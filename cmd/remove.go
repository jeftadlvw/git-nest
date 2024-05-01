package cmd

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/actions"
	cmdInternal "github.com/jeftadlvw/git-nest/cmd/internal"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/spf13/cobra"
)

func createRemoveCommand() *cobra.Command {
	var infoCmd = &cobra.Command{
		Use:     "remove path | url",
		Aliases: []string{"rm"},
		Short:   "Remove a submodule from this project",
		Run:     cmdInternal.RunWrapper(wrapRemoveSubmodule, cmdInternal.ArgExactN(1)),
	}

	return infoCmd
}

func wrapRemoveSubmodule(args []string) {
	err := removeSubmodule(args[0])
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func removeSubmodule(s string) error {

	// read context
	context, err := internal.EvaluateContext()
	if err != nil {
		return fmt.Errorf("internal context error: %w.\nPlease fix any configuration errors to proceed", err)
	}

	err = actions.RemoveSubmoduleFromContext(&context, s)
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
