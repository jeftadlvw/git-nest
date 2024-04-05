package internal

import (
	"fmt"
	"os"
	"syscall"
)

func HandleOsTerminationSignals(c <-chan os.Signal, callback func()) {
	for {
		// waiting for signal to arrive
		sig := <-c

		// call callback
		callback()

		// print exit information to stderr
		_, _ = fmt.Fprintf(os.Stderr, "\n%s\n", fmt.Sprintf("Exiting because of %s signal.", sig))

		// type assert signal to be of type syscall.Signal
		sysSig, ok := sig.(syscall.Signal)

		// apply value of type cast to exitCode
		exitCode := 1
		if ok {
			exitCode = int(sysSig)
		}

		// terminate application
		os.Exit(exitCode)
	}
}
