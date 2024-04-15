package models

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type Submodule struct {
	Path   Path
	Url    HttpUrl
	Ref    string
	Exists bool
}

/*
RemoteIdentifier returns a string to uniquely identify a submodule's remote origin.

Format: Submodule.Url@Submodule.Ref
*/
func (s *Submodule) RemoteIdentifier() string {
	identifier := s.Url.Host()
	if s.Ref != "" {
		identifier = identifier + "@" + s.Ref
	}

	return identifier
}

/*
Identifier returns a string to uniquely identify a submodule.

Format: Submodule.Url@Submodule.Ref>Submodule.Path
*/
func (s *Submodule) Identifier() string {
	return s.Identifier() + ">" + s.Path.String()
}

/*
CleanUp performs a data cleanup on this Submodule.
*/
func (s *Submodule) CleanUp() {
	s.Path = s.Path.Clean()
	s.Ref = strings.TrimSpace(s.Ref)
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

	// cleanup structure before validating
	s.CleanUp()

	// no empty values for s.Path and s.Url
	if s.Path.String() == "" {
		return fmt.Errorf("submodule path is required")
	}

	if _, err := url.Parse(s.Url.String()); err != nil {
		return fmt.Errorf("submodule url is invalid")
	}

	// s.Path may not escape project root directory (so no ../)
	if strings.Contains(string(s.Path), ".."+string(filepath.Separator)) {
		return fmt.Errorf("submodule path escapes root directory (%s)", s.Path)
	}

	// TODO validate s.Ref: may be empty, else no spaces
	if strings.Contains(s.Ref, " ") {
		return fmt.Errorf("submodule ref contains spaces (%s)", s.Ref)
	}

	return nil
}
