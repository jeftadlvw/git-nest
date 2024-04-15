package models

import (
	"fmt"
	"strings"
)

type SshUrl struct {
	Hostname string
	User     string
	Path     string
}

func NewSshUrl(hostname string, user string, path string) SshUrl {
	return SshUrl{hostname, user, path}
}

func (u *SshUrl) String() string {
	return fmt.Sprintf("ssh://%s@%s%s", u.User, u.Hostname, u.Path)
}

func (u *SshUrl) UnMarshalText(text []byte) error {
	sshUrl, err := SshUrlFromString(string(text))
	if err != nil {
		return err
	}
	*u = sshUrl
	return nil
}

func (u *SshUrl) MarshalText() (text []byte, err error) {
	return []byte(u.String()), nil
}

func SshUrlFromString(s string) (SshUrl, error) {
	u := SshUrl{}

	s = strings.TrimSpace(s)
	if s == "" {
		return u, fmt.Errorf("empty string")
	}

	// split :// for scheme --> check if scheme is correct (can be kept away)
	schemaSplit := strings.Split(s, "://")
	continueSplit := schemaSplit[0]
	if len(schemaSplit) == 2 {
		if schemaSplit[0] != "ssh" {
			return u, fmt.Errorf("unsupported protocol scheme %s, must be 'ssh'", schemaSplit[0])
		}
		continueSplit = schemaSplit[1]
	}

	// split at @ for [0]user and [1] host + path
	userHostSplit := strings.Split(continueSplit, "@")
	if len(userHostSplit) != 2 {
		return u, fmt.Errorf("ssh url must contain user@host")
	}

	u.User = userHostSplit[0]

	// split host at : for hostname and path
	hostPathSplit := strings.Split(userHostSplit[1], ":")
	if len(hostPathSplit) > 2 {
		return u, fmt.Errorf("ssh url must contain only one path")
	}

	u.Hostname = hostPathSplit[0]
	if len(hostPathSplit) == 2 {
		u.Path = hostPathSplit[1]
	}

	return NewSshUrl(u.Hostname, u.Path, u.User), nil
}
