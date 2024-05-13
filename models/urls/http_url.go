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
		HostnameS contains the domain name or ip address.
	*/
	HostnameS string

	/*
		Port contains the host's port.
	*/
	Port int

	/*
		PathS contains any further information encoded in the url path.
	*/
	PathS string

	/*
		Secure is a flag that defines whether this HttpUrl should be treated as an encrypted connection.
	*/
	Secure bool
}

/*
Clean cleans up struct values.
*/
func (u *HttpUrl) Clean() {
	u.HostnameS = strings.TrimSpace(u.HostnameS)

	u.PathS = strings.TrimSpace(u.PathS)
	u.PathS = "/" + strings.Trim(u.PathS, "/")
}

/*
Host returns the HostnameS and Port in a formatted matter.
*/
func (u *HttpUrl) Host(forcePort bool) string {
	var host string

	if strings.Contains(u.HostnameS, ":") {
		host = fmt.Sprintf("[%s]", u.HostnameS)
	} else {
		host = u.HostnameS
	}

	if forcePort || ((!u.Secure && u.Port != 80) || (u.Secure && u.Port != 443)) {
		host = fmt.Sprintf("%s:%d", host, u.Port)
	}

	return host
}

/*
IsEmpty returns whether HostnameS is empty or not. It calls Clean() beforehand.
*/
func (u *HttpUrl) IsEmpty() bool {
	u.Clean()
	return u.HostnameS == ""
}

/*
Hostname returns this Url's hostname.
*/
func (u *HttpUrl) Hostname() string {
	return u.HostnameS
}

/*
Path returns this Url's path.
*/
func (u *HttpUrl) Path() string {
	return u.PathS
}

/*
HostPathConcat returns the url without the protocol.
*/
func (u *HttpUrl) HostPathConcat() string {
	return u.hostPathConcat(false)
}

/*
HostPathConcatStrict returns the url forcing Port to be concatenated as well.
*/
func (u *HttpUrl) HostPathConcatStrict() string {
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
	if u.PathS == "/" {
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

/*
Validate validates this HttpUrl.
*/
func (u *HttpUrl) Validate() error {
	_, err := HttpUrlFromString(u.String())
	return err
}

func (u *HttpUrl) hostPathConcat(forcePort bool) string {
	if u.IsEmpty() {
		return ""
	}

	var path string

	if u.PathS != "/" {
		path = u.PathS
	}

	return fmt.Sprintf("%s%s", u.Host(forcePort), path)
}
