package test_env

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/utils"
	"strings"
)

func CreateTestEnvironment(origin string, ref string) (TestEnv, error) {

	if origin == "" {
		return TestEnv{}, errors.New("origin is empty")
	}
	ref = strings.TrimSpace(ref)

	tempDir, err := utils.CreateTempDir()
	if err != nil {
		return TestEnv{}, fmt.Errorf("error creating temporary directory: %w", err)
	}

	err = utils.CloneGitRepository(origin, tempDir, "temp")
	if err != nil {
		return TestEnv{}, fmt.Errorf("unable to clone git repository: %w", err)
	}

	if ref != "" {
		err = utils.ChangeGitHead(tempDir, ref)
		if err != nil {
			return TestEnv{}, fmt.Errorf("error while changing ref: %s", err)
		}
	}

	return TestEnv{
		Dir: tempDir,
	}, nil
}
