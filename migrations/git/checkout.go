package git

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
)

type Checkout struct {
	Path models.Path
	Ref  string
}

func (m Checkout) Migrate() error {
	err := utils.GitCheckout(m.Path, m.Ref)
	if err != nil {
		return fmt.Errorf("error while changing ref: %s", err)
	}

	return nil
}
