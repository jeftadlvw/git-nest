package utils

import (
	"errors"
	"os/exec"
	"strings"
)

func FindGitRoot() (string, error) {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
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
