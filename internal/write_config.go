package internal

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"os"
	"regexp"
	"strings"
)

const gitExcludeDirectory = ".git/info"
const gitExcludeFile = ".git/info/exclude"
const gitExcludePrefix string = "# git-nest configuration start"
const gitExcludeSuffix string = "# git-nest configuration end"
const gitExcludeInfo string = `# This part influences how git handles nested modules using git-nest.
# Do not touch except you know what you are doing!`

type WriteProjectConfigFilesReturn struct {
	ConfigWritten        bool
	GitExcludeWritten    bool
	ConfigWriteError     error
	GitExcludeWriteError error
}

/*
FmtSubmodulesGitIgnore returns a string that formats a slice of models.Submodule into a string that can be used by git to ignore the submodules' paths.
*/
func FmtSubmodulesGitIgnore(submodules []models.Submodule) string {
	sb := strings.Builder{}
	for _, submodule := range submodules {
		sb.WriteString(submodule.Path.UnixString())
		sb.WriteString("\n")
	}

	return strings.TrimSpace(sb.String())
}

/*
WriteSubmoduleIgnoreConfig uses internal.FmtSubmodulesGitIgnore, wraps it with some user information and writes that
into the passed file. Pre-existing configuration is replaced using utils.StringInsertAtFirst.
*/
func WriteSubmoduleIgnoreConfig(p models.Path, modules []models.Submodule) error {
	submoduleGitExcludeFmt := FmtSubmodulesGitIgnore(modules)
	submoduleGitExcludePart := gitExcludePrefix + "\n" + gitExcludeInfo + "\n"

	if submoduleGitExcludeFmt != "" {
		submoduleGitExcludePart = submoduleGitExcludePart + submoduleGitExcludeFmt + "\n"
	}
	submoduleGitExcludePart = submoduleGitExcludePart + gitExcludeSuffix

	fileContent := submoduleGitExcludePart

	existingContent, err := utils.ReadFileToStr(p)
	if err == nil {
		localFileContent, localErr := utils.StringInsertAtFirst(existingContent, submoduleGitExcludePart, gitExcludePrefix, gitExcludeSuffix)
		if localErr == nil {
			fileContent = localFileContent
		}
	}

	err = utils.WriteStrToFile(p, fileContent)
	return err
}

/*
WriteNestConfig writes models.Submodule configuration into the git-nest configuration file, preserving the first [config] section.
*/
func WriteNestConfig(p models.Path, modules []models.Submodule) error {
	existingContent := ""
	existingConfig := ""

	if p.Empty() {
		return fmt.Errorf("cannot write to empty path")
	}

	if p.IsDir() {
		return fmt.Errorf("passed path is a directory: %s", p)
	}

	if p.IsFile() {
		localExistingConfig, err := utils.ReadFileToStr(p)
		if err != nil {
			return fmt.Errorf("cannot read existing config: %w", err)
		}
		existingConfig = localExistingConfig

		// append '[[submodule]]' to in-memory config for regex to find existing config;
		// even if that header already exists, it'll be removed anyway by later processing
		existingConfig = existingConfig + "[[submodule]]"
	}

	// find any existing config section
	re := regexp.MustCompile(`\[config](?s:(.*?))\[\[submodule]]`)
	matches := re.FindStringSubmatch(existingConfig)

	if len(matches) > 1 {
		existingContent = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(matches[0]), "[[submodule]]"))
	}

	submodulesConfig := SubmodulesToTomlConfig("  ", modules...)

	if existingContent != "" {
		existingContent = existingContent + "\n\n"
	}

	finalConfigContent := existingContent + submodulesConfig + "\n"
	err := utils.WriteStrToFile(p, finalConfigContent)
	if err != nil {
		return fmt.Errorf("cannot write 'nestmodules.toml': %w", err)
	}

	return nil
}

/*
WriteProjectConfigFiles is a total wrapper function for internal.WriteSubmoduleIgnoreConfig and internal.WriteNestConfig,
calling both functions based on the passed models.NestContext.
*/
func WriteProjectConfigFiles(c models.NestContext) (WriteProjectConfigFilesReturn, error) {

	r := WriteProjectConfigFilesReturn{}

	// check if configuration file has been updated since initial
	// context evaluation
	if c.ConfigFile.IsFile() {
		localChecksum, err := utils.CalculateChecksumF(c.ConfigFile)
		if err != nil {
			return r, fmt.Errorf("internal error: could not calculate checksum: %w", err)
		}

		if c.Checksums.ConfigurationFile != localChecksum {
			return r, fmt.Errorf("configuration file checksum mismatch:\nThe configuration file has been changed since the initial start of the program.")
		}
	}

	// write nest config first, as
	// write to git-nest configuration file
	err := WriteNestConfig(c.ConfigFile, c.Config.Submodules)
	if err == nil {
		r.ConfigWritten = true
	}
	r.ConfigWriteError = err

	// write to git_exclude if project is a git repository
	if c.IsGitInstalled && c.IsGitRepository {

		// create gitExcludeDirectory if it does not exist
		gitExcludeDirectoryPath := c.GitRepositoryRoot.SJoin(gitExcludeDirectory)
		if gitExcludeDirectoryPath.IsFile() {
			r.GitExcludeWriteError = fmt.Errorf("%s is a file", gitExcludeDirectoryPath)
		}

		if !gitExcludeDirectoryPath.IsFile() && !gitExcludeDirectoryPath.IsDir() {
			err = os.MkdirAll(gitExcludeDirectoryPath.String(), os.ModePerm)
			if err != nil {
				r.GitExcludeWriteError = fmt.Errorf("cannot create directory %s: %w", gitExcludeDirectoryPath.String(), err)
			}
		}

		if gitExcludeDirectoryPath.IsDir() {
			err = WriteSubmoduleIgnoreConfig(c.GitRepositoryRoot.SJoin(gitExcludeFile), c.Config.Submodules)
			if err == nil {
				r.GitExcludeWritten = true
			} else {
				r.GitExcludeWriteError = err
			}
		}
	}

	return r, nil
}
