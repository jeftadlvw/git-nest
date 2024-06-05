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

func CreateTestEnvironment(root models.Path, settings test_env_models.EnvSettings) error {

	var err error

	tempDir := root.Clean()

	if !settings.NoGit {

		var repositoryDir models.Path

		if !settings.EmptyGit {
			settings.Origin = strings.TrimSpace(settings.Origin)

			if settings.Origin == "" {
				return errors.New("test environment origin is empty")
			}

			cloneDirName := strings.TrimSpace(settings.CloneDir)

			err = utils.CloneGitRepository(settings.Origin, tempDir, cloneDirName, nil)
			if err != nil {
				return fmt.Errorf("unable to clone git repository: %w", err)
			}
			gitDir := tempDir

			settings.Ref = strings.TrimSpace(settings.Ref)
			if cloneDirName != "" {
				gitDir = gitDir.SJoin(cloneDirName)
			} else {
				gitDir = gitDir.SJoin(strings.TrimSuffix(filepath.Base(settings.Origin), ".git"))
			}

			if settings.Ref != "" {
				err = utils.GitCheckout(gitDir, settings.Ref)
				if err != nil {
					return fmt.Errorf("error while changing ref: %s", err)
				}
			}

			repositoryDir = gitDir

		} else {
			_, err = utils.RunCommandCombinedOutput(tempDir, "git", "init")
			if err != nil {
				return fmt.Errorf("error initializing git repository: %w", err)
			}
			repositoryDir = tempDir
		}

		out, err := utils.RunCommandCombinedOutput(repositoryDir, "git", "config", "user.name", "foo")
		if err != nil {
			return fmt.Errorf("error setting local git user at %s: %w; %s", repositoryDir, err, out)
		}

		out, err = utils.RunCommandCombinedOutput(repositoryDir, "git", "config", "user.email", "foo@email.com")
		if err != nil {

			return fmt.Errorf("error setting local git email: %w; %s", err, out)
		}
	}

	return nil
}
