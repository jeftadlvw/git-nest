package internal

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/conversions"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/jeftadlvw/git-nest/io"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"os/exec"
)

/*
EvaluateContext evaluates the runtime context and writes it into constants.Context.
*/
func EvaluateContext() error {

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

	// get cwd
	cwdStr, err := os.Getwd()
	if err != nil {
		return err
	}
	cwd = models.Path(cwdStr)

	// evaluate project root
	projectRoot, err = FindProjectRoot(cwd)
	if err != nil {
		projectRoot = cwd
	}

	// evaluate configuration file path
	configFilePath = evaluateConfigFileFromProjectRoot(projectRoot)
	if configFilePath.IsDir() {
		return fmt.Errorf("configuration path to %s is a directory", configFilePath.String())
	}

	// read configuration file
	configFileExists = false
	configStr, err := io.ReadFileToStr(configFilePath)
	if err == nil {
		configFileExists = true
	}

	// populate configuration struct if a configuration file exists,
	// else nestConfig is an empty configuration
	nestConfig = models.NestConfig{}
	if configFileExists {
		err = conversions.PopulateNestConfigFromToml(&nestConfig, configStr)
		if err != nil {
			return err
		}
	}

	// check if project root is also a git repository
	gitRootStr, err := utils.FindGitRoot()
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

	constants.Context.WorkingDirectory = cwd
	constants.Context.ProjectRoot = projectRoot
	constants.Context.ConfigFileExists = configFileExists
	constants.Context.ConfigFile = configFilePath
	constants.Context.Config = nestConfig
	constants.Context.IsGitInstalled = IsGitInstalled
	constants.Context.IsGitProject = isGitProject

	return nil
}

func evaluateConfigFileFromProjectRoot(root models.Path) models.Path {
	if root.BContains(constants.ConfigSubDirFileName) {
		return root.Join(constants.ConfigSubDirFileName)
	}

	return root.Join(constants.ConfigFileName)
}
