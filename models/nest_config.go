package models

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"os"
)

type NestConfig struct {
	Config     Config
	Submodules []Submodule
}

func (c NestConfig) Validate() error {
	err := c.Config.Validate()
	if err != nil {
		return err
	}

	// validate each submodule
	for _, submodule := range c.Submodules {
		err = submodule.Validate()
		if err != nil {
			return fmt.Errorf("error at submodule %s@%s: %w", submodule.Url.String(), submodule.Ref, err)
		}
	}

	// check for duplicates
	var (
		identifierSet       = mapset.NewSet[string]()
		pathSet             = mapset.NewSet[string]()
		remoteUrlSet        = mapset.NewSet[string]()
		remoteIdentifierSet = mapset.NewSet[string]()

		duplicateOriginUrlSet = mapset.NewSet[string]()
		duplicateOriginSet    = mapset.NewSet[string]()
	)

	for _, submodule := range c.Submodules {
		var added bool

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

		// check if Submodules have duplicate remote origin urls
		submoduleRemoteUrl := submodule.Url.String()
		added = remoteUrlSet.Add(submoduleRemoteUrl)
		if !added && !c.Config.AllowDuplicateOrigins && !duplicateOriginUrlSet.Contains(submoduleRemoteUrl) {
			_, _ = fmt.Fprintf(os.Stderr, "submodule origin url %s defined multiple times", submoduleRemoteUrl)
			duplicateOriginUrlSet.Add(submoduleRemoteUrl)
		}

		// check if duplicate submodule is nested multiple times (but different directories)
		submoduleRemoteIdentifier := submodule.Identifier()
		added = remoteIdentifierSet.Add(submoduleRemoteIdentifier)
		if !added && !c.Config.AllowDuplicateOriginRefs && !duplicateOriginSet.Contains(submoduleRemoteIdentifier) {
			_, _ = fmt.Fprintf(os.Stderr, "submodule origin %s defined multiple times", submoduleRemoteIdentifier)
			duplicateOriginSet.Add(submoduleRemoteIdentifier)
		}
	}

	return nil
}
