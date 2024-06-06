package cmd

import (
	"fmt"
	cmdInternal "github.com/jeftadlvw/git-nest/cmd/internal"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/spf13/cobra"
)

func createVerifyCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "verify",
		Aliases: []string{"v"},
		Short:   "Verify configuration and nested modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyConfigAndSubmodules()
		},
	}

	return listCmd
}

func verifyConfigAndSubmodules() error {
	// read context
	context, err := cmdInternal.ErrorWrappedEvaluateContext()
	if err != nil {
		return err
	}

	submodulesExist := internal.SubmodulesExist(context.Config.Submodules, context.ProjectRoot)

	for index := range len(context.Config.Submodules) {
		submoduleExists := submodulesExist[index]
		existStr, err := internal.FmtSubmoduleExistOutput(submoduleExists.Status, submoduleExists.Payload, submoduleExists.Error)
		if err != nil {
			existStr = "internal error: " + err.Error()
		}

		if internal.SubmoduleStatusValid(submoduleExists.Status) {
			fmt.Printf("error for nested module at index %d: %s\n", index, existStr)
		}
	}

	return nil
}
