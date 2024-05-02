package internal

import (
	"fmt"
	"github.com/spf13/cobra"
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
