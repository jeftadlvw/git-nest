package test_env

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	test_env_models "github.com/jeftadlvw/git-nest/test_env/models"
	"github.com/jeftadlvw/git-nest/utils"
	"path/filepath"
	"strings"
)

func CreateTestEnvironment(settings test_env_models.EnvSettings) (test_env_models.TestEnv, error) {

	tempDir, err := utils.CreateTempDir()
	if err != nil {
		return test_env_models.TestEnv{}, fmt.Errorf("error creating temporary directory: %w", err)
	}

	if !settings.NoGit {

		var repositoryDir models.Path

		if !settings.EmptyGit {
			settings.Origin = strings.TrimSpace(settings.Origin)

			if settings.Origin == "" {
				return test_env_models.TestEnv{}, errors.New("test environment origin is empty")
			}

			cloneDirName := strings.TrimSpace(settings.CloneDir)

			err = utils.CloneGitRepository(settings.Origin, tempDir, cloneDirName)
			if err != nil {
				return test_env_models.TestEnv{}, fmt.Errorf("unable to clone git repository: %w", err)
			}
			gitDir := tempDir

			settings.Ref = strings.TrimSpace(settings.Ref)
			if cloneDirName != "" {
				gitDir = gitDir.SJoin(cloneDirName)
			} else {
				gitDir = gitDir.SJoin(strings.TrimSuffix(filepath.Base(settings.Origin), ".git"))
			}

			if settings.Ref != "" {
				err = utils.ChangeGitHead(gitDir, settings.Ref)
				if err != nil {
					return test_env_models.TestEnv{}, fmt.Errorf("error while changing ref: %s", err)
				}
			}

			repositoryDir = gitDir

		} else {
			_, err = utils.RunCommandCombinedOutput(tempDir, "git", "init")
			if err != nil {
				return test_env_models.TestEnv{}, fmt.Errorf("error initializing git repository: %w", err)
			}
			repositoryDir = tempDir
		}

		out, err := utils.RunCommandCombinedOutput(repositoryDir, "git", "config", "user.name", "foo")
		if err != nil {
			return test_env_models.TestEnv{}, fmt.Errorf("error setting local git user at %s: %w; %s", repositoryDir, err, out)
		}

		out, err = utils.RunCommandCombinedOutput(repositoryDir, "git", "config", "user.email", "foo@email.com")
		if err != nil {

			return test_env_models.TestEnv{}, fmt.Errorf("error setting local git email: %w; %s", err, out)
		}
	}

	return test_env_models.TestEnv{
		Dir: tempDir,
	}, nil
}
