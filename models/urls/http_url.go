package urls

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type HttpUrl struct {
	Hostname string
	Port     int
	Path     string
	Secure   bool
}

func (u *HttpUrl) Clean() {
	u.Hostname = strings.TrimSpace(u.Hostname)

	u.Path = strings.TrimSpace(u.Path)
	u.Path = "/" + strings.Trim(u.Path, "/")
}

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

func (u *HttpUrl) IsEmpty() bool {
	u.Clean()
	return u.Hostname == ""
}

func (u *HttpUrl) HostPathConcat() string {
	return u.hostPathConcat(false)
}

func (u *HttpUrl) HostPathConcatForcePort() string {
	return u.hostPathConcat(true)
}

func (u *HttpUrl) String() string {

	if u.IsEmpty() {
		return ""
	}

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

func (u *HttpUrl) UnmarshalText(text []byte) error {
	httpUrl, err := HttpUrlFromString(string(text))
	if err != nil {
		return err
	}
	*u = httpUrl
	return nil
}

func (u *HttpUrl) MarshalText() (text []byte, err error) {
	return []byte(u.String()), nil
}

func HttpUrlFromString(s string) (HttpUrl, error) {
	emptyUrl := HttpUrl{}

	u, err := url.Parse(s)
	if err != nil {
		return emptyUrl, err
	}

	if !strings.Contains(u.Scheme, "http") {
		return emptyUrl, errors.New("scheme is not http")
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
