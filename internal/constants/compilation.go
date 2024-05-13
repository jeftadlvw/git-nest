package constants

import (
	"strconv"
	"strings"
)

/*
Variable which values get injected by the `go build` compiler and the git-nest toolchain.
*/
var (
	version                 string
	ref                     string
	compilationTimestampStr string
	ephemeralBuildStr       string
)

var (
	compilationTimestampInt = -2
)

/*
Version returns the binary version. '[ephemeral]' if binary is compiled and run with `go run`,
'[dev]' if `go build` was run by the git-nest toolchain, or a version string if the version
has been injected during `go build`.
*/
func Version() string {
	if version == "" {
		if EphemeralBuild() {
			return "[ephemeral]"
		}
		return "[dev]"
	}
	return version
}

/*
Ref returns the version control location from which the binary was build on. 'unset' in case
it was not injected during `go build`, else the injected value.
*/
func Ref() string {
	if ref == "" {
		return "unset"
	}
	return ref
}

/*
EphemeralBuild returns whether the compiled binary was created by `go run` or the git-nest toolchain.
The toolchain injects a value at compile time that toggles the binary to be non-ephemeral.
*/
func EphemeralBuild() bool {
	return strings.ToLower(ephemeralBuildStr) != "false" || ephemeralBuildStr != "0"
}

/*
CompilationTimestamp returns the binary's compilation time in form of a unix timestamp, which is
injected by the build toolchain. -1 if no compilation time is set.
*/
func CompilationTimestamp() int {

	if compilationTimestampInt == -2 {
		compilationTime, err := strconv.Atoi(compilationTimestampStr)
		if err != nil {
			compilationTimestampInt = -1
		} else {
			compilationTimestampInt = compilationTime
		}
	}

	return compilationTimestampInt
}
