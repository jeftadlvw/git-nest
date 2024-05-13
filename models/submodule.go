package models

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models/urls"
	"net/url"
	"path/filepath"
	"strings"
)

type Submodule struct {
	Path Path
	Url  *urls.HttpUrl
	Ref  string
}

/*
Clean performs a data cleanup on this Submodule.
*/
func (s *Submodule) Clean() {
	s.Path = s.Path.Clean()
	s.Ref = strings.TrimSpace(s.Ref)
}

/*
RemoteIdentifier returns a string to uniquely identify a submodule's remote origin, incl. the reference.

Format: Submodule.Url@Submodule.Ref
*/
func (s *Submodule) RemoteIdentifier() string {
	s.Clean()

	hostPathConcat := s.Url.HostPathConcatStrict()
	if s.Ref != "" {
		hostPathConcat = hostPathConcat + "@" + s.Ref
	}

	return hostPathConcat
}

/*
Identifier returns a string to uniquely identify a submodule.

Format: Submodule.Url@Submodule.Ref>Submodule.Path
*/
func (s *Submodule) Identifier() string {
	s.Clean()

	pathSuffix := ""

	if s.Path.EmptyOrAtRoot() {
		pathSuffix = filepath.Base(s.Url.Path())
	} else {
		pathSuffix = s.Path.String()
	}

	return s.RemoteIdentifier() + ">" + pathSuffix
}

/*
String returns a string representation of this Submodule.
*/
func (s *Submodule) String() string {
	return fmt.Sprintf("Submodule %s", s.Identifier())
}

/*
Validate performs validation on this Submodule.
*/
func (s *Submodule) Validate() error {
	s.Clean()

	if s.Path.EmptyOrAtRoot() {
		return fmt.Errorf("submodule path must be set")
	}

	forbiddenCharacters := "!*"
	for _, char := range forbiddenCharacters {
		if strings.Contains(s.Path.String(), string(char)) {
			return fmt.Errorf("submodule path contains forbidden character '%c'", char)
		}
	}

	// url must be set
	if s.Url == nil || s.Url.String() == "" {
		return fmt.Errorf("submodule url is required")
	}

	if _, err := url.Parse(s.Url.String()); err != nil {
		return fmt.Errorf("submodule url is invalid")
	}

	// no whitespaces in ref
	if strings.Contains(s.Ref, " ") {
		return fmt.Errorf("submodule ref contains spaces (%s)", s.Ref)
	}

	return nil
}
