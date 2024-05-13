package git

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
)

type Clone struct {
	Url          interfaces.Url
	Path         models.Path
	CloneDirName string
}

func (m Clone) Migrate() error {
	err := utils.CloneGitRepository(m.Url.String(), m.Path, m.CloneDirName)
	if err != nil {
		return fmt.Errorf("error while cloning into %s: %s", m.Path.SJoin(m.CloneDirName), err)
	}

	return nil
}
