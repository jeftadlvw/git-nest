package internal

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/jeftadlvw/git-nest/models"
	"os"
)

/*
AcquireLockFile tries to acquire a context-specific lockfile.
It returns if the lockfile was acquired successfully.
*/
func AcquireLockFile(p models.Path) (*os.File, error) {
	lockFilePath := p.SJoin(constants.LockFileName)

	if lockFilePath.IsDir() {
		return nil, errors.New("lock file is directory")
	}

	infoText := "----\nAnother git-nest process might already be running in this project."

	// if file exist, leave it alone and act like it's locked
	if lockFilePath.IsFile() {
		return nil, fmt.Errorf("file is locked;\n%s", infoText)
	}

	// open a new file and prevent creation if it already exists
	file, err := os.OpenFile(lockFilePath.String(), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return nil, fmt.Errorf("error creating lock file:\n%w;\n%s", err, infoText)
	}

	return file, nil
}

/*
ReleaseLockFile releases the context-specific lockfile.
*/
func ReleaseLockFile(f *os.File) error {
	err := f.Close()
	if err != nil {
		return fmt.Errorf("error closing lock file: %w", err)
	}

	err = os.Remove(f.Name())
	if err != nil {
		return fmt.Errorf("error removing lock file: %w", err)
	}

	return nil
}
