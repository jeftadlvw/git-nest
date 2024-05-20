package internal

import (
	"fmt"
	"os"
	"sync"
)

var cleanupMutex sync.Mutex
var cleanupStack []func() error

/*
Cleanup is responsible for cleaning up any left-over files and things.
Gets called by the main function on exiting.
*/
func Cleanup() {
	cleanupMutex.Lock()

	var err error
	for _, f := range cleanupStack {
		err = f()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error while cleaning up: %s", err)
		}
	}

	cleanupMutex.Unlock()
}

/*
AddCleanup adds a function to the cleanup stack.
*/
func AddCleanup(f func() error) {
	cleanupMutex.Lock()
	cleanupStack = append(cleanupStack, f)
	cleanupMutex.Unlock()
}
