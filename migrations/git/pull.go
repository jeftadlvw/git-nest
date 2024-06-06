package git

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"golang.org/x/term"
	"os"
)

type Pull struct {
	Path models.Path
}

func (m Pull) Migrate() error {

	if !m.Path.Exists() {
		return errors.New("path does not exist")
	}

	var liveOutputFunc func(string) = nil
	terminalFd := 0
	var terminalWidth int
	baseOutput := fmt.Sprintf("%s", m.Path)

	if term.IsTerminal(terminalFd) {
		// execute once, but allows for early returns using break
		localTerminalWidth, _, err := term.GetSize(terminalFd)
		terminalWidth = localTerminalWidth

		for range 1 {
			if err != nil {
				break
			}

			if terminalWidth == 0 {
				break
			}

			liveOutputFunc = func(line string) {
				// shorten string and add ellipsis if it's shorter than the terminal width

				line = fmt.Sprintf("%s: %s", baseOutput, line)

				if len(line) > terminalWidth {
					line = line[:utils.MaxInt(0, terminalWidth-6)] + "..."
				}

				// pad content on right
				_, _ = fmt.Fprintf(os.Stderr, "\r%*s", -terminalWidth, line)
			}
		}
	}

	fmt.Printf("%s: busy", baseOutput)
	err := utils.GitPull(m.Path, liveOutputFunc)

	if liveOutputFunc != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\r%*s", -terminalWidth, "")
		_, _ = fmt.Fprintf(os.Stderr, "\r")
	}

	if err != nil {
		fmt.Printf("\r%s: error: %s", baseOutput, err)
		return fmt.Errorf("could not perform pull operation at %s: %w", m.Path, err)
	}

	fmt.Printf("\r%s: done.\n", baseOutput)
	return nil
}
