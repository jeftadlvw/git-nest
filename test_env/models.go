package test_env

import (
	"github.com/jeftadlvw/git-nest/models"
	"os"
)

type TestEnv struct {
	Dir models.Path
}

/*
	func (t *TestEnv) AddNestedModules(modules ...models.Submodule) {

	}
*/

func (t *TestEnv) Destroy() {
	_ = os.RemoveAll(t.Dir.String())
}
