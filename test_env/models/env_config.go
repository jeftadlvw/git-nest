package models

type EnvSettings struct {
	NoGit    bool
	EmptyGit bool
	Origin   string
	CloneDir string
	Ref      string
}
