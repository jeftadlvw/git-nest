package utils

import (
	"github.com/jeftadlvw/git-nest/models"
	"os/exec"
	"strings"
)

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

func RunCommandCombinedOutput(d models.Path, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	if !d.Empty() {
		cmd.Dir = d.String()
	}

	stdout, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(stdout)), err
}
