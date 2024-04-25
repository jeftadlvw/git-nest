package internal

import (
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"strings"
)

const (
	SUBMODULE_EXISTS_OK           = iota
	SUBMODULE_EXISTS_ERR_NO_EXIST = iota
	SUBMODULE_EXISTS_ERR_FILE     = iota
	SUBMODULE_EXISTS_ERR_NO_GIT   = iota
	SUBMODULE_EXISTS_ERR_REMOTE   = iota
	SUBMODULE_EXISTS_ERR_HEAD     = iota
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
		if remoteRefAbbrev != s.Ref {
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
