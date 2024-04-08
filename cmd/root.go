package cmd

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/cmd/internal"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     constants.ApplicationName,
	Version: constants.Version(),
	Short:   "Nest external repositories into your project without git knowing.",
	Long: `git-nest is a git command line extension for nesting external repositories
in your project without your parent repository noticing, using native features
and configurations files.`,
	Run: internal.PrintUsage,
}

func Execute() int {

	configureRootCommand(rootCmd)

	err := rootCmd.Execute()
	var exitCode = 0
	if err != nil {
		fmt.Println(err)
		exitCode = 1
	}

	return exitCode
}

func configureRootCommand(rootCmd *cobra.Command) {
	// add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(infoCmd)

	// miscellaneous configuration
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}