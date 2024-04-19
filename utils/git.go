package utils

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
CloneGitRepository clones a remote git repository.
*/
func CloneGitRepository(url string, p models.Path, cloneDirName string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return fmt.Errorf("git repository url is empty")
	}

	cloneDirName = strings.TrimSpace(cloneDirName)

	if cloneDirName != "" && filepath.Base(cloneDirName) != cloneDirName {
		return fmt.Errorf("repository clone directory name may not be path")
	}

	if !p.Exists() {
		return fmt.Errorf("%s does not exist", p)
	}

	commandArgsArr := []string{"clone", "--progress", url}
	if cloneDirName != "" {
		commandArgsArr = append(commandArgsArr, cloneDirName)
	}

	output, err := RunCommandCombinedOutput(p, "git", commandArgsArr...)
	if err != nil {
		return fmt.Errorf("error running git clone: %s; output: %s", err, output)
	}

	if strings.Contains(output, "ERROR: Repository not found.") {
		return fmt.Errorf("remote repository %s does not exist", url)
	}

	if strings.Contains(output, "fatal: destination path") {
		return fmt.Errorf("destination path already exists")
	}

	return nil
}

/*
ChangeGitHead changes a local repository's HEAD.
*/
func ChangeGitHead(repository models.Path, head string) error {
	head = strings.TrimSpace(head)
	if head == "" {
		return fmt.Errorf("head is empty")
	}

	if !repository.Exists() {
		return fmt.Errorf("%s does not exist", repository)
	}

	output, err := RunCommandCombinedOutput(repository, "git", "checkout", head, "--progress")
	if err != nil {
		return fmt.Errorf("error running git checkout: %s; output: %s", err, output)
	}

	if strings.Contains(output, "fatal: not a git repository") {
		return fmt.Errorf("%s is not a git repository", repository)
	}

	if strings.Contains(output, "error: pathspec") {
		return fmt.Errorf("head '%s' does not exist", head)
	}

	return nil
}

/*
GetGitRootDirectory retrieves the root of a git directory tree.
*/
func GetGitRootDirectory(d models.Path) (string, error) {
	if d.Empty() {
		return "", errors.New("path to repository may not be empty")
	}

	path, err := RunCommandCombinedOutput(d, "git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("error running git rev-parse: %s; output: %s", err, path)
	}

	if strings.HasPrefix(path, "fatal:") {
		return "", errors.New("git root not found")
	}

	return path, nil
}

/*
GetGitRemoteUrl retrieves the remote url from a git directory tree.
*/
func GetGitRemoteUrl(d models.Path) (string, error) {
	if d.Empty() {
		return "", errors.New("path to repository may not be empty")
	}

	path, err := RunCommandCombinedOutput(d, "git", "config", "--get", "remote.origin.url")
	if err != nil {
		return "", fmt.Errorf("error running git config: %s; output: %s", err, path)
	}

	if strings.HasPrefix(path, "fatal:") {
		return "", errors.New("git root not found")
	}

	return path, nil
}

/*
GetGitFetchHead retrieves the current HEAD of a local repository.
*/
func GetGitFetchHead(d models.Path) (string, string, error) {
	if d.Empty() {
		return "", "", errors.New("path to repository may not be empty")
	}

	longHead, err := RunCommandCombinedOutput(d, "git", "rev-parse", "--verify", "HEAD")
	if err != nil {
		return "", "", fmt.Errorf("error running git rev-parse: %s; output: %s", err, longHead)
	}

	if strings.HasPrefix(longHead, "fatal:") {
		return "", "", errors.New("git root not found")
	}

	abbrevHead, err := RunCommandCombinedOutput(d, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", "", fmt.Errorf("error running git rev-parse: %s; output: %s", err, abbrevHead)
	}

	if abbrevHead != "HEAD" {
		return longHead, abbrevHead, nil
	}

	return longHead, "", nil
}

/*
GetGitVersion retrieves the git version installed in the current environment. Can also be used to check if git is installed.
*/
func GetGitVersion() (string, error) {
	version, err := exec.Command("git", "--version").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(version)), nil
}
