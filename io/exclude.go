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

func InsertInto(original string, insert string, startDetermination string, endDetermination string) (string, error) {

	startDeterminationCount := strings.Count(original, startDetermination)
	endDeterminationCount := strings.Count(original, endDetermination)

	if startDeterminationCount == 0 {
		return "", fmt.Errorf("cannot to find starting determinator")
	}

	if endDeterminationCount == 0 {
		return "", fmt.Errorf("cannot to find ending determinator")
	}

	if startDeterminationCount > 1 {
		return "", fmt.Errorf("start determinator found multiple times")
	}

	if endDeterminationCount > 1 {
		return "", fmt.Errorf("end determinator found multiple times")
	}

	before := strings.Split(original, startDetermination)
	after := strings.Split(original, endDetermination)

	/*
		TODO: replace with regex-version
		// the regex here is untested
		var re = regexp.MustCompile(`(startDetermination).*(endDetermination)`)
		s := re.ReplaceAllString(sample, `${1}foo`)
	*/

	return before[0] + insert + after[len(after)-1], nil
}
