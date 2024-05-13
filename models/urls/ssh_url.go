package urls

import (
	"fmt"
	"strings"
)

/*
SshUrl represents a regular ssh url.
*/
type SshUrl struct {
	/*
		HostnameS contains the domain name or ip address.
	*/
	HostnameS string

	/*
		User contains the username that want to connect to the host.
	*/
	User string

	/*
		PathS contains anything else specifying the connection.
	*/
	PathS string
}

/*
Clean cleans up struct values.
*/
func (u *SshUrl) Clean() {
	u.HostnameS = strings.TrimSpace(u.HostnameS)
	u.User = strings.TrimSpace(u.User)

	u.PathS = strings.TrimSpace(u.PathS)
	if u.PathS != "/" {
		u.PathS = strings.TrimSuffix(u.PathS, "/")
	}
	u.PathS = strings.TrimPrefix(u.PathS, ":")
}

/*
IsEmpty returns whether HostnameS or User is empty or not. It calls Clean() beforehand.
*/
func (u *SshUrl) IsEmpty() bool {
	u.Clean()
	return u.HostnameS == "" || u.User == ""
}

/*
Hostname returns this Url's hostname.
*/
func (u *SshUrl) Hostname() string {
	return u.HostnameS
}

/*
Path returns this Url's path.
*/
func (u *SshUrl) Path() string {
	return u.PathS
}

/*
HostPathConcat returns the url without the protocol.
*/
func (u *SshUrl) HostPathConcat() string {
	if u.IsEmpty() {
		return ""
	}

	var path string

	if u.PathS != "" {
		path = ":" + u.PathS
	}

	return fmt.Sprintf("%s%s", u.HostnameS, path)
}

/*
HostPathConcatStrict returns the url forcing Port to be concatenated as well.
*/
func (u *SshUrl) HostPathConcatStrict() string {
	return u.HostPathConcat()
}

/*
String returns this SshUrl back as a usable url.
*/
func (u *SshUrl) String() string {
	if u.IsEmpty() {
		return ""
	}

	return fmt.Sprintf("ssh://%s@%s", u.User, u.HostPathConcat())
}

/*
UnmarshalText implements the encoding.TextUnmarshaler interface.
*/
func (u *SshUrl) UnmarshalText(text []byte) error {
	sshUrl, err := SshUrlFromString(string(text))
	if err != nil {
		return err
	}
	*u = sshUrl
	return nil
}

/*
MarshalText implements the encoding.TextMarshaler interface.
*/
func (u *SshUrl) MarshalText() (text []byte, err error) {
	return []byte(u.String()), nil
}

/*
SshUrlFromString creates a SshUrl from a string, returning an error if validation fails.
*/
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

	// split host at : for hostname and :path
	if strings.HasPrefix(userHostSplit[1], ":") {
		return u, fmt.Errorf("hostname does not exists or starts with a colon")
	}

	hostPathColonSplit := strings.Split(userHostSplit[1], ":")
	if len(hostPathColonSplit) > 2 {
		return SshUrl{}, fmt.Errorf("ssh urls may contain only one port")
	}

	u.HostnameS = hostPathColonSplit[0]
	if len(hostPathColonSplit) == 2 {
		u.PathS = hostPathColonSplit[1]
	}

	if strings.Contains(u.HostnameS, "/") {
		return SshUrl{}, fmt.Errorf("ssh hostname must not contain '/'")
	}

	return SshUrl{u.HostnameS, u.User, u.PathS}, nil
}
