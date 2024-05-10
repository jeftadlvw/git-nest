package context

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/models"
)

type WriteConfigFiles struct {
	Context models.NestContext
}

func (m WriteConfigFiles) Migrate() error {
	var err error

	_, _, err1, err2 := internal.WriteProjectConfigFiles(m.Context)
	if err1 != nil {
		err = err1
	} else if err2 != nil {
		err = err2
	}

	if err != nil {
		return fmt.Errorf("error writing configuration files: %w", err)
	}

	return nil
}
