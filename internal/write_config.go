package internal

import (
	"bytes"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"regexp"
	"strings"
	"text/template"
)

const gitExcludeFile = ".git/info/exclude"
const gitExcludeStartString string = "# git-nest configuration start"
const gitExcludeEndString string = "# git-nest configuration end"
const excludeTemplate string = `# This part influences how git handles nested modules using git-nest.
# Do not touch except you know what you are doing!
{{- range . }}
{{ .Path }}
{{- end }}`

func FmtSubmodulesGitExclude(submodules []models.Submodule) string {
	buffer := bytes.NewBufferString("")

	t := template.Must(template.New("exclude").Parse(excludeTemplate))
	err := t.Execute(buffer, submodules)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(buffer.String())
}

func WriteSubmodulePathIgnoreConfig(p models.Path, modules []models.Submodule) error {
	submoduleGitExcludePart := FmtSubmodulesGitExclude(modules)
	submoduleGitExcludePart = gitExcludeStartString + "\n" + submoduleGitExcludePart + "\n" + gitExcludeEndString

	fileContent := submoduleGitExcludePart

	existingContent, err := utils.ReadFileToStr(p)
	if err == nil {
		localFileContent, localErr := utils.StringInsertAtFirst(existingContent, submoduleGitExcludePart, gitExcludeStartString, gitExcludeEndString)
		if localErr == nil {
			fileContent = localFileContent
		}
	}

	err = utils.WriteStrToFile(p, fileContent)
	return err
}

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

func WriteProjectConfigFiles(c models.NestContext) (bool, bool, error, error) {

	var gitExcludeWritten, configWritten bool
	var gitExcludeWriteError, configWriteError error

	// write to git_exclude if project is a git repository
	if c.IsGitInstalled && c.IsGitRepository {
		err := WriteSubmodulePathIgnoreConfig(c.GitRepositoryRoot.SJoin(gitExcludeFile), c.Config.Submodules)
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
