package urls

import (
	"fmt"
	"strings"
)

type SshUrl struct {
	Hostname string
	User     string
	Path     string
}

func (u *SshUrl) Clean() {
	u.Hostname = strings.TrimSpace(u.Hostname)
	u.User = strings.TrimSpace(u.User)

	u.Path = strings.TrimSpace(u.Path)
	u.Path = strings.Trim(u.Path, "/")

	if !strings.HasPrefix(u.Path, ":") {
		u.Path = "/" + u.Path
	}
}

func (u *SshUrl) IsEmpty() bool {
	u.Clean()
	return u.Hostname == "" || u.User == ""
}

func (u *SshUrl) HostPathConcat() string {
	if u.IsEmpty() {
		return ""
	}

	var path string

	if u.Path != "/" {
		path = u.Path
	}

	return fmt.Sprintf("%s%s", u.Hostname, path)
}

func (u *SshUrl) String() string {
	if u.IsEmpty() {
		return ""
	}

	return fmt.Sprintf("ssh://%s@%s", u.User, u.HostPathConcat())
}

func (u *SshUrl) UnmarshalText(text []byte) error {
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
			return SshUrl{}, fmt.Errorf("unsupported protocol scheme %s, must be 'ssh'", schemaSplit[0])
		}
		continueSplit = schemaSplit[1]
	}

	// split at @ for [0]user and [1] host + path
	userHostSplit := strings.Split(continueSplit, "@")
	if len(userHostSplit) != 2 {
		return u, fmt.Errorf("ssh urls must contain one user@host")
	}

	u.User = userHostSplit[0]

	// split host at : for hostname and :port/path
	if strings.HasPrefix(userHostSplit[1], ":") {
		return u, fmt.Errorf("hostname does not exists or starts with a colon")
	}

	hostPathColonSplit := strings.Split(userHostSplit[1], ":")
	if len(hostPathColonSplit) > 2 {
		return SshUrl{}, fmt.Errorf("ssh urls may contain only one port")
	}

	if len(hostPathColonSplit) == 2 {
		u.Hostname = hostPathColonSplit[0]
		u.Path = ":" + hostPathColonSplit[1]
	} else {
		// split host at / for hostname and /part
		hostPathSlashSplit := strings.SplitN(hostPathColonSplit[0], "/", 2)
		u.Hostname = hostPathSlashSplit[0]

		if len(hostPathSlashSplit) == 2 {
			u.Path = "/" + hostPathSlashSplit[1]
		}
	}

	if strings.Contains(u.Hostname, "/") {
		return SshUrl{}, fmt.Errorf("ssh hostname must not contain '/'")
	}

	return SshUrl{u.Hostname, u.User, u.Path}, nil
}
