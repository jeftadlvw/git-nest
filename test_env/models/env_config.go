package models

type EnvSettings struct {
	NoGit    bool
	EmptyGit bool
	Path     string
	Origin   string
	CloneDir string
	Ref      string
}
