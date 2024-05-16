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
CreateContext returns a fresh evaluated models.NestContext for the passed path.
*/
func CreateContext(p models.Path) (models.NestContext, error) {

	var (
		projectRoot      models.Path
		configFileExists bool
		configFilePath   models.Path
		nestConfig       models.NestConfig
		gitRoot          models.Path
		IsGitInstalled   bool
		isGitProject     bool
		err              error
	)

	nestContext := models.NestContext{}

	if !p.IsDir() {
		return nestContext, fmt.Errorf("%s is not a directory", p.String())
	}

	// evaluate project root
	projectRoot, err = FindProjectRoot(p)
	if err != nil {
		projectRoot = p
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
	} else {
		configStr = ""
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
			_, _ = fmt.Fprintf(os.Stderr, "git-nest root and git repository root directories do not match: %s != %s\n", gitRoot.String(), projectRoot.String())
		}
	}

	// calculate checksum of configuration file content
	configFileChecksum := utils.CalculateChecksumS(configStr)

	nestContext.WorkingDirectory = p
	nestContext.ProjectRoot = projectRoot
	nestContext.ConfigFileExists = configFileExists
	nestContext.ConfigFile = configFilePath
	nestContext.Config = nestConfig
	nestContext.IsGitInstalled = IsGitInstalled
	nestContext.IsGitRepository = isGitProject
	nestContext.GitRepositoryRoot = gitRoot
	nestContext.Checksums.ConfigurationFile = configFileChecksum

	return nestContext, nil
}

/*
CreateContextFromCurrentWorkingDir returns a fresh evaluated models.NestContext from the current working directory.
*/
func CreateContextFromCurrentWorkingDir() (models.NestContext, error) {
	cwdStr, err := os.Getwd()
	if err != nil {
		return models.NestContext{}, err
	}
	cwd := models.Path(cwdStr)

	return CreateContext(cwd)
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
