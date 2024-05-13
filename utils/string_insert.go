package utils

import (
	"errors"
	"strings"
)

/*
StringInsert inserts a string between the start and end delimiter strings.
Use StringInsertAtFirst if multiple occurrences of both delimiters are allowed in the input string.
*/
func StringInsert(original string, insert string, startDelimiter string, endDelimiter string) (string, error) {
	return stringInsert(original, insert, startDelimiter, endDelimiter, false)
}

/*
StringInsertAtFirst is the same as StringInsert, but allows the occurrence of multiple delimiter string in the input string
by only replacing the first start-end combination.
*/
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
