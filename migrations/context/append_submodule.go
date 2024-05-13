package context

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
)

type AppendSubmodule struct {
	Context   *models.NestContext
	Submodule models.Submodule
}

func (m AppendSubmodule) Migrate() error {
	if m.Context == nil {
		return fmt.Errorf("migration contained nil context")
	}

	// validate submodule
	err := m.Submodule.Validate()
	if err != nil {
		return fmt.Errorf("submodule validation error: %w", err)
	}

	// check for duplicates
	newSubmoduleSlice := append(m.Context.Config.Submodules, m.Submodule)
	err = models.CheckForDuplicateSubmodules(m.Context.Config.Config.AllowDuplicateOrigins, newSubmoduleSlice...)
	if err != nil {
		return fmt.Errorf("duplicate checking: %s", err)
	}

	// set submodule slice
	m.Context.Config.Submodules = newSubmoduleSlice
	return nil
}
