package utils

import (
	"github.com/jeftadlvw/git-nest/models"
	"os/exec"
	"strings"
)

func RunCommand(d models.Path, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	if !d.Empty() {
		cmd.Dir = d.String()
	}
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}
