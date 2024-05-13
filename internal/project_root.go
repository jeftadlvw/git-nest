package internal

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/jeftadlvw/git-nest/models"
)

/*
FindProjectRoot takes a directory and searches this and upper directories for a
`nestmodules.toml` configuration file in its root- or a `.config` subdirectory. The first directory that
fits these conditions is chosen as the project root. Returns an error if no root directory found.

FindProjectRoot will not check the config file's integrity. Its existence is enough.
*/
func FindProjectRoot(p models.Path) (models.Path, error) {

	if !p.Exists() {
		return "", fmt.Errorf("%s does not exist", p)
	}

	// if p is a file, set p to the file's directory
	if p.IsFile() {
		p = p.Parent()
	}

	for {
		if p.BContains(constants.ConfigFileName) || p.BContains(constants.ConfigSubDirFileName) {
			break
		}

		if p.EmptyOrAtRoot() {
			return "", fmt.Errorf("%s reached root of original path", p)
		}
		p = p.Parent()
	}

	return p, nil
}
