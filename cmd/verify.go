package cmd

import (
	"bytes"
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/spf13/cobra"
	"text/tabwriter"
)

func createVerifyCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "verify",
		Aliases: []string{"v"},
		Short:   "Verify configuration and nested modules",
		Run: func(cmd *cobra.Command, args []string) {
			verifyConfigAndSubmodules()
		},
	}

	return listCmd
}

func verifyConfigAndSubmodules() {
	context, err := internal.EvaluateContext()
	if err != nil {
		fmt.Println("context error:", err)
		return
	}

	submodulesExist := internal.SubmodulesExist(context.Config.Submodules, context.ProjectRoot)

	buffer := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(buffer, 5, 0, 0, '.', tabwriter.TabIndent)

	_, _ = fmt.Fprintf(tabWriter, "i\tpath\torigin\tref\tstatus")
	for index := range len(context.Config.Submodules) {
		submoduleExists := submodulesExist[index]
		existStr, err := internal.FmtSubmoduleExistOutput(submoduleExists.Status, submoduleExists.Payload, submoduleExists.Error)
		if err != nil {
			existStr = "internal error: " + err.Error()
		}

		if existStr != "ok" {
			fmt.Printf("error for nested module at index %d: %s\n", index, existStr)
		}
	}
}
