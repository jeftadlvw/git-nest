package internal

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"os/exec"
)

/*
CreateContext returns a fresh evaluated models.NestContext.
*/
func CreateContext() (models.NestContext, error) {

	var (
		cwd              models.Path
		projectRoot      models.Path
		configFileExists bool
		configFilePath   models.Path
		nestConfig       models.NestConfig
		gitRoot          models.Path
		IsGitInstalled   bool
		isGitProject     bool
	)

	nestContext := models.NestContext{}

	// get cwd
	cwdStr, err := os.Getwd()
	if err != nil {
		return nestContext, err
	}
	cwd = models.Path(cwdStr)

	// evaluate project root
	projectRoot, err = FindProjectRoot(cwd)
	if err != nil {
		projectRoot = cwd
	}

	// evaluate configuration file path
	configFilePath = evaluateConfigFileFromDir(projectRoot)
	if configFilePath.IsDir() {
		return nestContext, fmt.Errorf("configuration path to %s is a directory", configFilePath.String())
	}

	// read configuration file
	configFileExists = false
	configStr, err := utils.ReadFileToStr(configFilePath)
	if err == nil {
		configFileExists = true
	}

	// populate configuration struct if a configuration file exists,
	// else nestConfig is an empty configuration
	nestConfig = models.NestConfig{}
	if configFileExists {
		err = PopulateNestConfigFromToml(&nestConfig, configStr, false)
		if err != nil {
			return nestContext, err
		}
	}

	// check if project root is also a git repository
	gitRootStr, err := utils.GetGitRootDirectory(projectRoot)
	IsGitInstalled = false
	isGitProject = false
	if err != nil {
		IsGitInstalled = !errors.Is(err, exec.ErrNotFound)
	} else {
		IsGitInstalled = true

		gitRoot = models.Path(gitRootStr)
		if gitRoot == projectRoot {
			isGitProject = true
		} else if !nestConfig.Config.AllowUnequalRoots {
			_, _ = fmt.Fprintf(os.Stderr, "git-nest root and git repository root directories do not match:%s != %s", gitRoot.String(), projectRoot.String())
		}
	}

	nestContext.WorkingDirectory = cwd
	nestContext.ProjectRoot = projectRoot
	nestContext.ConfigFileExists = configFileExists
	nestContext.ConfigFile = configFilePath
	nestContext.Config = nestConfig
	nestContext.IsGitInstalled = IsGitInstalled
	nestContext.IsGitRepository = isGitProject
	nestContext.GitRepositoryRoot = gitRoot

	return nestContext, nil
}

/*
evaluateConfigFileFromDir returns an absolute path to the git-nest configuration file for a given directory.
Defaults to constants.ConfigFileName if the file at constants.ConfigSubDirFileName does not exist.
*/
func evaluateConfigFileFromDir(d models.Path) models.Path {
	if d.BContains(constants.ConfigSubDirFileName) {
		return d.SJoin(constants.ConfigSubDirFileName)
	}

	return d.SJoin(constants.ConfigFileName)
}
