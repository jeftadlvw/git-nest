package internal

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"path/filepath"
	"strings"
)

/*
PathContainsUp returns whether the passed models.Path contains any "..".
*/
func PathContainsUp(p models.Path) bool {
	return strings.Contains(p.String(), "..")
}

/*
PathOutsideRoot returns whether a given path is not located within a root path.
This check is purely lexical.
*/
func PathOutsideRoot(root models.Path, p models.Path) bool {
	relative, err := root.Relative(p)
	if err != nil {
		return true
	}

	return PathContainsUp(relative)
}

/*
PathRelativeToRootButOtherOriginIfNotAbs is a somewhat complicated function.
It calculates the relative path between root and p if p is an absolute path.
But in case p is not absolute, it is first joined on the origin path.

TODO: improve this function name.
*/
func PathRelativeToRootButOtherOriginIfNotAbs(root models.Path, origin models.Path, p models.Path) (models.Path, error) {
	var (
		err            error
		comparePath    models.Path
		relativeToRoot models.Path
	)

	if filepath.IsAbs(p.String()) {
		comparePath = p
	} else {
		comparePath = origin.Join(p)
	}

	relativeToRoot, err = root.Relative(comparePath)
	if err != nil {
		return "", fmt.Errorf("internal error: could not find relative to project root: %w", err)
	}

	return relativeToRoot, nil
}
