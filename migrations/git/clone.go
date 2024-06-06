package git

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"golang.org/x/term"
	"os"
)

type Clone struct {
	Url          interfaces.Url
	Path         models.Path
	CloneDirName string
}

func (m Clone) Migrate() error {

	if !m.Path.Exists() {
		err := os.MkdirAll(m.Path.String(), os.ModePerm)
		if err != nil {
			return fmt.Errorf("internal error: could not create directory %s: %w", m.Path, err)
		}
	}

	var liveOutputFunc func(string) = nil
	terminalFd := 0
	var terminalWidth int

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
				if len(line) > terminalWidth {
					line = line[:utils.MaxInt(0, terminalWidth-6)] + "..."
				}

				_, _ = fmt.Fprintf(os.Stderr, "\r%*s", -terminalWidth, line)
			}
		}
	}

	err := utils.CloneGitRepository(m.Url.String(), m.Path, m.CloneDirName, liveOutputFunc)

	if term.IsTerminal(terminalFd) {
		_, _ = fmt.Fprintf(os.Stderr, "\r%*s", -terminalWidth, "")
		_, _ = fmt.Fprintf(os.Stderr, "\r")
	}

	if err != nil {
		return fmt.Errorf("error while cloning into %s: %s", m.Path.SJoin(m.CloneDirName), err)
	}

	return nil
}
