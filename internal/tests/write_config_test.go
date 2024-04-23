package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/internal/constants"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"strings"
	"testing"
)

const testRepoUrl = "https://github.com/jeftadlvw/example-repository.git"
const gitExcludeFile = ".git/info/exclude"
const gitExcludePrefix string = "# git-nest configuration start"
const gitExcludeSuffix string = "# git-nest configuration end"
const gitExcludeInfo string = `# This part influences how git handles nested modules using git-nest.
# Do not touch except you know what you are doing!`

func TestFmtSubmodulesGitExclude(t *testing.T) {
	cases := []struct {
		submodules []models.Submodule
		expected   string
	}{
		{submodules: []models.Submodule{}, expected: ""},
		{submodules: []models.Submodule{{Path: "submodule1"}}, expected: "submodule1"},
		{submodules: []models.Submodule{{Path: "submodule1"}, {Path: "submodule2"}, {Path: "submodule3"}}, expected: "submodule1\nsubmodule2\nsubmodule3"},
	}

	for _, tc := range cases {
		output := internal.FmtSubmodulesGitIgnore(tc.submodules)
		if output != tc.expected {
			t.Errorf("FmtSubmodulesGitIgnore() for %v\nExpected:\n>%s<\n\nActual:\n>%s<", tc, tc.expected, output)
		}
	}
}

func TestWriteSubmodulePathIgnoreConfig(t *testing.T) {
	tempFile, err := utils.CreateTempFile("")
	if err != nil {
		t.Errorf("failed to create temp file")
		return
	}
	defer os.Remove(tempFile.String())

	utils.WriteStrToFile(tempFile, "Foocontent\n\n\t\n")

	// cases are run sequentially on same file
	cases := []struct {
		submodules []models.Submodule
		expected   string
		error      bool
	}{
		{[]models.Submodule{}, "", false},
		{[]models.Submodule{{Path: "submodule1"}}, "submodule1", false},
		{[]models.Submodule{{Path: "submodule1"}, {Path: "submodule2"}, {Path: "submodule3"}}, "submodule1\nsubmodule2\nsubmodule3", false},
		{[]models.Submodule{}, "", false},
	}

	for _, tc := range cases {
		expectedContains := gitExcludeInfo
		if tc.expected != "" {
			expectedContains = expectedContains + "\n" + tc.expected
		}
		expectedContains = gitExcludePrefix + "\n" + expectedContains + "\n" + gitExcludeSuffix

		err = internal.WriteSubmoduleIgnoreConfig(tempFile, tc.submodules)
		if tc.error && err == nil {
			t.Errorf("WriteSubmoduleIgnoreConfig() for %v returned no error but expected one", tc)
			continue
		}
		if !tc.error && err != nil {
			t.Errorf("WriteSubmoduleIgnoreConfig() for %v returned error, but should've not -> %s", tc, err)
			continue
		}

		fileContents, err := utils.ReadFileToStr(tempFile)
		if err != nil {
			t.Errorf("WriteSubmoduleIgnoreConfig() failed to read temp file")
			continue
		}

		if !strings.Contains(fileContents, expectedContains) {
			t.Errorf("WriteSubmoduleIgnoreConfig() for %v expected string in file\nExpected:\n>%s<\n\nFile:\n>%s<", tc, expectedContains, fileContents)
			continue
		}
	}
}

func TestWriteNestConfig(t *testing.T) {
	tempFile, err := utils.CreateTempFile("")
	if err != nil {
		t.Errorf("failed to create temp file")
		return
	}
	defer os.Remove(tempFile.String())

	tempDir, err := utils.CreateTempDir()
	if err != nil {
		t.Errorf("failed to create temp dir: %s", err)
		return
	}
	defer os.RemoveAll(tempDir.String())

	// cases are run sequentially on same file
	cases := []struct {
		path       models.Path
		submodules []models.Submodule
		error      bool
	}{
		{"", []models.Submodule{}, true},
		{tempDir, []models.Submodule{}, true},
		{tempFile, []models.Submodule{}, false},
		{tempFile, []models.Submodule{{Path: "submodule1"}}, false},
		{tempFile, []models.Submodule{{Path: "submodule1"}, {Path: "submodule2"}, {Path: "submodule3"}}, false},
		{tempFile, []models.Submodule{}, false},
	}

	// test without config header in file
	for _, tc := range cases {
		expectedContains := internal.SubmodulesToTomlConfig("  ", tc.submodules...)

		err = internal.WriteNestConfig(tc.path, tc.submodules)
		if tc.error && err == nil {
			t.Errorf("WriteNestConfig() for %v returned no error but expected one", tc)
			continue
		}
		if !tc.error && err != nil {
			t.Errorf("WriteNestConfig() for %v returned error, but should've not -> %s", tc, err)
			continue
		}

		fileContents, err := utils.ReadFileToStr(tempFile)
		if err != nil {
			t.Errorf("WriteNestConfig() failed to read temp file")
			continue
		}

		if !strings.Contains(fileContents, expectedContains) {
			t.Errorf("WriteNestConfig() for %v expected string in file\nExpected:\n>%s<\n\nFile:\n>%s<", tc, expectedContains, fileContents)
			continue
		}
	}

	// add config header and test for it's existence
	const mockConfig = "[config]\n   foo = \"foo\"\n   bar = \"bar\"\n\t\n\n"
	const mockConfigInFile = "[config]\n   foo = \"foo\"\n   bar = \"bar\"\n"
	err = utils.WriteStrToFile(tempFile, mockConfig)
	if err != nil {
		t.Errorf("failed to write test config section: %s", err)
		return
	}
	for _, tc := range cases {
		expectedContains := internal.SubmodulesToTomlConfig("  ", tc.submodules...)

		err = internal.WriteNestConfig(tc.path, tc.submodules)
		if tc.error && err == nil {
			t.Errorf("WriteNestConfig() for %v returned no error but expected one", tc)
			continue
		}
		if !tc.error && err != nil {
			t.Errorf("WriteNestConfig() for %v returned error, but should've not -> %s", tc, err)
			continue
		}

		fileContents, err := utils.ReadFileToStr(tempFile)
		if err != nil {
			t.Errorf("WriteNestConfig() failed to read temp file")
			continue
		}

		if !strings.Contains(fileContents, expectedContains) {
			t.Errorf("WriteNestConfig() for %v expected string in file\nExpected:\n>%s<\n\nFile:\n>%s<", tc, expectedContains, fileContents)
			continue
		}

		if !strings.Contains(fileContents, mockConfigInFile) {
			t.Errorf("WriteNestConfig() for %v removed config section from config file: %s", tc, fileContents)
			continue
		}
	}
}

