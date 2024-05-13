package submodules

import (
	"errors"
	"github.com/jeftadlvw/git-nest/models"
	"strings"
)

type UpdateRef struct {
	Submodule *models.Submodule
	Ref       string
}

func (m UpdateRef) Migrate() error {
	if m.Submodule == nil {
		return errors.New("migration contained nil submodule")
	}

	m.Ref = strings.TrimSpace(m.Ref)
	if m.Ref == "" {
		return errors.New("migration contained empty ref.\nThis functionality is currently unsupported")
	}

	m.Submodule.Ref = m.Ref
	return nil
}
