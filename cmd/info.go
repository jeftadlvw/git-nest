package cmd

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/spf13/cobra"
	"runtime"
	"strings"
	"time"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Various information (useful for debugging)",
	Run: func(cmd *cobra.Command, args []string) {
		printCompilationInformation()
	},
}

func printCompilationInformation() {
	// beautify compilation time output
	compilationTime := "unknown"
	if constants.CompilationTimestamp() != -1 {
		layout := "Mon Jan 02 15:04:05 2006"
		compilationTime = time.Unix(int64(constants.CompilationTimestamp()), 0).Format(layout)
	}

	fmt.Printf("%s info dump:\n", constants.ApplicationName)
	fmt.Printf("  Version:\t%s\n", constants.Version())
	fmt.Printf("  Git ref:\t%s\n", constants.RefHash())
	fmt.Printf("  Runtime:\t%s\n", runtime.Version())
	fmt.Printf("  Built:\t%s\n", compilationTime)
	fmt.Printf("  OS/Arch:\t%s/%s\n", runtime.GOOS, runtime.GOARCH)

	// estimate if binary is local dev build
	if strings.HasPrefix(constants.Version(), "[") || constants.RefHash() == "unset" || constants.CompilationTimestamp() == -1 {
		fmt.Printf("\nThis binary is most likely a local development build.\n")
	}
}
