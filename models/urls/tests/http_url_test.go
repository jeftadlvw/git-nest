package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models/urls"
	"testing"
)

func TestHttpUrlClean(t *testing.T) {
	tests := []struct {
		url      urls.HttpUrl
		expected urls.HttpUrl
	}{
		{
			url:      urls.HttpUrl{HostnameS: "  example.com  ", PathS: "  /path/to/resource  "},
			expected: urls.HttpUrl{HostnameS: "example.com", PathS: "/path/to/resource"}},
		{
			url:      urls.HttpUrl{HostnameS: "example.com", PathS: ""},
			expected: urls.HttpUrl{HostnameS: "example.com", PathS: "/"},
		},
		{
			url:      urls.HttpUrl{HostnameS: "", PathS: "path/to/resource"},
			expected: urls.HttpUrl{HostnameS: "", PathS: "/path/to/resource"},
		},
		{
			url:      urls.HttpUrl{HostnameS: "example.com", PathS: "/path/to/resource/"},
			expected: urls.HttpUrl{HostnameS: "example.com", PathS: "/path/to/resource"},
		},
		{
			url:      urls.HttpUrl{HostnameS: "example.com", PathS: "/path/to/resource"},
			expected: urls.HttpUrl{HostnameS: "example.com", PathS: "/path/to/resource"},
		},
	}

	for index, tc := range tests {
		t.Run(fmt.Sprintf("TestHttpUrlClean-%d", index+1), func(t *testing.T) {
			tc.url.Clean()
			if tc.url.Hostname() != tc.expected.Hostname() {
				t.Fatalf("unexpected hostname %s, expected %s", tc.url.Hostname(), tc.expected.Hostname())
			}

			if tc.url.Path() != tc.expected.Path() {
				t.Fatalf("unexpected path %s, expected %s", tc.url.Path(), tc.expected.Path())
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestHttpUrlHost-%d", index+1), func(t *testing.T) {
			result := tc.url.Host(false)
			if result != tc.expected {
				t.Fatalf("returned >%s<, expected >%s<", result, tc.expected)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestHttpUrlHostForcePort-%d", index+1), func(t *testing.T) {
			result := tc.url.Host(true)
			if result != tc.expected {
				t.Fatalf("returned >%s<, expected >%s<", result, tc.expected)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestHttpUrlIsEmpty-%d", index+1), func(t *testing.T) {
			result := tc.url.IsEmpty()
			if result != tc.expected {
				t.Fatalf("returned %t, expected %t", result, tc.expected)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestHttpUrlHostPathConcat-%d", index+1), func(t *testing.T) {
			result := tc.url.HostPathConcat()
			if result != tc.expected {
				t.Fatalf("returned %s, expected %s", result, tc.expected)
			}
		})
	}
}

func TestHttpUrlHostPathConcatForcePort(t *testing.T) {
	cases := []struct {
		url      urls.HttpUrl
		expected string
	}{
		{urls.HttpUrl{"", 80, "/", false}, ""},
		{urls.HttpUrl{"example.com", 80, "/", false}, "example.com:80"},
		{urls.HttpUrl{"example.com", 443, "/", true}, "example.com:443"},
		{urls.HttpUrl{"example.com", 8080, "/", false}, "example.com:8080"},
		{urls.HttpUrl{"example.com", 8080, "/path", false}, "example.com:8080/path"},
		{urls.HttpUrl{"example.com", 443, "/path/", true}, "example.com:443/path"},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestHttpUrlHostPathConcatForcePort-%d", index+1), func(t *testing.T) {
			result := tc.url.HostPathConcatStrict()
			if result != tc.expected {
				t.Fatalf("returned %s, expected %s", result, tc.expected)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestHttpUrlString-%d", index+1), func(t *testing.T) {
			result := tc.url.String()
			if result != tc.expected {
				t.Fatalf("returned %s, expected %s", result, tc.expected)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestHttpUrlUnMarshalText-%d", index+1), func(t *testing.T) {
			var url urls.HttpUrl
			err := url.UnmarshalText([]byte(tc.input))
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if url != tc.expected {
				t.Fatalf("returned %v, expected %v", url, tc.expected)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestHttpUrlMarshalText-%d", index+1), func(t *testing.T) {
			result, err := tc.url.MarshalText()
			if err != nil {
				t.Fatalf("MarshalText() for %v returned error: %s", tc.url, err)
			}
			if string(result) != tc.expected {
				t.Errorf("MarshalText() for %v returned %s, expected %s", tc.url, result, tc.expected)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestHttpUrlFromString-%d", index+1), func(t *testing.T) {
			result, err := urls.HttpUrlFromString(tc.input)
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if result != tc.expected {
				t.Errorf("returned %v, expected %v", result, tc.expected)
			}
		})
	}
}
