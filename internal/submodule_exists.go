package internal

import (
	"errors"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"strings"
)

const (
	SUBMODULE_EXISTS_OK = iota
	SUBMODULE_EXISTS_UNDEFINED_REF
	SUBMODULE_EXISTS_ERR_NO_EXIST
	SUBMODULE_EXISTS_ERR_FILE
	SUBMODULE_EXISTS_ERR_NO_GIT
	SUBMODULE_EXISTS_ERR_REMOTE
	SUBMODULE_EXISTS_ERR_HEAD
)

type SubmoduleExistsMapValue struct {
	Status  int
	Payload string
	Error   error
}

/*
SubmoduleExists returns whether and in what state a possible submodule exists.
*/
func SubmoduleExists(s models.Submodule, root models.Path) (int, string, error) {
	submodulePath := root.Join(s.Path)

	if !submodulePath.Exists() {
		return SUBMODULE_EXISTS_ERR_NO_EXIST, "", nil
	}

	if !submodulePath.IsDir() {
		return SUBMODULE_EXISTS_ERR_FILE, "", nil
	}

	submoduleGitRemoteUrl, err := utils.GetGitRemoteUrl(submodulePath)
	if err != nil {
		return SUBMODULE_EXISTS_ERR_NO_GIT, "", nil
	}

	if submoduleGitRemoteUrl != s.Url.String() {
		return SUBMODULE_EXISTS_ERR_REMOTE, submoduleGitRemoteUrl, nil
	}

	if s.Ref == "" {
		return SUBMODULE_EXISTS_UNDEFINED_REF, "", nil
	}

	var (
		returnFlag    = SUBMODULE_EXISTS_OK
		returnPayload string
		returnErr     error
	)

	remoteRef, remoteRefAbbrev, err := utils.GetGitFetchHead(submodulePath)
	if err != nil {
		returnFlag = SUBMODULE_EXISTS_ERR_HEAD
		returnErr = err
		return returnFlag, returnPayload, returnErr
	}

	if len(remoteRef) == 0 {
		if !strings.HasPrefix(remoteRef, s.Ref) {
			returnFlag = SUBMODULE_EXISTS_ERR_HEAD
			returnPayload = remoteRef
		}
	} else {
		if s.Ref != "" && remoteRefAbbrev != s.Ref {
			returnFlag = SUBMODULE_EXISTS_ERR_HEAD
			returnPayload = remoteRefAbbrev
		}
	}

	return returnFlag, returnPayload, returnErr
}

/*
SubmodulesExist takes multiple submodules and verifies their existence in bulk
*/
func SubmodulesExist(submodules []models.Submodule, root models.Path) []SubmoduleExistsMapValue {

	var existMapping []SubmoduleExistsMapValue

	for _, submodule := range submodules {
		status, payload, err := SubmoduleExists(submodule, root)
		existMapping = append(existMapping, SubmoduleExistsMapValue{
			Status:  status,
			Payload: payload,
			Error:   err,
		})
	}

	return existMapping
}

/*
SubmoduleValid takes multiple submodules and verifies their existence in bulk
*/
func SubmoduleValid(submodule models.Submodule, root models.Path) bool {
	status, _, _ := SubmoduleExists(submodule, root)
	return SubmoduleStatusValid(status)
}

/*
ValidSubmodulesCount how many passed submodules are valid.
*/
func ValidSubmodulesCount(submodules []models.Submodule, root models.Path) int {
	valid := 0

	for _, submodule := range submodules {
		if SubmoduleValid(submodule, root) {
			valid++
		}
	}

	return valid
}

/*
SubmoduleStatusValid returns if a status belongs to being valid
*/
func SubmoduleStatusValid(status int) bool {
	return status == SUBMODULE_EXISTS_OK || status == SUBMODULE_EXISTS_UNDEFINED_REF
}

func FmtSubmoduleExistOutput(status int, payload string, err error) (string, error) {
	existStr := ""

	switch status {
	case SUBMODULE_EXISTS_OK:
		existStr = "ok"
	case SUBMODULE_EXISTS_UNDEFINED_REF:
		existStr = "ok, empty ref"
	case SUBMODULE_EXISTS_ERR_NO_EXIST:
		existStr = "no exist"
	case SUBMODULE_EXISTS_ERR_FILE:
		existStr = "error: path is a file"
	case SUBMODULE_EXISTS_ERR_NO_GIT:
		existStr = "error: git not installed"
	case SUBMODULE_EXISTS_ERR_REMOTE:
		existStr = "error: unequal remote urls: " + payload
	case SUBMODULE_EXISTS_ERR_HEAD:
		if err != nil {
			existStr = "error: unable to fetch HEAD: " + err.Error()
		} else {
			existStr = "error: unequal ref HEADs: " + payload
		}
	default:
		return "", errors.New("invalid exist state")
	}

	return existStr, nil
}
