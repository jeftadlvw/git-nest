package internal

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/conversions"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/jeftadlvw/git-nest/io"
	"github.com/jeftadlvw/git-nest/models"
	"os"
)

/*
EvaluateContext evaluates the runtime context and writes it into constants.Context.
*/
func EvaluateContext() error {

	// get cwd
	cwdStr, err := os.Getwd()
	if err != nil {
		return err
	}
	cwd := models.Path(cwdStr)
	constants.Context.WorkingDirectory = cwd

	// evaluate project root
	var projectRoot models.Path
	projectRoot, err = FindProjectRoot(cwd)
	if err != nil {
		projectRoot = cwd
	}
	constants.Context.ProjectRoot = projectRoot

	// evaluate configuration file path
	configFilePath := evaluateConfigFileFromProjectRoot(projectRoot)
	if configFilePath.IsDir() {
		return fmt.Errorf("configuration path to %s is a directory", configFilePath.String())
	}
	constants.Context.ConfigFile = configFilePath

	// read configuration file
	configStr, err := io.ReadFileToStr(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading configuration file %s: %w", configFilePath.String(), err)
	}

	// populate configuration struct
	nestConfig := models.NestConfig{}
	err = conversions.PopulateNestConfigFromToml(&nestConfig, configStr)
	if err != nil {
		return err
	}
	constants.Context.Config = nestConfig

	return nil
}

func evaluateConfigFileFromProjectRoot(root models.Path) models.Path {
	if root.BContains(constants.ConfigSubDirFileName) {
		return root.Join(constants.ConfigSubDirFileName)
	}

	return root.Join(constants.ConfigFileName)
}