func TestWriteProjectConfigFiles(t *testing.T) {

	prepareGitRepository := func() (models.Path, error) {
		tempDir, err := utils.CreateTempDir()
		if err != nil {
			return "", fmt.Errorf("failed to create temp dir: %w", err)
		}

		err = utils.CloneGitRepository(testRepoUrl, tempDir, ".")
		if err != nil {
			return "", fmt.Errorf("failed to clone git repository: %w", err)
		}

		return tempDir, nil
	}

	// git exists
	t.Run("TestWriteProjectConfigFiles-1", func(t *testing.T) {
		t.Parallel()

		tempDir, err := prepareGitRepository()
		if err != nil {
			t.Errorf("failed to prepare git repository: %s", err)
			return
		}
		defer os.RemoveAll(tempDir.String())

		var configFile = tempDir.SJoin(constants.ConfigFileName)
		var absGitExcludeFile = tempDir.SJoin(gitExcludeFile)

		context := models.NestContext{
			ProjectRoot:       tempDir,
			GitRepositoryRoot: tempDir,
			ConfigFileExists:  false,
			ConfigFile:        configFile,
			Config:            models.NestConfig{},
			IsGitInstalled:    true,
			IsGitRepository:   true,
		}

		modulesExist := len(context.Config.Submodules) != 0

		gitExcludeWritten, configWritten, gitExcludeWriteErr, configWriteErr := internal.WriteProjectConfigFiles(context)

		if gitExcludeWriteErr != nil {
			t.Errorf("Writing to git exclude failed: %s", gitExcludeWriteErr)
		}

		if configWriteErr != nil {
			t.Errorf("Writing to config file failed: %s", configWriteErr)
		}

		if !gitExcludeWritten {
			t.Errorf("WriteProjectConfigFiles() should've written to git exclude")
		}

		if !configWritten {
			t.Errorf("WriteProjectConfigFiles() should've written to configuration file")
		}

		configContents, err := utils.ReadFileToStr(configFile)
		if err != nil {
			t.Errorf("error reading configuration file: %s", err)
		}

		gitExcludeContents, err := utils.ReadFileToStr(absGitExcludeFile)
		if err != nil {
			t.Errorf("error reading git exclude file: %s", err)
		}

		if modulesExist {
			if !strings.Contains(configContents, internal.SubmodulesToTomlConfig("  ", context.Config.Submodules...)) {
				t.Errorf("configuration file does not contain submodules")
			}

			if !strings.Contains(gitExcludeContents, internal.FmtSubmodulesGitIgnore(context.Config.Submodules)) {
				t.Errorf("git exclude file does not contain submodules")
			}
		}
	})

	// git does not exist
	t.Run("TestWriteProjectConfigFiles-2", func(t *testing.T) {
		t.Parallel()

		tempDir, err := prepareGitRepository()
		if err != nil {
			t.Errorf("failed to prepare git repository: %s", err)
			return
		}
		defer os.RemoveAll(tempDir.String())

		var configFile = tempDir.SJoin(constants.ConfigFileName)
		var absGitExcludeFile = tempDir.SJoin(gitExcludeFile)

		context := models.NestContext{
			ProjectRoot:       tempDir,
			GitRepositoryRoot: tempDir,
			ConfigFileExists:  false,
			ConfigFile:        configFile,
			Config:            models.NestConfig{},
			IsGitInstalled:    false, // setting this should be enough!
			IsGitRepository:   true,
		}

		modulesExist := len(context.Config.Submodules) != 0

		gitExcludeWritten, configWritten, gitExcludeWriteErr, configWriteErr := internal.WriteProjectConfigFiles(context)

		if gitExcludeWriteErr != nil {
			t.Errorf("Writing to git exclude failed, and should've never happened: %s", gitExcludeWriteErr)
		}

		if configWriteErr != nil {
			t.Errorf("Writing to config file failed: %s", configWriteErr)
		}

		if gitExcludeWritten {
			t.Errorf("WriteProjectConfigFiles() should not have written to git exclude")
		}

		if !configWritten {
			t.Errorf("WriteProjectConfigFiles() should've written to configuration file")
		}

		configContents, err := utils.ReadFileToStr(configFile)
		if err != nil {
			t.Errorf("error reading configuration file: %s", err)
		}

		gitExcludeContents, err := utils.ReadFileToStr(absGitExcludeFile)
		if err != nil {
			t.Errorf("error reading git exclude file: %s", err)
		}

		if modulesExist {
			if !strings.Contains(configContents, internal.SubmodulesToTomlConfig("  ", context.Config.Submodules...)) {
				t.Errorf("configuration file does not contain submodules")
			}

			if strings.Contains(gitExcludeContents, internal.FmtSubmodulesGitIgnore(context.Config.Submodules)) {
				t.Errorf("git exclude file does contains submodules")
			}
		}
	})

}