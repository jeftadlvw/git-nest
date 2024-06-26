package models

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	NO_EXIST = iota
	FILE     = iota
	DIR      = iota
)

var unixStringReplacePattern = regexp.MustCompile(`/+`)

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
	cleanedPath := filepath.Clean(strings.TrimSpace(string(*p)))
	cleanedPath = strings.ReplaceAll(cleanedPath, strings.Repeat(string(filepath.Separator), 2), string(filepath.Separator))
	return Path(cleanedPath)
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
Relative returns the relative path from another Path to this Path, based on the stdlib filepath.Rel.
*/
func (p *Path) Relative(o Path) (Path, error) {
	rp, err := filepath.Rel(p.String(), o.String())
	return Path(rp), err
}

/*
EmptyOrAtRoot returns whether this Path is an empty string, at a root-level-directory
or at its top-most parent directory of its original path.
*/
func (p *Path) EmptyOrAtRoot() bool {
	return p.Empty() || p.AtRoot()
}

/*
Empty returns whether this Path is an empty string.
*/
func (p *Path) Empty() bool {
	return strings.TrimSpace(string(*p)) == ""
}

/*
AtRoot returns whether this Path is at a root-level-directory or at its top-most
parent directory of its original path.
*/
func (p *Path) AtRoot() bool {
	s := string(*p)
	return s == "/" || s == "." || s == ".."
}

/*
Base returns the last element of this Path, based on the stdlib filepath.Base.
*/
func (p *Path) Base() string {
	return filepath.Base(string(*p))
}

/*
Join returns a new Path with all passed Path structs joined together using filepath.Join.
If strings should be joined on this Path, use SJoin.
*/
func (p *Path) Join(paths ...Path) Path {
	pathsStr := make([]string, len(paths))
	for i, path := range paths {
		pathsStr[i] = string(path)
	}

	return Path(filepath.Join(append([]string{string(*p)}, pathsStr...)...))
}

/*
SJoin returns a new Path with all passed strings joined together using filepath.Join.
*/
func (p *Path) SJoin(paths ...string) Path {
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
Equals returns whether this and another path are the same
*/
func (p *Path) Equals(other Path) bool {

	pLocal, otherLocal := p.Clean(), other.Clean()

	// lowercase the strings and compare them
	thisLowerCase := strings.ToLower(string(pLocal))
	otherLowerCase := strings.ToLower(string(otherLocal))

	// if not equal in lowercase, then they are not the same path
	if thisLowerCase != otherLowerCase {
		return false
	}

	// if equal in lowercase, proceed to check if path is on a
	// case-sensitive filesystem or not
	caseSensitive, err := IsOnCaseSensitiveFilesystem(pLocal)
	if err != nil {
		// return false in case of an error
		return false
	}

	// if case-sensitive, compare both original strings
	if caseSensitive {
		return pLocal == otherLocal
	}

	// if case-insensitive, return true
	return true
}

/*
String returns this Path as its base type, removing the need for specific casts.
*/
func (p *Path) String() string {
	if p.EmptyOrAtRoot() {
		return ""
	}
	return string(p.Clean())
}

/*
UnixString ensures the path is formatted as a UNIX path (using '/' as path separators).
*/
func (p *Path) UnixString() string {
	str := p.String()

	str = strings.ReplaceAll(str, string(filepath.Separator), "/")
	str = strings.ReplaceAll(str, "\\", "/")
	str = unixStringReplacePattern.ReplaceAllString(str, "/")

	return str
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

/*
IsOnCaseSensitiveFilesystem returns whether a given path is on a case-sensitive filesystem.
*/
func IsOnCaseSensitiveFilesystem(p Path) (bool, error) {
	alt := p.Parent()
	alt = alt.SJoin(flipCase(p.Base()))

	// get file stat for passed file
	pathInfo, err := os.Stat(string(p))
	if err != nil {
		return false, err
	}

	// if file does not exist, assume to be on case-sensitive filesystem
	altInfo, err := os.Stat(string(alt))
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}

		return false, err
	}

	// if both file exist, check if they are the same
	return !os.SameFile(pathInfo, altInfo), nil
}

func pathCheck(p Path) int {
	fileInfo, err := os.Stat(string(p))
	if err != nil {
		if os.IsNotExist(err) {
			return NO_EXIST
		}
	}

	if fileInfo == nil {
		return NO_EXIST
	}

	if fileInfo.IsDir() {
		return DIR
	}

	return FILE
}

func flipCase(s string) string {
	if s == "" {
		return s
	}
	firstChar := string(s[0])
	if strings.ToLower(firstChar) == firstChar {
		return strings.ToUpper(firstChar) + s[1:]
	}
	return strings.ToLower(firstChar) + s[1:]
}
