package urls

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

/*
HttpUrl represents a regular http url.
*/
type HttpUrl struct {
	/*
		Hostname contains the domain name or ip address.
	*/
	Hostname string

	/*
		Port contains the host's port.
	*/
	Port int

	/*
		Path contains any further information encoded in the url path.
	*/
	Path string

	/*
		Secure is a flag that defines whether this HttpUrl should be treated as an encrypted connection.
	*/
	Secure bool
}

/*
Clean cleans up struct values.
*/
func (u *HttpUrl) Clean() {
	u.Hostname = strings.TrimSpace(u.Hostname)

	u.Path = strings.TrimSpace(u.Path)
	u.Path = "/" + strings.Trim(u.Path, "/")
}

/*
Host returns the Hostname and Port in a formatted matter.
*/
func (u *HttpUrl) Host(forcePort bool) string {
	var host string

	if strings.Contains(u.Hostname, ":") {
		host = fmt.Sprintf("[%s]", u.Hostname)
	} else {
		host = u.Hostname
	}

	if forcePort || ((!u.Secure && u.Port != 80) || (u.Secure && u.Port != 443)) {
		host = fmt.Sprintf("%s:%d", host, u.Port)
	}

	return host
}

/*
IsEmpty returns whether Hostname is empty or not. It calls Clean() beforehand.
*/
func (u *HttpUrl) IsEmpty() bool {
	u.Clean()
	return u.Hostname == ""
}

/*
HostPathConcat returns the url without the protocol.
*/
func (u *HttpUrl) HostPathConcat() string {
	return u.hostPathConcat(false)
}

/*
HostPathConcatForcePort returns the url forcing Port to be concatenated as well.
*/
func (u *HttpUrl) HostPathConcatForcePort() string {
	return u.hostPathConcat(true)
}

/*
String returns this HttpUrl back as a usable url.
*/
func (u *HttpUrl) String() string {

	if u.IsEmpty() {
		return ""
	}

	u.Clean()

	var scheme string

	if u.Secure {
		scheme = "https"
	} else {
		scheme = "http"
	}

	hostPathConcat := u.HostPathConcat()
	if u.Path == "/" {
		hostPathConcat = hostPathConcat + "/"
	}

	return fmt.Sprintf("%s://%s", scheme, hostPathConcat)
}

/*
UnmarshalText implements the encoding.TextUnmarshaler interface.
*/
func (u *HttpUrl) UnmarshalText(text []byte) error {
	httpUrl, err := HttpUrlFromString(string(text))
	if err != nil {
		return err
	}
	*u = httpUrl
	return nil
}

/*
MarshalText implements the encoding.TextMarshaler interface.
*/
func (u *HttpUrl) MarshalText() (text []byte, err error) {
	return []byte(u.String()), nil
}

/*
HttpUrlFromString creates a HttpUrl from a string, returning an error if validation fails.
*/
func HttpUrlFromString(s string) (HttpUrl, error) {
	emptyUrl := HttpUrl{}

	u, err := url.Parse(s)
	if err != nil {
		return emptyUrl, err
	}

	if !strings.Contains(u.Scheme, "http") {
		return emptyUrl, errors.New("scheme is not http-based")
	}

	var (
		secure bool
		port   int
	)

	secure = u.Scheme == "https"

	if u.Port() == "" {
		if secure {
			port = 443
		} else {
			port = 80
		}
	} else {
		port, err = strconv.Atoi(u.Port())
		if err != nil || port < 0 || port >= 65535 {
			return emptyUrl, fmt.Errorf("invalid port number: %s", u.Port())
		}
	}

	return HttpUrl{u.Hostname(), port, u.Path, secure}, nil
}

func (u *HttpUrl) hostPathConcat(forcePort bool) string {
	if u.IsEmpty() {
		return ""
	}

	var path string

	if u.Path != "/" {
		path = u.Path
	}

	return fmt.Sprintf("%s%s", u.Host(forcePort), path)
}
