package io

import (
	"bytes"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"strings"
	"text/template"
)

const startString string = "# git-nest configuration start"
const endString string = "# git-nest configuration end"

const excludeTemplate string = `# This part influences how git handles nested modules using git-nest.
# Do not touch except you know what you are doing!
{{- range . }}
	{{ .Path }}
{{- end }}`

// TODO extract git-nest configuration part location from .git/info/exclude and override it with new configuration
func FmtSubmodulesGitExclude(submodules []models.Submodule) string {
	buffer := bytes.NewBufferString("")

	t := template.Must(template.New("exclude").Parse(excludeTemplate))
	err := t.Execute(buffer, submodules)
	if err != nil {
		return ""
	}

	return buffer.String()
}
