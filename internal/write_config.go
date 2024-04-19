package internal

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"regexp"
	"strings"
)

const gitExcludeFile = ".git/info/exclude"
const gitExcludePrefix string = "# git-nest configuration start"
const gitExcludeSuffix string = "# git-nest configuration end"
const gitExcludeInfo string = `# This part influences how git handles nested modules using git-nest.
# Do not touch except you know what you are doing!`

/*
FmtSubmodulesGitIgnore returns a string that formats a slice of models.Submodule into a string that can be used by git to ignore the submodules' paths.
*/
func FmtSubmodulesGitIgnore(submodules []models.Submodule) string {
	sb := strings.Builder{}
	for _, submodule := range submodules {
		sb.WriteString(submodule.Path.String())
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

		// append '[[submodule]]' to in-memory config for regex to find existing config
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
func WriteProjectConfigFiles(c models.NestContext) (bool, bool, error, error) {

	var gitExcludeWritten, configWritten bool
	var gitExcludeWriteError, configWriteError error

	// write to git_exclude if project is a git repository
	if c.IsGitInstalled && c.IsGitRepository {
		err := WriteSubmoduleIgnoreConfig(c.GitRepositoryRoot.SJoin(gitExcludeFile), c.Config.Submodules)
		if err == nil {
			gitExcludeWritten = true
		}
		gitExcludeWriteError = err
	}

	// write to git-nest configuration file
	err := WriteNestConfig(c.ConfigFile, c.Config.Submodules)
	if err == nil {
		configWritten = true
	}
	configWriteError = err

	return gitExcludeWritten, configWritten, gitExcludeWriteError, configWriteError
}
