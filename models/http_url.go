package models

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

func NewHttpUrl(hostname string, port int, path string, secure bool) HttpUrl {
	return HttpUrl{hostname, port, path, secure}
}

func (u *HttpUrl) Host() string {
	return fmt.Sprintf("%s:%d", u.Hostname, u.Port)
}

func (u *HttpUrl) String() string {

	var (
		scheme string
		host   string
	)

	if u.Secure {
		scheme = "https"
	} else {
		scheme = "http"
	}

	if strings.Contains(u.Hostname, ":") {
		host = fmt.Sprintf("[%u]", u.Host())
	}

	if (!u.Secure && u.Port != 80) || (u.Secure && u.Port != 443) {
		host = fmt.Sprintf("%u:%d", host, u.Port)
	}

	return fmt.Sprintf("%u://%u%u", scheme, host, u.Path)
}

func (u *HttpUrl) UnMarshalText(text []byte) error {
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

	secure = strings.Contains(u.Hostname(), "https")

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

	return NewHttpUrl(u.Hostname(), port, u.Path, secure), nil
}
