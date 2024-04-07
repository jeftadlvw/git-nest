package internal

import (
	"github.com/spf13/cobra"
)

func PrintUsage(cmd *cobra.Command, args []string) {
	_ = cmd.Usage()
}
