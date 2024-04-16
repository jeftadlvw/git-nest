package tests

import (
	"github.com/jeftadlvw/git-nest/models/urls"
	"testing"
)

func TestHttpUrl_Clean(t *testing.T) {
	tests := []struct {
		url      urls.HttpUrl
		expected urls.HttpUrl
	}{
		{
			url:      urls.HttpUrl{Hostname: "  example.com  ", Path: "  /path/to/resource  "},
			expected: urls.HttpUrl{Hostname: "example.com", Path: "/path/to/resource"}},
		{
			url:      urls.HttpUrl{Hostname: "example.com", Path: ""},
			expected: urls.HttpUrl{Hostname: "example.com", Path: "/"},
		},
		{
			url:      urls.HttpUrl{Hostname: "", Path: "path/to/resource"},
			expected: urls.HttpUrl{Hostname: "", Path: "/path/to/resource"},
		},
		{
			url:      urls.HttpUrl{Hostname: "example.com", Path: "/path/to/resource/"},
			expected: urls.HttpUrl{Hostname: "example.com", Path: "/path/to/resource"},
		},
		{
			url:      urls.HttpUrl{Hostname: "example.com", Path: "/path/to/resource"},
			expected: urls.HttpUrl{Hostname: "example.com", Path: "/path/to/resource"},
		},
	}

	for _, tc := range tests {
		tc.url.Clean()
		if tc.url.Hostname != tc.expected.Hostname {
			t.Errorf("Expected hostname %s in %v, got: %s", tc.expected.Hostname, tc.url, tc.url.Hostname)
		}

		if tc.url.Path != tc.expected.Path {
			t.Errorf("Expected path %s in %v, got: %s", tc.expected.Path, tc.url, tc.url.Path)
		}
	}
}

func TestHttpUrlHost(t *testing.T) {
	cases := []struct {
		url      urls.HttpUrl
		expected string
	}{
		{urls.HttpUrl{"example.com", -1, "", false}, "example.com:-1"},
		{urls.HttpUrl{"example.com", 80, "", false}, "example.com"},
		{urls.HttpUrl{"example.com", 80, "", true}, "example.com:80"},
		{urls.HttpUrl{"example.com", 443, "/", true}, "example.com"},
		{urls.HttpUrl{"example.com", 443, "/", false}, "example.com:443"},
		{urls.HttpUrl{"example.com", 8080, "/", false}, "example.com:8080"},
	}

	for _, tc := range cases {
		result := tc.url.Host(false)
		if result != tc.expected {
			t.Errorf("Host() for %v returned >%s<, expected >%s<", tc.url, result, tc.expected)
		}
	}
}

func TestHttpUrlHostForcePort(t *testing.T) {
	cases := []struct {
		url      urls.HttpUrl
		expected string
	}{
		{urls.HttpUrl{"example.com", -1, "", false}, "example.com:-1"},
		{urls.HttpUrl{"example.com", 80, "", false}, "example.com:80"},
		{urls.HttpUrl{"example.com", 80, "", true}, "example.com:80"},
		{urls.HttpUrl{"example.com", 443, "/", true}, "example.com:443"},
		{urls.HttpUrl{"example.com", 443, "/", false}, "example.com:443"},
		{urls.HttpUrl{"example.com", 8080, "/", false}, "example.com:8080"},
	}

	for _, tc := range cases {
		result := tc.url.Host(true)
		if result != tc.expected {
			t.Errorf("Host() for %v returned >%s<, expected >%s<", tc.url, result, tc.expected)
		}
	}
}

func TestHttpUrlIsEmpty(t *testing.T) {
	cases := []struct {
		url      urls.HttpUrl
		expected bool
	}{
		{urls.HttpUrl{"", -1, "", false}, true},
		{urls.HttpUrl{"  \n", 80, "", false}, true},
		{urls.HttpUrl{"  example.com  ", 80, "", false}, false},
		{urls.HttpUrl{"example.com", 80, "", false}, false},
	}

	for _, tc := range cases {
		result := tc.url.IsEmpty()
		if result != tc.expected {
			t.Errorf("Host() for %v returned %t, expected %t", tc.url, result, tc.expected)
		}
	}
}

