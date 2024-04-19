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

	// validate each submodule
	for index, submodule := range c.Submodules {
		err = submodule.Validate()
		if err != nil {
			return fmt.Errorf("error at submodule index %d: %w", index, err)
		}
	}

	// check for duplicates
	var (
		identifierSet = mapset.NewSet[string]()
		pathSet       = mapset.NewSet[string]()
		remoteUrlSet  = mapset.NewSet[string]()
	)

	for _, submodule := range c.Submodules {
		var added bool

		// submodules may not escape project root (by having / or ../ as prefix)
		if strings.HasPrefix(submodule.Path.String(), string(filepath.Separator)) {
			return fmt.Errorf("submodule path is relative to project root")
		}

		if strings.Contains(submodule.Path.String(), string(filepath.Separator)) {
			return fmt.Errorf("submodule path escapes project root (%s)", submodule.Path)
		}

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
		if !added && !c.Config.AllowDuplicateOrigins {
			return fmt.Errorf("submodule origin urls %s defined multiple times", submoduleRemoteUrl)
		}
	}

	return nil
}
