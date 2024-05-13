package tests

import (
	"fmt"
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

	for index, tc := range tests {
		t.Run(fmt.Sprintf("TestSshUrlClean-%d", index+1), func(t *testing.T) {
			tc.url.Clean()
			if tc.url.HostnameS != tc.expected.HostnameS {
				t.Fatalf("unexpected hostname >%s<, expected >%s<", tc.url.HostnameS, tc.expected.HostnameS)
			}

			if tc.url.User != tc.expected.User {
				t.Fatalf("unexpected user >%s<, expected >%s<", tc.url.User, tc.expected.User)
			}

			if tc.url.PathS != tc.expected.PathS {
				t.Fatalf("unexpected path >%s<, expected >%s<", tc.url.PathS, tc.expected.PathS)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestSshUrlIsEmpty-%d", index+1), func(t *testing.T) {
			result := tc.url.IsEmpty()
			if result != tc.empty {
				t.Fatalf("returned %t, expected %t", result, tc.empty)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestSshUrlHostPathConcat-%d", index+1), func(t *testing.T) {
			result := tc.url.HostPathConcat()
			if result != tc.expected {
				t.Fatalf("returned %s, expected %s", result, tc.expected)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestSshUrlString-%d", index+1), func(t *testing.T) {
			result := tc.url.String()
			if result != tc.expected {
				t.Fatalf("returned >%s<, expected >%s<", result, tc.expected)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestSshUrlUnmarshalText-%d", index+1), func(t *testing.T) {
			var url urls.SshUrl
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

func TestSshUrlMarshalText(t *testing.T) {
	cases := []struct {
		url      urls.SshUrl
		expected string
	}{
		{urls.SshUrl{"example.com", "user", "/path"}, "ssh://user@example.com:/path"},
		{urls.SshUrl{"example.com", "user", ""}, "ssh://user@example.com"},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestSshUrlMarshalText-%d", index+1), func(t *testing.T) {
			result, err := tc.url.MarshalText()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if string(result) != tc.expected {
				t.Fatalf("returned %s, expected %s", result, tc.expected)
			}
		})
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

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestSshUrlFromString-%d", index+1), func(t *testing.T) {
			result, err := urls.SshUrlFromString(tc.input)
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if result != tc.expected {
				t.Fatalf("returned %v, expected %v", result, tc.expected)
			}
		})
	}
}
