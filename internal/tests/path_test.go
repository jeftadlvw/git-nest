package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/models"
	"testing"
)

func TestPathContainsUp(t *testing.T) {
	tests := []struct {
		path     models.Path
		expected bool
	}{
		{"/path/to/file", false},
		{"/path/../file", false},
		{"/path/../../file", false},
		{"/path/foo/../../file", false},
		{"../file", true},
		{"./../file", true},
		{"../../file", true},
	}

	for index, test := range tests {
		t.Run(fmt.Sprintf("TestPathContainsUp-%d", index+1), func(t *testing.T) {
			result := internal.PathContainsUp(test.path)
			if result != test.expected {
				t.Errorf("returned %v, expected %v", result, test.expected)
			}
		})
	}
}

func TestPathOutsideRoot(t *testing.T) {
	root := models.Path("/path/to/root")
	tests := []struct {
		rootPath models.Path
		path     models.Path
		expected bool
	}{
		{"/path/to/root", "/path/to/root/file", false},
		{"/path/to/root", "/path/to/root/../file", true},
		{"/path/to/root", "/other/path", true},
		{"/path/to/root", "/path/to/other/root/file", true},
	}

	for index, test := range tests {
		t.Run(fmt.Sprintf("TestPathOutsideRoot-%d", index+1), func(t *testing.T) {
			result := internal.PathOutsideRoot(root, test.path)
			if result != test.expected {
				t.Errorf("returned %v, expected %v", result, test.expected)
			}
		})
	}
}

func TestPathRelativeToRootWithJoinedOriginIfNotAbs(t *testing.T) {
	tests := []struct {
		root     models.Path
		origin   models.Path
		path     models.Path
		expected models.Path
		err      bool
	}{
		{"/path/to/root", "", "file", "file", true},
		{"/path/to/root", "relative/path", "file", "file", true},
		{"/path/to/root", "foo", "file", "file", true},

		{"/path/to/root", "/origin", "file", "../../../origin/file", false},
		{"/path/to/root", "does not matter", "/path/to/root/file", "file", false},
		{"/path/to/root", "", "/path/to/file", "../file", false},
		{"/path/to/root", "", "/path/../file", "../../../file", false},

		{"/path/to/root", "/path/to/root/foo", "file", "foo/file", false},
		{"/path/to/root", "/path/to/root/foo", "bar/file", "foo/bar/file", false},
		{"/path/to/root", "/path/to/root/foo", "../bar/file", "bar/file", false},
		{"/path/to/root", "/path/to/root/foo", "../bar/../file", "file", false},
		{"/path/to/root", "/path/to/root/foo", "../bar/./../file", "file", false},
		{"/path/to/root", "/path/to/root/foo", ".", "foo", false},
		{"/path/to/root", "/path/to/root/foo", "", "foo", false},

		{"/path/to/root", "/path/to", "file", "../file", false},
		{"/path/to/root", "/path/to", "bar/file", "../bar/file", false},
		{"/path/to/root", "/path/to", "../bar/file", "../../bar/file", false},
		{"/path/to/root", "/path/to", "../bar/../file", "../../file", false},
		{"/path/to/root", "/path/to", "../bar/./../file", "../../file", false},
		{"/path/to/root", "/path/to", ".", "..", false},
		{"/path/to/root", "/path/to", "", "..", false},

		{"/path/to/root", "/path/to/baz", "file", "../baz/file", false},
		{"/path/to/root", "/path/to/baz", "bar/file", "../baz/bar/file", false},
		{"/path/to/root", "/path/to/baz", "../bar/file", "../bar/file", false},
		{"/path/to/root", "/path/to/baz", "../bar/../file", "../file", false},
		{"/path/to/root", "/path/to/baz", "../bar/./../file", "../file", false},
		{"/path/to/root", "/path/to/baz", ".", "../baz", false},
		{"/path/to/root", "/path/to/baz", "", "../baz", false},
	}

	for index, tc := range tests {
		t.Run(fmt.Sprintf("TestPathRelativeToRootWithJoinedOriginIfNotAbs-%d", index+1), func(t *testing.T) {
			result, err := internal.PathRelativeToRootWithJoinedOriginIfNotAbs(tc.root, tc.origin, tc.path)
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !tc.err && result != tc.expected {
				t.Errorf("returned %s, expected %s", result, tc.expected)
			}
		})
	}
}
