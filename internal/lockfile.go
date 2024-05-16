package internal

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
)

/*
AcquireLockFile tries to acquire a context-specific lockfile.
It returns if the lockfile was acquired successfully.
*/
func AcquireLockFile(c models.NestContext) (bool, error) {
	lockFilePath := getLockFilePath(c.IsGitRepository, c.ProjectRoot)

	if lockFilePath.IsDir() {
		return false, errors.New("lock file is directory")
	}

	if lockFilePath.IsFile() {
		return false, nil
	}

	err := utils.WriteStrToFile(lockFilePath, constants.LockFileContents)
	if err != nil {
		return false, fmt.Errorf("could not write lockfile: %w", err)
	}

	return true, nil
}

/*
ReleaseLockFile releases the context-specific lockfile.
*/
func ReleaseLockFile(c models.NestContext) error {
	lockFilePath := getLockFilePath(c.IsGitRepository, c.ProjectRoot)

	if lockFilePath.IsDir() {
		return errors.New("lock file is directory")
	}

	if lockFilePath.IsFile() {
		err := os.Remove(lockFilePath.String())
		if err != nil {
			return fmt.Errorf("could not remove lock file: %w", err)
		}
	}

	return nil

}

func getLockFilePath(isGitRepository bool, projectRoot models.Path) models.Path {
	if isGitRepository {
		return projectRoot.SJoin(constants.LockFileNameGitRepo)
	}

	return projectRoot.SJoin(constants.LockFileName)
}
