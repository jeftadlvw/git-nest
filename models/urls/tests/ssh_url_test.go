package tests

import (
	"github.com/jeftadlvw/git-nest/models/urls"
	"testing"
)

func TestSshUrlClean(t *testing.T) {
	tests := []struct {
		url      urls.SshUrl
		expected urls.SshUrl
	}{
		{
			url:      urls.SshUrl{HostnameS: "example.com", User: "", PathS: ""},
			expected: urls.SshUrl{HostnameS: "example.com", User: "", PathS: ""},
		},
		{
			url:      urls.SshUrl{HostnameS: "  example.com  ", User: "  hello  ", PathS: "  /path/to/resource  "},
			expected: urls.SshUrl{HostnameS: "example.com", User: "hello", PathS: "/path/to/resource"}},
		{
			url:      urls.SshUrl{HostnameS: "", User: "", PathS: "path/to/resource"},
			expected: urls.SshUrl{HostnameS: "", User: "", PathS: "path/to/resource"},
		},
		{
			url:      urls.SshUrl{HostnameS: "example.com", User: "user", PathS: "/path/to/resource/"},
			expected: urls.SshUrl{HostnameS: "example.com", User: "user", PathS: "/path/to/resource"},
		},
	}

	for _, tc := range tests {
		tc.url.Clean()
		if tc.url.HostnameS != tc.expected.HostnameS {
			t.Errorf("Expected hostname >%s< in %v, got: >%s<", tc.expected.HostnameS, tc.url, tc.url.HostnameS)
		}

		if tc.url.User != tc.expected.User {
			t.Errorf("Expected user >%s< in %v, got: >%s<", tc.expected.User, tc.url, tc.url.PathS)
		}

		if tc.url.PathS != tc.expected.PathS {
			t.Errorf("Expected path >%s< in %v, got: >%s<", tc.expected.PathS, tc.url, tc.url.PathS)
		}
	}
}

func TestSshUrlIsEmpty(t *testing.T) {
	cases := []struct {
		url   urls.SshUrl
		empty bool
	}{
		{urls.SshUrl{"", "", "/path"}, true},
		{urls.SshUrl{"example.com", "", "/path"}, true},
		{urls.SshUrl{"", "user", "/path"}, true},
		{urls.SshUrl{"example.com", "user", ""}, false},
	}

	for _, tc := range cases {
		result := tc.url.IsEmpty()
		if result != tc.empty {
			t.Errorf("String() for %v returned %t, expected %t", tc.url, result, tc.empty)
		}
	}
}

func TestSshUrlHostPathConcat(t *testing.T) {
	cases := []struct {
		url      urls.SshUrl
		expected string
	}{
		{urls.SshUrl{"", "", "/"}, ""},
		{urls.SshUrl{"example.com", "", "/"}, ""},
		{urls.SshUrl{"example.com", "user", ""}, "example.com"},
		{urls.SshUrl{"example.com", "user", "/"}, "example.com:/"},
		{urls.SshUrl{"example.com", "user", "path"}, "example.com:path"},
		{urls.SshUrl{"example.com", "user", "/path"}, "example.com:/path"},
	}

	for _, tc := range cases {
		result := tc.url.HostPathConcat()
		if result != tc.expected {
			t.Errorf("HostPathConcat() for %v returned %s, expected %s", tc.url, result, tc.expected)
		}
	}
}

func TestSshUrlString(t *testing.T) {
	cases := []struct {
		url      urls.SshUrl
		expected string
	}{
		{urls.SshUrl{"", "", "/path"}, ""},
		{urls.SshUrl{"example.com", "user", "/path"}, "ssh://user@example.com:/path"},
		{urls.SshUrl{"example.com", "user", "/"}, "ssh://user@example.com:/"},
		{urls.SshUrl{"example.com", "user", ""}, "ssh://user@example.com"},
	}

	for _, tc := range cases {
		result := tc.url.String()
		if result != tc.expected {
			t.Errorf("String() for %v returned >%s<, expected >%s<", tc.url, result, tc.expected)
		}
	}
}

func TestSshUrlUnmarshalText(t *testing.T) {
	cases := []struct {
		input    string
		expected urls.SshUrl
		err      bool
	}{
		{"user@example.com:path", urls.SshUrl{"example.com", "user", "path"}, false},
		{"user@example.com:/path", urls.SshUrl{"example.com", "user", "/path"}, false},
		{"ssh://user@example.com:port/path", urls.SshUrl{"example.com", "user", "port/path"}, false},
		{"ssh://user@example.com/path", urls.SshUrl{}, true},
		{"invalid_url", urls.SshUrl{}, true},
		{"", urls.SshUrl{}, true},
		{"ssh://user_example.com", urls.SshUrl{}, true},
		{"ssh://:path", urls.SshUrl{}, true},
		{"ssh://:", urls.SshUrl{}, true},
		{"ssh://user@example.com:path1:path2", urls.SshUrl{}, true},
	}

	for _, tc := range cases {
		var url urls.SshUrl
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

func TestSshUrlMarshalText(t *testing.T) {
	cases := []struct {
		url      urls.SshUrl
		expected string
	}{
		{urls.SshUrl{"example.com", "user", "/path"}, "ssh://user@example.com:/path"},
		{urls.SshUrl{"example.com", "user", ""}, "ssh://user@example.com"},
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

func TestSshUrlFromString(t *testing.T) {
	cases := []struct {
		input    string
		expected urls.SshUrl
		err      bool
	}{
		{"user@example.com:path", urls.SshUrl{"example.com", "user", "path"}, false},
		{"user@example.com:/path", urls.SshUrl{"example.com", "user", "/path"}, false},
		{"ssh://user@example.com:port/path", urls.SshUrl{"example.com", "user", "port/path"}, false},
		{"ssh://user@example.com/path", urls.SshUrl{}, true},
		{"invalid_url", urls.SshUrl{}, true},
		{"", urls.SshUrl{}, true},
		{"ssh://user_example.com", urls.SshUrl{}, true},
		{"ssh://:path", urls.SshUrl{}, true},
		{"ssh://:", urls.SshUrl{}, true},
		{"ssh://user@example.com:path1:path2", urls.SshUrl{}, true},
	}

	for _, tc := range cases {
		result, err := urls.SshUrlFromString(tc.input)
		if tc.err && err == nil {
			t.Errorf("SshUrlFromString() for %s returned no error, expected error", tc.input)
		}
		if !tc.err && err != nil {
			t.Errorf("SshUrlFromString() for %s returned error: %s", tc.input, err)
		}
		if result != tc.expected {
			t.Errorf("SshUrlFromString() for %s returned %v, expected %v", tc.input, result, tc.expected)
		}
	}
}
