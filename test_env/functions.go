package test_env

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/test_env/models"
	"github.com/jeftadlvw/git-nest/utils"
	"path/filepath"
	"strings"
)

func CreateTestEnvironment(settings models.EnvSettings) (models.TestEnv, error) {

	tempDir, err := utils.CreateTempDir()
	if err != nil {
		return models.TestEnv{}, fmt.Errorf("error creating temporary directory: %w", err)
	}

	if !settings.NoGit {

		if !settings.EmptyGit {
			settings.Origin = strings.TrimSpace(settings.Origin)
			if settings.Origin != "" {
				cloneDirName := strings.TrimSpace(settings.CloneDir)

				err = utils.CloneGitRepository(settings.Origin, tempDir, cloneDirName)
				if err != nil {
					return models.TestEnv{}, fmt.Errorf("unable to clone git repository: %w", err)
				}

				settings.Ref = strings.TrimSpace(settings.Ref)
				gitDir := tempDir
				if cloneDirName != "" {
					gitDir = gitDir.SJoin(cloneDirName)
				} else {
					gitDir = gitDir.SJoin(filepath.Base(settings.Origin))
				}

				if settings.Ref != "" {
					err = utils.ChangeGitHead(gitDir, settings.Ref)
					if err != nil {
						return models.TestEnv{}, fmt.Errorf("error while changing ref: %s", err)
					}
				}
			}
		} else {
			_, err = utils.RunCommandCombinedOutput(tempDir, "git", "init")
			if err != nil {
				return models.TestEnv{}, fmt.Errorf("error initializing git repository: %w", err)
			}
		}
	}

	return models.TestEnv{
		Dir: tempDir,
	}, nil
}
