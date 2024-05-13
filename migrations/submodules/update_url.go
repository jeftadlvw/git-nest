package submodules

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
)

type UpdateUrl struct {
	Submodule *models.Submodule
	Url       urls.HttpUrl
}

func (m UpdateUrl) Migrate() error {
	if m.Submodule == nil {
		return errors.New("migration contained nil submodule")
	}

	err := m.Url.Validate()
	if err != nil {
		return fmt.Errorf("validation error for url: %w", err)
	}

	m.Submodule.Url = &m.Url
	return nil
}
