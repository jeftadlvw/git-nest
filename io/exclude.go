package io

import (
	"bytes"
	"github.com/jeftadlvw/git-nest/models"
	template2 "html/template"
)

const startString string = "# git-nest configuration start"
const endString string = "# git-nest configuration end"

const excludeTemplate string = `# This part influences how git handles nested modules using git-nest.
# Do not touch except you know what you are doing!
{{- range . }}
	{{ .Path }}
{{- end }}`

// TODO extract git-nest configuration part location from .git/info/exclude and override it with new configuration
func FmtSubmodulesGitIgnore(submodules []models.Submodule) string {
	buffer := bytes.NewBufferString("")

	template := template2.Must(template2.New("exclude").Parse(excludeTemplate))
	err := template.Execute(buffer, submodules)
	if err != nil {
		return ""
	}

	return buffer.String()
}
