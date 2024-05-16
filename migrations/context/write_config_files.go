package context

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/models"
)

type WriteConfigFiles struct {
	Context *models.NestContext
}

func (m WriteConfigFiles) Migrate() error {
	r, err := internal.WriteProjectConfigFiles(*m.Context)
	if err == nil {
		if r.ConfigWriteError != nil {
			err = r.ConfigWriteError
		} else if r.GitExcludeWriteError != nil {
			err = r.GitExcludeWriteError
		}
	}

	if err != nil {
		return fmt.Errorf("error writing configuration files: %w", err)
	}

	return nil
}
