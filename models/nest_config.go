package models

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"path/filepath"
	"strings"
)

/*
NestConfig represents a higher-order configuration for git-nest.
*/
type NestConfig struct {
	/*
		Config contains every configuration flag.
	*/
	Config Config `toml:"config"`

	/*
		Submodule contains all Submodule for this configuration.
	*/
	Submodules []Submodule `toml:"submodule"`
}

/*
Validate performs validation on this NestConfig.
*/
func (c NestConfig) Validate() error {
	err := c.Config.Validate()
	if err != nil {
		return err
	}

	// validate each submodule and check for duplicates
	for index, submodule := range c.Submodules {
		err = submodule.Validate()
		if err != nil {
			return fmt.Errorf("error at submodule index %d: %w", index, err)
		}

		// submodules may not escape project root (by having / or ../ as prefix)
		if strings.HasPrefix(submodule.Path.String(), string(filepath.Separator)) {
			return fmt.Errorf("submodule path must be relative to project root")
		}

		if strings.HasPrefix(submodule.Path.String(), "..") {
			return fmt.Errorf("submodule path escapes project root (%s)", submodule.Path)
		}
	}

	err = CheckForDuplicateSubmodules(c.Config.AllowDuplicateOrigins, c.Submodules...)
	if err != nil {
		return err
	}

	return nil
}

/*
CheckForDuplicateSubmodules syntactically checks if any duplicate submodules exist within a slice of Submodule.
*/
func CheckForDuplicateSubmodules(allowDuplicateOrigins bool, submodules ...Submodule) error {
	var added bool

	var (
		identifierSet = mapset.NewSet[string]()
		pathSet       = mapset.NewSet[string]()
		remoteUrlSet  = mapset.NewSet[string]()
	)

	for _, submodule := range submodules {
		// check for 100% duplicates
		added = identifierSet.Add(submodule.Identifier())
		if !added {
			return fmt.Errorf("submodule %s defined multiple times", submodule.Identifier())
		}

		// a directory cannot be used twice
		added = pathSet.Add(submodule.Path.String())
		if !added {
			return fmt.Errorf("submodule directory %s used multiple times", submodule.Path)
		}

		// check if submodules have duplicate remote origin urls
		// it's enough to check for the url, because if AllowDuplicateOrigins is set to true,
		// duplicate refs are automatically allowed too (no ref is a ref to something default!)
		submoduleRemoteUrl := submodule.Url.String()
		added = remoteUrlSet.Add(submoduleRemoteUrl)
		if !added && !allowDuplicateOrigins {
			return fmt.Errorf("submodule origin urls %s defined multiple times", submoduleRemoteUrl)
		}
	}

	return nil
}
