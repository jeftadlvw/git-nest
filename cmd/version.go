package cmd

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Print %s version", constants.ApplicationName),
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	fmt.Println(constants.Version())
}
