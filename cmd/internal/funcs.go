package internal

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/spf13/cobra"
	"os"
)

/*
PrintUsage is a wrapper function around the default cobra.Command Usage() function.
*/
func PrintUsage(cmd *cobra.Command, args []string) {
	_ = cmd.Usage()
}

/*
RunWrapper wraps the function set for the 'Run' field in a cobra.Command.
It takes a runner function and an argument count validation function. If the latter
is not nil, it is executed first and checked for returned errors. If no errors
were returned, the runner function is executed.
*/
func RunWrapper(run func(cmd *cobra.Command, args []string), validateArgCount ...func(c int) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if validateArgCount != nil {
			for _, validationFunc := range validateArgCount {
				err := validationFunc(len(args))
				if err != nil {
					fmt.Printf("fatal: argument count error: %s\n", err)
					return
				}
			}
		}

		run(cmd, args)
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
ErrorWrappedLockAcquiringAtProjectRootFromCwd is a wrapper for internal.AcquireLockFile to remove boilerplate code.
It returns an error if the lockfile at the project root starting from the current working directory
could not be acquired.
*/
func ErrorWrappedLockAcquiringAtProjectRootFromCwd() (*os.File, error) {
	projectRoot, err := GetProjectRootFromCwd()
	if err != nil {
		return nil, fmt.Errorf("could not get project root: %w", err)
	}

	lockFile, err := internal.AcquireLockFile(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("error while acquiring lock: %w", err)
	}

	return lockFile, nil
}

/*
ErrorWrappedLockReleasing is a wrapper for internal.ReleaseLockFile to remove boilerplate code.
It returns an error if the lockfile could not be released.
*/
func ErrorWrappedLockReleasing(f *os.File) error {
	err := internal.ReleaseLockFile(f)
	if err != nil {
		return fmt.Errorf("could not release lock: %w", err)
	}

	return nil
}
