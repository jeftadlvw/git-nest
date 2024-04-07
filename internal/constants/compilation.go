package constants

import (
	"strconv"
)

var (
	version                 string
	refHash                 string
	compilationTimestampStr string
	ephemeralBuildStr       string
)

var (
	compilationTimestampInt = -2
)

func Version() string {
	if version == "" {
		if EphemeralBuild() {
			return "[ephemeral]"
		}
		return "[dev]"
	}
	return version
}

func RefHash() string {
	if refHash == "" {
		return "unset"
	}
	return refHash
}

func EphemeralBuild() bool {
	return ephemeralBuildStr != "false"
}

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
