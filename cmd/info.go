package cmd

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/jeftadlvw/git-nest/utils"
	"github.com/spf13/cobra"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func createInfoCmd() *cobra.Command {
	var infoCmd = &cobra.Command{
		Use:     "info",
		Aliases: []string{"i"},
		Short:   "Print various debug information",
		Run: func(cmd *cobra.Command, args []string) {
			redact, _ := cmd.Flags().GetBool("redact")
			printDebugInformation(redact)
		},
	}
	infoCmd.Flags().BoolP("redact", "r", false, "hide personal info")

	return infoCmd
}

func printDebugInformation(redact bool) {
	context, err := internal.EvaluateContext()
	if err != nil {
		fmt.Println(err)
		return
	}

	workingDir := context.WorkingDirectory.String()
	rootDir := context.ProjectRoot.String()
	repositoryRoot := context.GitRepositoryRoot.String()
	validNestedModules := 0

	// beautify compilation time output
	compilationTime := "unknown"
	if constants.CompilationTimestamp() != -1 {
		layout := "Mon Jan 02 15:04:05 2006"
		compilationTime = time.Unix(int64(constants.CompilationTimestamp()), 0).Format(layout)
	}

	configurationFileString := ""
	if context.ConfigFileExists {
		configurationFileString = string(context.ConfigFile)
		if redact {
			configurationFileString, err = filepath.Rel(workingDir, configurationFileString)
			if err != nil {
				configurationFileString = "error"
			} else {
				configurationFileString = "." + string(filepath.Separator) + configurationFileString
			}
		}
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

	if redact {
		rootDir, err = filepath.Rel(workingDir, rootDir)
		if err != nil {
			rootDir = "error"
		}

		repositoryRoot, err = filepath.Rel(workingDir, repositoryRoot)
		if err != nil {
			repositoryRoot = "error"
		}

		workingDir = "."
	}

	for _, validNestedModule := range internal.SubmodulesExist(context.Config.Submodules, context.ProjectRoot) {
		if validNestedModule.Status == internal.SUBMODULE_EXISTS_OK {
			validNestedModules++
		}
	}

	infoMap := []utils.Node{
		{"Binary", []utils.Node{
			{"Version", constants.Version()},
			{"Git ref", constants.Ref()},
			{"Runtime", runtime.Version()},
			{"Build", compilationTime},
			{"OS/Arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)},
		}},
		{"Context", []utils.Node{
			{"Working directory", workingDir},
			{"Root directory", rootDir},
			{"Configuration file", configurationFileString},
			{"Git installed", gitInstalledString},
			{"Git repository", context.IsGitRepository},
			{"Repository root", repositoryRoot},
		}},

		{"Valid modules", fmt.Sprintf("%d/%d", validNestedModules, len(context.Config.Submodules))},
	}

	fmt.Println(utils.FmtTree(utils.FmtTreeConfig{
		Indent:             "   ",
		NewLinesAtTopLevel: true,
	}, infoMap...))

	// estimate if binary is local dev build
	if strings.HasPrefix(constants.Version(), "[") || constants.Ref() == "unset" || constants.CompilationTimestamp() == -1 {
		fmt.Printf("\nThis binary is most likely a local development built.\n")
	}
}
