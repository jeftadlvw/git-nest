package utils

import (
	"fmt"
	"strings"
)

func StringInsert(original string, insert string, startDetermination string, endDetermination string) (string, error) {

	startDeterminationCount := strings.Count(original, startDetermination)
	endDeterminationCount := strings.Count(original, endDetermination)

	if startDeterminationCount == 0 {
		return "", fmt.Errorf("cannot to find starting determinator")
	}

	if endDeterminationCount == 0 {
		return "", fmt.Errorf("cannot to find ending determinator")
	}

	if startDeterminationCount > 1 {
		return "", fmt.Errorf("start determinator found multiple times")
	}

	if endDeterminationCount > 1 {
		return "", fmt.Errorf("end determinator found multiple times")
	}

	before := strings.Split(original, startDetermination)
	after := strings.Split(original, endDetermination)

	/*
		TODO: replace with regex-version
		// the regex here is untested
		var re = regexp.MustCompile(`(startDetermination).*(endDetermination)`)
		s := re.ReplaceAllString(sample, `${1}foo`)
	*/

	return before[0] + insert + after[len(after)-1], nil
}
