package utils

import (
	"errors"
	"strings"
)

func StringInsert(original string, insert string, startDelimiter string, endDelimiter string) (string, error) {
	return stringInsert(original, insert, startDelimiter, endDelimiter, false)
}

func StringInsertAtFirst(original string, insert string, startDelimiter string, endDelimiter string) (string, error) {
	return stringInsert(original, insert, startDelimiter, endDelimiter, true)
}

func stringInsert(original string, insert string, startDelimiter string, endDelimiter string, ignoreMultipleOccurrences bool) (string, error) {

	if startDelimiter == "" {
		return "", errors.New("startDelimiter is empty")
	}

	if endDelimiter == "" {
		return "", errors.New("endDelimiter is empty")
	}

	startDeterminationCount := strings.Count(original, startDelimiter)
	endDeterminationCount := strings.Count(original, endDelimiter)

	if startDeterminationCount == 0 {
		return "", errors.New("cannot find starting delimiter")
	}

	if endDeterminationCount == 0 {
		return "", errors.New("cannot find ending delimiter")
	}

	if !ignoreMultipleOccurrences {
		if startDeterminationCount > 1 {
			return "", errors.New("start delimiter found multiple times")
		}

		if endDeterminationCount > 1 {
			return "", errors.New("end delimiter found multiple times")
		}
	}

	before := strings.SplitN(original, startDelimiter, 2)
	after := strings.SplitN(original, endDelimiter, 2)

	return before[0] + insert + after[len(after)-1], nil
}
