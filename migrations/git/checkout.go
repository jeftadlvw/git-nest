package git

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"strings"
)

type Checkout struct {
	Path models.Path
	Ref  string
}

func (m Checkout) Migrate() error {

	if !m.Path.Exists() {
		return fmt.Errorf("%s does not exist", m.Path)
	}

	if !m.Path.IsDir() {
		return fmt.Errorf("%s is not a directory", m.Path)
	}

	ref := strings.TrimSpace(m.Ref)
	if ref == "" {
		return fmt.Errorf("ref cannot be blank")
	}

	err := utils.ChangeGitHead(m.Path, ref)
	if err != nil {
		return fmt.Errorf("error while changing ref: %s", err)
	}

	return nil
}
