package internal

import (
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"strings"
)

type SubmoduleExist struct {
	Flag    int
	Payload string
	Error   error
}

const (
	SUBMODULE_EXISTS_OK           = iota
	SUBMODULE_EXISTS_ERR_NO_EXIST = iota
	SUBMODULE_EXISTS_ERR_FILE     = iota
	SUBMODULE_EXISTS_ERR_NO_GIT   = iota
	SUBMODULE_EXISTS_ERR_REMOTE   = iota
	SUBMODULE_EXISTS_ERR_HEAD     = iota
)

func SubmoduleExists(s models.Submodule, root models.Path) SubmoduleExist {
	returnExists := SubmoduleExist{SUBMODULE_EXISTS_OK, "", nil}
	submodulePath := root.Join(s.Path)

	if !submodulePath.Exists() {
		returnExists.Flag = SUBMODULE_EXISTS_ERR_NO_EXIST
		return returnExists
	}

	if !submodulePath.IsDir() {
		returnExists.Flag = SUBMODULE_EXISTS_ERR_FILE
		return returnExists
	}

	submoduleGitRemoteUrl, err := utils.GetGitRemoteUrl(submodulePath)
	if err != nil {
		returnExists.Flag = SUBMODULE_EXISTS_ERR_NO_GIT
		returnExists.Error = err
		return returnExists
	}

	if submoduleGitRemoteUrl != s.Url.String() {
		returnExists.Flag = SUBMODULE_EXISTS_ERR_REMOTE
		returnExists.Payload = submoduleGitRemoteUrl
	}

	remoteRef, remoteRefAbbrev, err := utils.GetGitFetchHead(submodulePath)
	if err != nil {
		returnExists.Flag = SUBMODULE_EXISTS_ERR_HEAD
		returnExists.Error = err
	}

	if len(remoteRef) == 0 {
		if !strings.HasPrefix(remoteRef, s.Ref) {
			returnExists.Flag = SUBMODULE_EXISTS_ERR_HEAD
			returnExists.Payload = remoteRef
		}
	} else {
		if remoteRefAbbrev != s.Ref {
			returnExists.Flag = SUBMODULE_EXISTS_ERR_HEAD
			returnExists.Payload = remoteRefAbbrev
		}
	}

	return returnExists
}
