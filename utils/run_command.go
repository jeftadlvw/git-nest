package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"os/exec"
	"strings"
	"sync"
)

/*
RunCommand is a subset-wrapper for exec.Command, providing separate return values for stdout and stderr.
*/
func RunCommand(d models.Path, command string, args ...string) (string, string, error) {
	var stdout strings.Builder
	var stderr strings.Builder

	generateStringBuildFunc := func(sb *strings.Builder) func(s string) {
		return func(s string) {
			sb.WriteString(s + "\n")
		}
	}

	err := RunCommandLiveOutput(generateStringBuildFunc(&stdout), generateStringBuildFunc(&stderr), d, command, args...)
	if err != nil {
		return "", "", err
	}

	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}

/*
RunCommandCombinedOutput is a subset-wrapper for exec.Command, returning both stdout and stderr in one string.
*/
func RunCommandCombinedOutput(d models.Path, command string, args ...string) (string, error) {
	var stdout strings.Builder

	fStdout := func(s string) {
		stdout.WriteString(s + "\n")
	}

	err := RunCommandLiveOutputCombinedOutput(fStdout, d, command, args...)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(stdout.String()), err
}

/*
RunCommandLiveOutput is a wrapper for exec.Command that takes a callback function for updates in both stdout and stderr streams.
*/
func RunCommandLiveOutput(fStdout func(string), fStderr func(string), wd models.Path, command string, args ...string) error {
	return runCommand(fStdout, fStderr, wd, command, args...)
}

/*
RunCommandLiveOutputCombinedOutput is a wrapper for exec.Command that takes a callback function for updates in a combined stdout and stderr stream.
*/
func RunCommandLiveOutputCombinedOutput(fStdout func(string), wd models.Path, command string, args ...string) error {
	var mutex = &sync.Mutex{}

	fStdoutWrapper := func(s string) {
		mutex.Lock()
		fStdout(s)
		mutex.Unlock()
	}

	return runCommand(fStdoutWrapper, fStdoutWrapper, wd, command, args...)
}

func runCommand(fStdout func(string), fStderr func(string), wd models.Path, command string, args ...string) error {

	// configure command
	cmd := exec.Command(command, args...)
	addEnglishLocaleEnv(cmd)
	if !wd.Empty() {
		cmd.Dir = wd.String()
	}

	// obtain pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("could not obtain stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("could not obtain stderr pipe: %w", err)
	}

	// Create readers for stdout and stderr.
	stdoutScanner := bufio.NewScanner(stdoutPipe)
	stdoutScanner.Split(scanLines)

	stderrScanner := bufio.NewScanner(stderrPipe)
	stderrScanner.Split(scanLines)

	// Function to read and print each line from a given Reader.
	readAndPrint := func(scanner *bufio.Scanner, callback func(string)) {
		for scanner.Scan() {
			t := scanner.Text()
			callback(t)
		}
	}

	// Read and print stdout and stderr concurrently.
	go readAndPrint(stdoutScanner, fStdout)
	go readAndPrint(stderrScanner, fStderr)

	// start and wait for command to finish
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("could not start command: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// scanLines is a split function for a [Scanner] that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is either a '\r' or a '\n.
// The last non-empty line of input will be returned even if it has no
// newline.
//
// This function is inspired by https://stackoverflow.com/a/41433698
// and based on by bufio.ScanLines.
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		if data[i] == '\n' {
			// We have a line terminated by single newline.
			return i + 1, data[0:i], nil
		}

		advance = i + 1
		if len(data) > i+1 && data[i+1] == '\n' {
			advance += 1
		}
		return advance, data[0:i], nil
	}

	/// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

/*
addEnglishLocaleEnv adds an environment variables that causes some programs to force their output language to english.
Works on unix only, but is added for every platform.
*/
func addEnglishLocaleEnv(cmd *exec.Cmd) {
	cmd.Env = append(cmd.Env, "LANGUAGE=en_US.UTF-8")
}
