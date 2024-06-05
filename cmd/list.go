package cmd

import (
	"bytes"
	"fmt"
	cmdInternal "github.com/jeftadlvw/git-nest/cmd/internal"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/spf13/cobra"
	"text/tabwriter"
)

func createListCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List nested modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printSubmodules()
		},
	}

	return listCmd
}

func printSubmodules() error {
	// read context
	context, err := cmdInternal.ErrorWrappedEvaluateContext()
	if err != nil {
		return err
	}

	if len(context.Config.Submodules) == 0 {
		fmt.Println("no nested modules defined in current context")
		return nil
	}

	submodulesExist := internal.SubmodulesExist(context.Config.Submodules, context.ProjectRoot)

	buffer := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(buffer, 5, 0, 1, ' ', tabwriter.TabIndent)

	_, _ = fmt.Fprintf(tabWriter, "i\tpath\torigin\tref\tstatus\n")
	for index, submodule := range context.Config.Submodules {
		submoduleExists := submodulesExist[index]
		existStr, err := internal.FmtSubmoduleExistOutput(submoduleExists.Status, submoduleExists.Payload, submoduleExists.Error)
		if err != nil {
			existStr = "internal error: " + err.Error()
		}

		_, _ = fmt.Fprintf(tabWriter, "%d\t%s\t%s\t%s\t%s\n", index+1, submodule.Path, submodule.Url.String(), submodule.Ref, existStr)
	}
	_ = tabWriter.Flush()

	fmt.Println(buffer.String())
	return nil
}
