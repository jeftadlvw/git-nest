package utils

import (
	"github.com/jeftadlvw/git-nest/models"
	"os/exec"
	"strings"
)

/*
RunCommand is a subset-wrapper for exec.Command, providing seperate return values for stdout and stderr.
*/
func RunCommand(d models.Path, command string, args ...string) (string, string, error) {
	cmd := exec.Command(command, args...)
	if !d.Empty() {
		cmd.Dir = d.String()
	}

	var stderr strings.Builder
	cmd.Stderr = &stderr

	stdout, err := cmd.Output()
	return strings.TrimSpace(string(stdout)), strings.TrimSpace(stderr.String()), err
}

/*
RunCommandCombinedOutput is a subset-wrapper for exec.Command, returning both stdout and stderr in one string.
*/
func RunCommandCombinedOutput(d models.Path, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	if !d.Empty() {
		cmd.Dir = d.String()
	}

	stdout, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(stdout)), err
}
