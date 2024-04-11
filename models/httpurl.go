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

func NewHttpUrl(host string, port int, path string, secure bool) *HttpUrl {
	return &HttpUrl{host, port, path, secure}
}

func (h *HttpUrl) Host() string {
	return fmt.Sprintf("%s:%d", h.Hostname, h.Port)
}

func (h *HttpUrl) String() string {

	var (
		scheme string
		host   string
	)

	if h.Secure {
		scheme = "https"
	} else {
		scheme = "http"
	}

	if strings.Contains(h.Hostname, ":") {
		host = fmt.Sprintf("[%s]", h.Host())
	}

	if (!h.Secure && h.Port != 80) || (h.Secure && h.Port != 443) {
		host = fmt.Sprintf("%s:%d", host, h.Port)
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, h.Path)
}

func FromString(s string) (*HttpUrl, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	if !strings.Contains(u.Scheme, "http") {
		return nil, errors.New("scheme is not http")
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
			return nil, fmt.Errorf("invalid port number: %s", u.Port())
		}
	}

	return NewHttpUrl(u.Hostname(), port, u.Path, secure), nil

}

func (h *HttpUrl) UnMarshalText(text []byte) error {
	httpUrl, err := FromString(string(text))
	if err != nil {
		return err
	}
	*h = *httpUrl
	return nil
}

func (h *HttpUrl) MarshalText() (text []byte, err error) {
	return []byte(h.String()), nil
}
