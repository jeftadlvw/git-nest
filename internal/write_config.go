package internal

import (
	"bytes"
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

func WriteConfig(c models.NestContext) {

	// write to git_exclude if project is a git repository
	if c.IsGitRepository {
		_ = WriteGitExclude(c.GitRepositoryRoot.Join(gitExcludeFile), c.Config.Submodules)
	}

	// write to git-nest configuration file
	_ = WriteNestConfig(c.ConfigFile, c.Config.Submodules)
}

func WriteGitExclude(p models.Path, modules []models.Submodule) error {
	gitExcludeFileContent := ""
	submoduleGitExcludePart := FmtSubmodulesGitExclude(modules)
	submoduleGitExcludePart = gitExcludeStartString + submoduleGitExcludePart + gitExcludeEndString

	existingContent, err := utils.ReadFileToStr(p)
	if err == nil {
		gitExcludeFileContent, err = utils.StringInsert(existingContent, submoduleGitExcludePart, gitExcludeStartString, gitExcludeEndString)
	} else {
		gitExcludeFileContent = submoduleGitExcludePart
	}

	err = utils.WriteStrToFile(p, gitExcludeFileContent)
	return err
}

func WriteNestConfig(p models.Path, modules []models.Submodule) error {
	existingContent := ""
	existingConfig, err := utils.ReadFileToStr(p)

	if err == nil {
		re := regexp.MustCompile(`\[config](?s:(.*?))(?=\[\[submodule]])`)
		matches := re.FindStringSubmatch(existingConfig)
		if len(matches) > 1 {
			existingContent = matches[1]
		}
	}

	submodulesConfig := SubmodulesToTomlConfig("  ", modules...)

	finalConfigContent := strings.TrimSpace(existingContent) + "\n\n" + strings.TrimSpace(submodulesConfig) + "\n"
	err = utils.WriteStrToFile(p, finalConfigContent)
	return err
}

func FmtSubmodulesGitExclude(submodules []models.Submodule) string {
	buffer := bytes.NewBufferString("")

	t := template.Must(template.New("exclude").Parse(excludeTemplate))
	err := t.Execute(buffer, submodules)
	if err != nil {
		return ""
	}

	return buffer.String()
}
