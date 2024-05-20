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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestFmtSubmodulesGitExclude-%d", index+1), func(t *testing.T) {
			output := internal.FmtSubmodulesGitIgnore(tc.submodules)
			if output != tc.expected {
				t.Fatalf("Expected:\n>%s<\n\nActual:\n>%s<", tc.expected, output)
			}
		})
	}
}

func TestWriteSubmodulePathIgnoreConfig(t *testing.T) {
	tempFile, err := utils.CreateTempFile("")
	if err != nil {
		t.Fatalf("failed to create temp file")
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestWriteSubmodulePathIgnoreConfig-%d", index+1), func(t *testing.T) {
			expectedContains := gitExcludeInfo
			if tc.expected != "" {
				expectedContains = expectedContains + "\n" + tc.expected
			}
			expectedContains = gitExcludePrefix + "\n" + expectedContains + "\n" + gitExcludeSuffix

			err = internal.WriteSubmoduleIgnoreConfig(tempFile, tc.submodules)
			if tc.error && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.error && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			fileContents, err := utils.ReadFileToStr(tempFile)
			if err != nil {
				t.Fatalf("failed to read temp file")
			}

			if !strings.Contains(fileContents, expectedContains) {
				t.Fatalf("Wexpected string in file\nExpected:\n>%s<\n\nFile:\n>%s<", expectedContains, fileContents)
			}
		})
	}
}

func TestWriteNestConfig(t *testing.T) {
	tempFile, err := utils.CreateTempFile("")
	if err != nil {
		t.Fatalf("failed to create temp file")
		return
	}
	defer os.Remove(tempFile.String())

	tempDir, err := utils.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %s", err)
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
	for index, tc := range cases {

		t.Run(fmt.Sprintf("TestWriteNestConfig-%d", index+1), func(t *testing.T) {
			expectedContains := internal.SubmodulesToTomlConfig("  ", tc.submodules...)

			err = internal.WriteNestConfig(tc.path, tc.submodules)
			if tc.error && err == nil {
				t.Fatalf("WriteNestConfig() for %v returned no error but expected one", tc)
			}
			if !tc.error && err != nil {
				t.Fatalf("WriteNestConfig() for %v returned error, but should've not -> %s", tc, err)
			}

			fileContents, err := utils.ReadFileToStr(tempFile)
			if err != nil {
				t.Fatalf("WriteNestConfig() failed to read temp file")
			}

			if !strings.Contains(fileContents, expectedContains) {
				t.Fatalf("WriteNestConfig() for %v expected string in file\nExpected:\n>%s<\n\nFile:\n>%s<", tc, expectedContains, fileContents)
			}
		})
	}

	// add config header and test for it's existence
	const mockConfig = "[config]\n   foo = \"foo\"\n   bar = \"bar\"\n\t\n\n"
	const mockConfigInFile = "[config]\n   foo = \"foo\"\n   bar = \"bar\"\n"
	err = utils.WriteStrToFile(tempFile, mockConfig)
	if err != nil {
		t.Fatalf("failed to write test config section: %s", err)
		return
	}
	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestWriteNestConfig-%d", index+1), func(t *testing.T) {
			expectedContains := internal.SubmodulesToTomlConfig("  ", tc.submodules...)

			err = internal.WriteNestConfig(tc.path, tc.submodules)
			if tc.error && err == nil {
				t.Fatalf("WriteNestConfig() for %v returned no error but expected one", tc)
			}
			if !tc.error && err != nil {
				t.Fatalf("WriteNestConfig() for %v returned error, but should've not -> %s", tc, err)
			}

			fileContents, err := utils.ReadFileToStr(tempFile)
			if err != nil {
				t.Fatalf("WriteNestConfig() failed to read temp file")
			}

			if !strings.Contains(fileContents, expectedContains) {
				t.Fatalf("WriteNestConfig() for %v expected string in file\nExpected:\n>%s<\n\nFile:\n>%s<", tc, expectedContains, fileContents)
			}

			if !strings.Contains(fileContents, mockConfigInFile) {
				t.Fatalf("WriteNestConfig() for %v removed config section from config file: %s", tc, fileContents)
			}
		})
	}
}

func TestWriteProjectConfigFiles(t *testing.T) {

	prepareGitRepository := func() (models.Path, error) {
		tempDir, err := utils.CreateTempDir()
		if err != nil {
			return "", fmt.Errorf("failed to create temp dir: %w", err)
		}

		err = utils.CloneGitRepository(testRepoUrl, tempDir, ".", nil)
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
			t.Fatalf("failed to prepare git repository: %s", err)
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

		r, err := internal.WriteProjectConfigFiles(context)

		if err != nil {
			t.Fatalf("failed to write project config files: %s", err)
		}

		if r.GitExcludeWriteError != nil {
			t.Fatalf("writing to git exclude failed: %s", r.GitExcludeWriteError)
		}

		if r.ConfigWriteError != nil {
			t.Fatalf("writing to config file failed: %s", r.ConfigWriteError)
		}

		if !r.GitExcludeWritten {
			t.Fatalf("should've written to git exclude")
		}

		if !r.ConfigWritten {
			t.Fatalf("should've written to configuration file")
		}

		configContents, err := utils.ReadFileToStr(configFile)
		if err != nil {
			t.Fatalf("error reading configuration file: %s", err)
		}

		gitExcludeContents, err := utils.ReadFileToStr(absGitExcludeFile)
		if err != nil {
			t.Fatalf("error reading git exclude file: %s", err)
		}

		if modulesExist {
			if !strings.Contains(configContents, internal.SubmodulesToTomlConfig("  ", context.Config.Submodules...)) {
				t.Fatalf("configuration file does not contain submodules")
			}

			if !strings.Contains(gitExcludeContents, internal.FmtSubmodulesGitIgnore(context.Config.Submodules)) {
				t.Fatalf("git exclude file does not contain submodules")
			}
		}
	})

	// git does not exist
	t.Run("TestWriteProjectConfigFiles-2", func(t *testing.T) {
		t.Parallel()

		tempDir, err := prepareGitRepository()
		if err != nil {
			t.Fatalf("failed to prepare git repository: %s", err)
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

		r, err := internal.WriteProjectConfigFiles(context)

		if err != nil {
			t.Fatalf("failed to write project config files: %s", err)
		}

		if r.GitExcludeWriteError != nil {
			t.Fatalf("writing to git exclude failed, and should've never happened: %s", r.GitExcludeWriteError)
		}

		if r.ConfigWriteError != nil {
			t.Fatalf("writing to config file failed: %s", r.ConfigWriteError)
		}

		if r.GitExcludeWritten {
			t.Fatalf("should've written to git exclude")
		}

		if !r.ConfigWritten {
			t.Fatalf("should've written to configuration file")
		}

		configContents, err := utils.ReadFileToStr(configFile)
		if err != nil {
			t.Fatalf("error reading configuration file: %s", err)
		}

		gitExcludeContents, err := utils.ReadFileToStr(absGitExcludeFile)
		if err != nil {
			t.Fatalf("error reading git exclude file: %s", err)
		}

		if modulesExist {
			if !strings.Contains(configContents, internal.SubmodulesToTomlConfig("  ", context.Config.Submodules...)) {
				t.Fatalf("configuration file does not contain submodules")
			}

			if strings.Contains(gitExcludeContents, internal.FmtSubmodulesGitIgnore(context.Config.Submodules)) {
				t.Fatalf("git exclude file does contains submodules")
			}
		}
	})
}
