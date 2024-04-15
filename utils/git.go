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

	if strings.HasPrefix(path, "fatal:") {
		return "", errors.New("git root not found")
	}

	return path, nil
}

func GetGitFetchHead(d models.Path) ([]string, error) {
	longHead, err := RunCommand(d, "git", "rev-parse", "--verify", "HEAD")
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(longHead, "fatal:") {
		return nil, errors.New("git root not found")
	}

	abbrevHead, err := RunCommand(d, "git", "rev-parse", "--abbref-rev", "HEAD")
	if err != nil {
		return nil, err
	}

	returnArr := []string{longHead}

	if abbrevHead != "HEAD" {
		returnArr = append(returnArr, longHead)
		return returnArr, nil
	}

	return returnArr, nil
}

func GetGitVersion() (string, error) {
	version, err := exec.Command("git", "--version").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(version)), nil
}
