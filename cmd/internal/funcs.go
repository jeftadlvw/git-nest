package internal

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	application_internal "github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/spf13/cobra"
	"os"
)

/*
PrintUsage is a wrapper function around the default cobra.Command Usage() function.
*/
func PrintUsage(cmd *cobra.Command, args []string) error {
	_ = cmd.Usage()
	return nil
}

/*
RunWrapper wraps the function set for the 'Run' field in a cobra.Command.
It takes a runner function and an argument count validation function. If the latter
is not nil, it is executed first and checked for returned errors. If no errors
were returned, the runner function is executed.
*/
func RunWrapper(run func(cmd *cobra.Command, args []string) error, validateArgCount ...func(c int) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if validateArgCount != nil {
			for _, validationFunc := range validateArgCount {
				err := validationFunc(len(args))
				return fmt.Errorf("argument count error: %s\n", err)
			}
		}

		return run(cmd, args)
	}
}

/*
GetProjectRootFromCwd returns the project root directory, starting from the current directory.
*/
func GetProjectRootFromCwd() (models.Path, error) {
	cwdStr, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current working directory: %w", err)
	}
	cwd := models.Path(cwdStr)
	projectRoot, err := internal.FindProjectRoot(cwd)
	if err != nil {
		projectRoot = cwd
	}

	return projectRoot, nil
}

/*
ErrorWrappedEvaluateContext is a wrapper for the cmd package to remove repetitive boilerplate code.
It returns the evaluated context or a preformatted error.
*/
func ErrorWrappedEvaluateContext() (models.NestContext, error) {
	context, err := internal.EvaluateContext()
	if err != nil {
		return models.NestContext{}, fmt.Errorf("internal context error: %w.\nPlease fix any configuration errors to proceed", err)
	}

	return context, nil
}

/*
GetApplicationMutex is a wrapper function to get a git-nest lockfile at the project root directory.
*/
func GetApplicationMutex() (internal.LockFile, error) {
	// get project root
	projectRoot, err := GetProjectRootFromCwd()
	if err != nil {
		return internal.LockFile{}, fmt.Errorf("could not get project root: %w", err)
	}

	// create lockfile
	lf, err := application_internal.CreateLockFile(projectRoot.SJoin(constants.LockFileName))
	if err != nil {
		infoText := "Another git-nest process might already be running in this project."
		return internal.LockFile{}, fmt.Errorf("unable to acquire lockfile: %s\n%s", err, infoText)
	}

	return lf, nil
}
