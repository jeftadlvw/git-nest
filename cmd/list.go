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
		existStr := ""
		submoduleExists := submodulesExist[index]

		// translate flag, payload and error
		switch submoduleExists.Status {
		case internal.SUBMODULE_EXISTS_OK:
			existStr = "ok"
		case internal.SUBMODULE_EXISTS_ERR_NO_EXIST:
			existStr = "no exist"
		case internal.SUBMODULE_EXISTS_ERR_FILE:
			existStr = "error: path is a file"
		case internal.SUBMODULE_EXISTS_ERR_NO_GIT:
			existStr = "error: git not installed"
		case internal.SUBMODULE_EXISTS_ERR_REMOTE:
			existStr = "error: unequal remote urls: " + submoduleExists.Payload
		case internal.SUBMODULE_EXISTS_ERR_HEAD:
			if submoduleExists.Error != nil {
				existStr = "error: unable to fetch HEAD: " + submoduleExists.Error.Error()
			} else {
				existStr = "error: unequal ref HEADs: " + submoduleExists.Payload
			}
		}

		_, _ = fmt.Fprintf(tabWriter, "%d\t%s\t%s\t%s\t%s", index+1, submodule.Path, submodule.Url.String(), submodule.Ref, existStr)
	}
	_ = tabWriter.Flush()

	fmt.Println(buffer.String())
}
