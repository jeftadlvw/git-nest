package fs

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"os"
)

type DeleteDirectory struct {
	Path   models.Path
	DryRun bool
}

func (m DeleteDirectory) Migrate() error {
	if m.Path.Empty() {
		return errors.New("path is empty")
	}

	if m.Path.AtRoot() {
		return errors.New("cannot delete system root directory")
	}

	if !m.DryRun {
		err := os.RemoveAll(m.Path.String())
		if err != nil {
			return fmt.Errorf("could not delete directory: %w", err)
		}
	}

	return nil
}