func TestHttpUrlHostPathConcat(t *testing.T) {
	cases := []struct {
		url      urls.HttpUrl
		expected string
	}{
		{urls.HttpUrl{"", 80, "/", false}, ""},
		{urls.HttpUrl{"example.com", 80, "/", false}, "example.com"},
		{urls.HttpUrl{"example.com", 443, "/", true}, "example.com"},
		{urls.HttpUrl{"example.com", 8080, "/", false}, "example.com:8080"},
		{urls.HttpUrl{"example.com", 8080, "/path", false}, "example.com:8080/path"},
		{urls.HttpUrl{"example.com", 443, "/path/", true}, "example.com/path"},
	}

	for _, tc := range cases {
		result := tc.url.HostPathConcat(false)
		if result != tc.expected {
			t.Errorf("HostPathConcat() for %v returned %s, expected %s", tc.url, result, tc.expected)
		}
	}
}

func TestHttpUrlString(t *testing.T) {
	cases := []struct {
		url      urls.HttpUrl
		expected string
	}{
		{urls.HttpUrl{"", 80, "/", false}, ""},
		{urls.HttpUrl{"example.com", 80, "/", false}, "http://example.com/"},
		{urls.HttpUrl{"example.com", 443, "/", true}, "https://example.com/"},
		{urls.HttpUrl{"example.com", 8080, "/", false}, "http://example.com:8080/"},
		{urls.HttpUrl{"example.com", 8080, "/path", false}, "http://example.com:8080/path"},
		{urls.HttpUrl{"example.com", 443, "/path/", true}, "https://example.com/path"},
	}

	for _, tc := range cases {
		result := tc.url.String()
		if result != tc.expected {
			t.Errorf("String() for %v returned %s, expected %s", tc.url, result, tc.expected)
		}
	}
}

func TestHttpUrlUnMarshalText(t *testing.T) {
	cases := []struct {
		input    string
		expected urls.HttpUrl
		err      bool
	}{
		{"http://example.com", urls.HttpUrl{"example.com", 80, "", false}, false},
		{"https://example.com:8080/path", urls.HttpUrl{"example.com", 8080, "/path", true}, false},
		{"invalid_url", urls.HttpUrl{}, true},
	}

	for _, tc := range cases {
		var url urls.HttpUrl
		err := url.UnmarshalText([]byte(tc.input))
		if tc.err && err == nil {
			t.Errorf("UnmarshalText() for %s returned no error, expected error", tc.input)
		}
		if !tc.err && err != nil {
			t.Errorf("UnmarshalText() for %s returned error: %s", tc.input, err)
		}
		if url != tc.expected {
			t.Errorf("UnmarshalText() for %s returned %v, expected %v", tc.input, url, tc.expected)
		}
	}
}

func TestHttpUrlMarshalText(t *testing.T) {
	cases := []struct {
		url      urls.HttpUrl
		expected string
	}{
		{urls.HttpUrl{"example.com", 80, "/", false}, "http://example.com/"},
		{urls.HttpUrl{"example.com", 80, "/", true}, "https://example.com:80/"},
		{urls.HttpUrl{"example.com", 8080, "/path", false}, "http://example.com:8080/path"},
		{urls.HttpUrl{"example.com", 443, "/path", true}, "https://example.com/path"},
		{urls.HttpUrl{"example.com", 443, "/path", false}, "http://example.com:443/path"},
	}

	for _, tc := range cases {
		result, err := tc.url.MarshalText()
		if err != nil {
			t.Errorf("MarshalText() for %v returned error: %s", tc.url, err)
		}
		if string(result) != tc.expected {
			t.Errorf("MarshalText() for %v returned %s, expected %s", tc.url, result, tc.expected)
		}
	}
}

func TestHttpUrlFromString(t *testing.T) {
	cases := []struct {
		input    string
		expected urls.HttpUrl
		err      bool
	}{
		{"http://example.com", urls.HttpUrl{"example.com", 80, "", false}, false},
		{"https://example.com:8080/path", urls.HttpUrl{"example.com", 8080, "/path", true}, false},
		{"invalid_url", urls.HttpUrl{}, true},
		{"ftp://example.com", urls.HttpUrl{}, true},
		{"http://example.com:abc", urls.HttpUrl{}, true},
		{"http://example.com:99999", urls.HttpUrl{}, true},
	}

	for _, tc := range cases {
		result, err := urls.HttpUrlFromString(tc.input)
		if tc.err && err == nil {
			t.Errorf("HttpUrlFromString() for %s returned no error, expected error", tc.input)
		}
		if !tc.err && err != nil {
			t.Errorf("HttpUrlFromString() for %s returned error: %s", tc.input, err)
		}
		if result != tc.expected {
			t.Errorf("HttpUrlFromString() for %s returned %v, expected %v", tc.input, result, tc.expected)
		}
	}
}
