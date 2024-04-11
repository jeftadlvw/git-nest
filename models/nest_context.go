package models

/*
NestContext bundles all relevant information of a project that utilizes git-nest.
*/
type NestContext struct {
	/*
		Config is the configuration of the project.
	*/
	Config NestConfig

	/*
		WorkingDirectory contains the Path to the working directory the binary was executed from.
	*/
	WorkingDirectory Path

	/*
		ProjectRoot is a Path to the project's root directory.

		A directory counts as project root if a directory contains a `nestmodules.toml` file directly at the
		directory or in a `.Config` subdirectory. If no configuration file could be found, the directory tree
		is traversed up to find the next possible parent project. If the current directory is not part of a
		git-nest project, the string is set to the current working directory.
	*/
	ProjectRoot Path

	/*
		ConfigFile is a Path that points to the project's `nestmodules.toml`.

		If no configuration files has found, it points to `[ProjectRoot]/nestmodules.toml`.
	*/
	ConfigFile Path
}