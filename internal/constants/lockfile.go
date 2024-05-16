package constants

/*
LockFileName contains the path string to the project-local lockfile.
*/
const LockFileName = "~git-nest.lock"

/*
LockFileContents contains the contents of the lockfile in case someone opens it.
*/
const LockFileContents = `git-nest lockfile
Manual removal could lead to data loss in the git-nest configuration.
Do not remove except you know what you're' doing!\n`
