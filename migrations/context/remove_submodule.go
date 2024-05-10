package context

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"slices"
)

type RemoveSubmodule struct {
	Context        *models.NestContext
	SubmoduleIndex int
}

func (m RemoveSubmodule) Migrate() error {
	if m.Context == nil {
		return fmt.Errorf("migration contained nil context")
	}

	if m.SubmoduleIndex < 0 || m.SubmoduleIndex >= len(m.Context.Config.Submodules) {
		return fmt.Errorf("submodule index out of range (%d, but can be [0, %d])", m.SubmoduleIndex, len(m.Context.Config.Submodules)-1)
	}

	// remove submodule from submodule slice
	m.Context.Config.Submodules = slices.Delete(m.Context.Config.Submodules, m.SubmoduleIndex, m.SubmoduleIndex+1)

	return nil
}
