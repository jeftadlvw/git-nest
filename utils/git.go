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
func CloneGitRepository(url string, p models.Path, cloneDirName string, liveOutput func(string)) error {
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

	outputBuilder := strings.Builder{}
	liveOutputFunc := func(line string) {
		outputBuilder.WriteString(line + "\n")

		if liveOutput != nil {
			liveOutput(line)
		}
	}

	err := RunCommandLiveOutputCombinedOutput(liveOutputFunc, p, "git", commandArgsArr...)

	output := outputBuilder.String()
	if err != nil {
		return fmt.Errorf("error running git clone: %w; output: %s", err, output)
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
GitCheckout changes a local repository's HEAD.
*/
func GitCheckout(repository models.Path, ref string) error {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return fmt.Errorf("ref cannot be blank")
	}

	if repository.Empty() {
		return errors.New("path is empty")
	}

	if !repository.Exists() {
		return fmt.Errorf("%s does not exist", repository)
	}

	if !repository.IsDir() {
		return fmt.Errorf("%s is not a directory", repository)
	}

	output, err := RunCommandCombinedOutput(repository, "git", "checkout", ref, "--progress")

	if strings.Contains(output, "fatal: not a git repository") {
		return fmt.Errorf("%s is not a git repository", repository)
	}

	if strings.Contains(output, "error: pathspec") {
		return fmt.Errorf("ref '%s' does not exist", ref)
	}

	if err != nil {
		return fmt.Errorf("error running git checkout: %w; output: %s", err, output)
	}

	return nil
}

/*
GitPull changes a local repository's HEAD.
*/
func GitPull(repository models.Path, liveOutput func(string)) error {
	if repository.Empty() {
		return errors.New("empty repository value")
	}

	if !repository.Exists() {
		return fmt.Errorf("%s does not exist", repository)
	}

	if !repository.IsDir() {
		return fmt.Errorf("%s is not a directory", repository)
	}

	outputBuilder := strings.Builder{}
	liveOutputFunc := func(line string) {
		outputBuilder.WriteString(line + "\n")

		if liveOutput != nil {
			liveOutput(line)
		}
	}

	err := RunCommandLiveOutputCombinedOutput(liveOutputFunc, repository, "git", "pull", "--progress")
	output := outputBuilder.String()

	if strings.Contains(output, "fatal: not a git repository") {
		return fmt.Errorf("%s is not a git repository", repository)
	}

	if strings.Contains(output, "There is no tracking information for the current branch.") {
		return fmt.Errorf("repository has no configured remote")
	}

	if strings.Contains(output, "error: The following untracked working tree files would be overwritten by merge:") {
		return fmt.Errorf("repository contains untracked contents")
	}

	if strings.Contains(output, "fatal: Need to specify how to reconcile divergent branches.") {
		return fmt.Errorf("repository does not know how to merge")
	}

	if strings.Contains(output, "fatal: Not possible to fast-forward, aborting") {
		return fmt.Errorf("unable to fast-forward")
	}

	if strings.Contains(output, "Merge conflict") {
		abortOutput, abortErr := RunCommandCombinedOutput(repository, "git", "merge", "--abort")
		if abortErr != nil {
			return fmt.Errorf("error aborting merge: %w; output: %s", abortErr, abortOutput)
		}

		return fmt.Errorf("pull would cause merge conflict")
	}

	if err != nil {
		return fmt.Errorf("error running git checkout: %w; output: %s", err, output)
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
		return "", fmt.Errorf("error running git rev-parse: %w; output: %s", err, path)
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

	remoteUrl, err := RunCommandCombinedOutput(d, "git", "config", "--get", "remote.origin.url")
	if err != nil {
		return "", fmt.Errorf("error running git config: %w; output: %s", err, remoteUrl)
	}

	if strings.HasPrefix(remoteUrl, "fatal:") {
		return "", errors.New("git root not found")
	}

	return remoteUrl, nil
}

/*
GetGitFetchHead retrieves the current HEAD of a local repository.
Returns long commit hash, commit abbreviation and/or error.
*/
func GetGitFetchHead(d models.Path) (string, string, error) {
	if d.Empty() {
		return "", "", errors.New("path to repository may not be empty")
	}

	longHead, err := RunCommandCombinedOutput(d, "git", "rev-parse", "--verify", "HEAD")
	if err != nil {
		return "", "", fmt.Errorf("error running git rev-parse: %w; output: %s", err, longHead)
	}

	if strings.HasPrefix(longHead, "fatal: not a git repository") {
		return "", "", errors.New("no git repository")
	}

	abbrevHead, err := RunCommandCombinedOutput(d, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", "", fmt.Errorf("error running git rev-parse: %w; output: %s", err, abbrevHead)
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

/*
GetGitHasUntrackedChanges returns whether a local repository has uncommitted changes.
In case of any errors, true is returned.
*/
func GetGitHasUntrackedChanges(d models.Path) (bool, error) {
	if d.Empty() {
		return true, errors.New("path to repository may not be empty")
	}

	out, err := RunCommandCombinedOutput(d, "git", "status", "--porcelain=v1")
	if err != nil {
		return true, err
	}

	// if not a repository, then return false
	if strings.HasPrefix(out, "fatal: not a git repository") {
		return false, nil
	}

	// in case another error occurs, return it
	if strings.HasPrefix(out, "fatal:") {
		return true, fmt.Errorf("git error: %s", out)
	}

	return strings.TrimSpace(out) != "", nil
}

/*
GetGitHasUnpublishedChanges returns whether a local repository has unpushed commits.
In case of any errors, true is returned.
*/
func GetGitHasUnpublishedChanges(d models.Path) (bool, error) {
	if d.Empty() {
		return true, errors.New("path to repository may not be empty")
	}

	out, err := RunCommandCombinedOutput(d, "git", "status")
	if err != nil {
		return true, err
	}

	// if not a repository, then return false
	if strings.HasPrefix(out, "fatal: not a git repository") {
		return false, nil
	}

	// in case another error occurs, return it
	if strings.HasPrefix(out, "fatal:") {
		return true, fmt.Errorf("git error: %s", out)
	}

	return strings.Contains(out, "(use \"git push\" to publish your local commits)"), nil
}
