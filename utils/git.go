package utils

import (
	"errors"
	"github.com/jeftadlvw/git-nest/models"
	"os/exec"
	"strings"
)

func GetGitRemoteUrl(d models.Path) (string, error) {
	path, err := RunCommand(d, "git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}

	pathStr := string(path)
	pathStr = strings.TrimSpace(pathStr)

	if strings.HasPrefix(pathStr, "fatal:") {
		return "", errors.New("git root not found")
	}

	return pathStr, nil
}

func GetGitVersion() (string, error) {
	version, err := exec.Command("git", "--version").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(version)), nil
}
