package cmd

import (
	"github.com/jeftadlvw/git-nest/cmd/internal"
	application_internal "github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/spf13/cobra"
)

func Execute() (int, error) {

	var err error = nil

	var rootCmd = &cobra.Command{
		Use:     constants.ApplicationName,
		Version: constants.Version(),
		Short:   "Nest external repositories into your project without git knowing.",
		Long: `git-nest is a git command line extension for nesting external repositories
in your project without your parent repository noticing, using native features
and configurations files.`,
		RunE:          internal.PrintUsage,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	configureRootCommand(rootCmd)

	// ensure application mutex
	lf, err := internal.GetApplicationMutex()
	if err != nil {
		return -1, err
	}
	application_internal.AddCleanup(lf.Release)

	// execute command handler
	err = rootCmd.Execute()
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func configureRootCommand(rootCmd *cobra.Command) {
	// add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(createInfoCmd())
	rootCmd.AddCommand(createAddCmd())
	rootCmd.AddCommand(createRemoveCommand())
	rootCmd.AddCommand(createListCmd())
	rootCmd.AddCommand(createVerifyCmd())
	rootCmd.AddCommand(syncCmd)

	// miscellaneous configuration
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
