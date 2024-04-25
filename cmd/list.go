package cmd

import (
	"bytes"
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/spf13/cobra"
	"text/tabwriter"
)

func createListCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List nested modules",
		Run: func(cmd *cobra.Command, args []string) {
			printSubmodules()
		},
	}

	return listCmd
}

func printSubmodules() {
	context, err := internal.EvaluateContext()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(context.Config.Submodules) == 0 {
		fmt.Println("no nested modules defined in current context")
		return
	}

	submodulesExist := internal.SubmodulesExist(context.Config.Submodules, context.ProjectRoot)

	buffer := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(buffer, 5, 0, 0, '.', tabwriter.TabIndent)

	_, _ = fmt.Fprintf(tabWriter, "i\tpath\torigin\tref\tstatus")
	for index, submodule := range context.Config.Submodules {
		submoduleExists := submodulesExist[index]
		existStr, err := internal.FmtSubmoduleExistOutput(submoduleExists.Status, submoduleExists.Payload, submoduleExists.Error)
		if err != nil {
			existStr = "internal error: " + err.Error()
		}

		_, _ = fmt.Fprintf(tabWriter, "%d\t%s\t%s\t%s\t%s", index+1, submodule.Path, submodule.Url.String(), submodule.Ref, existStr)
	}
	_ = tabWriter.Flush()

	fmt.Println(buffer.String())
}
