package cmd

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/jeftadlvw/git-nest/utils"
	"github.com/spf13/cobra"
	"runtime"
	"strings"
	"time"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print various debug information",
	Run: func(cmd *cobra.Command, args []string) {
		printDebugInformation()
	},
}

func printDebugInformation() {
	context, err := internal.EvaluateContext()
	if err != nil {
		fmt.Println(err)
		return
	}

	// beautify compilation time output
	compilationTime := "unknown"
	if constants.CompilationTimestamp() != -1 {
		layout := "Mon Jan 02 15:04:05 2006"
		compilationTime = time.Unix(int64(constants.CompilationTimestamp()), 0).Format(layout)
	}

	configurationFileString := ""
	if context.ConfigFileExists {
		configurationFileString = string(context.ConfigFile)
	} else {
		configurationFileString = "none"
	}

	gitInstalledString := fmt.Sprintf("%t", context.IsGitInstalled)
	if context.IsGitInstalled {
		gitVersion, err := utils.GetGitVersion()
		if err == nil {
			gitInstalledString = fmt.Sprintf("%s; %s", gitInstalledString, gitVersion)
		}
	}

	infoMap := []utils.Node{
		{"Binary", []utils.Node{
			{"Version", constants.Version()},
			{"Git ref", constants.RefHash()},
			{"Runtime", runtime.Version()},
			{"Build", compilationTime},
			{"OS/Arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)},
		}},
		{"Context", []utils.Node{
			{"Working directory", context.WorkingDirectory},
			{"Root directory", context.ProjectRoot},
			{"Configuration file", configurationFileString},
			{"Git installed", gitInstalledString},
			{"Git project", context.IsGitProject},
		}},
	}

	fmt.Printf(utils.FmtTree("", infoMap...))

	// estimate if binary is local dev build
	if strings.HasPrefix(constants.Version(), "[") || constants.RefHash() == "unset" || constants.CompilationTimestamp() == -1 {
		fmt.Printf("This binary is most likely a local development built.\n")
	}
}
