package main

import (
	"github.com/jeftadlvw/git-nest/cmd"
	"github.com/jeftadlvw/git-nest/internal"
	"os"
	"os/signal"
	"sync"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)

	var execOnce sync.Once

	go internal.HandleOsTerminationSignals(c, func() {
		execOnce.Do(cmd.MockCleanupFunc)
	})

	exitCode := cmd.MockRunFun()
	execOnce.Do(cmd.MockCleanupFunc)
	os.Exit(exitCode)
}
