package internal

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"os"
)

type LockFile struct {
	file *os.File
}

func (l *LockFile) Release() error {

	// close the file
	err := l.file.Close()
	if err != nil {
		return fmt.Errorf("error closing lockfile: %w", err)
	}

	// remove the file
	err = os.Remove(l.file.Name())
	if err != nil {
		return fmt.Errorf("error removing lockfile: %w", err)
	}

	return nil
}

func CreateLockFile(p models.Path) (LockFile, error) {
	if p.IsDir() {
		return LockFile{}, errors.New("lockfile is directory")
	}

	// if file exist, leave it alone and act like it's locked
	if p.IsFile() {
		return LockFile{}, fmt.Errorf("lockfile already exists")
	}

	// open a new file and prevent creation if it already exists
	f, err := os.OpenFile(p.String(), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return LockFile{}, fmt.Errorf("error creating lockfile: %w", err)
	}

	// create LockFile struct
	lockFile := LockFile{file: f}

	return lockFile, nil
}
