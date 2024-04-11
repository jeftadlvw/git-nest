package models

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	NO_EXIST = iota
	FILE     = iota
	DIR      = iota
)

/*
Path is a string-typed type that abstracts path operations.
*/
type Path string

/*
IsFile returns whether this Path is an existing file.
*/
func (p *Path) IsFile() bool {
	return pathCheck(*p) == FILE
}

/*
IsDir returns whether this Path is an existing directory.
*/
func (p *Path) IsDir() bool {
	return pathCheck(*p) == DIR
}

/*
Exists returns whether this Path exists.
*/
func (p *Path) Exists() bool {
	return pathCheck(*p) != NO_EXIST
}

/*
Clean returns an OS-fitting, cleaned up copy of this Path, based on the stdlib filepath.Clean.
*/
func (p *Path) Clean() Path {
	return Path(filepath.Clean(strings.TrimSpace(string(*p))))
}

/*
Up returns a copy of this Path in the parent directory, based on the stdlib filepath.Dir.
*/
func (p *Path) Up() Path {
	return Path(filepath.Dir(string(*p)))
}

/*
Parent is a synonym for Up().
*/
func (p *Path) Parent() Path {
	return p.Up()
}

/*
Parts returns all single of the Path.
It uses filepath.Separator to split the path string.
*/
func (p *Path) Parts() []string {
	return strings.Split(string(*p), string(filepath.Separator))
}

/*
EmptyOrAtRoot returns whether this Path is an empty string, at a root-level-directory
or at its top-most parent directory of its original path.
*/
func (p *Path) EmptyOrAtRoot() bool {
	s := string(*p)
	return s == "/" || s == "." || s == ".." || s == ""
}

/*
Base returns the last element of this Path, based on the stdlib filepath.Base.
*/
func (p *Path) Base() Path {
	return Path(filepath.Base(string(*p)))
}

/*
Join returns a new Path with all passed subpaths joined together using filepath.Join.
*/
func (p *Path) Join(paths ...string) Path {
	return Path(filepath.Join(append([]string{string(*p)}, paths...)...))
}

/*
Contains returns whether the passed pattern exist within this Path's directory, based on the stdlib filepath.Glob.
*/
func (p *Path) Contains(pattern string) (bool, error) {

	if !p.IsDir() {
		return false, errors.New("path does not exist or is not a directory")
	}

	matches, err := filepath.Glob(filepath.Join(string(*p), pattern))
	if err != nil {
		return false, err
	}

	return len(matches) != 0, nil
}

/*
BContains is a wrapper for Path.Contains and only returns the success boolean value.
*/
func (p *Path) BContains(pattern string) bool {
	contains, _ := p.Contains(pattern)
	return contains
}

/*
String returns this Path as its base type, removing the need for specific casts.
*/
func (p *Path) String() string {
	return string(*p)
}

/*
UnmarshalText unmarshalls any byte array into a Path type.
This is done so that Path implements the encoding.TextUnmarshaler interface.
*/
func (p *Path) UnmarshalText(text []byte) error {
	*p = Path(strings.TrimSpace(string(text)))
	return nil
}

/*
MarshalText unmarshalls this Path into a string.
This is done so that Path implements the encoding.TextMarshaler interface.
*/
func (p *Path) MarshalText() (text []byte, err error) {
	return []byte(*p), nil
}

func pathCheck(p Path) int {
	fileInfo, err := os.Stat(string(p))
	if err != nil {
		if os.IsNotExist(err) {
			return NO_EXIST
		}
	}

	if fileInfo.IsDir() {
		return DIR
	}

	return FILE
}
