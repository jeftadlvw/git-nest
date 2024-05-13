package tests

import (
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"testing"
)

func TestSubmoduleCleanUp(t *testing.T) {
	submodule := models.Submodule{
		Path: "/valid/../path",
		Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
		Ref:  "  main  ",
	}

	submodule.Clean()

	if submodule.Path != "/path" {
		t.Errorf("Expected cleaned up path to be 'path', got %s", submodule.Path)
	}

	if submodule.Ref != "main" {
		t.Errorf("Expected cleaned up ref to be main, got %s", submodule.Ref)
	}
}

func TestSubmoduleRemoteIdentifier(t *testing.T) {
	tests := []struct {
		submodule models.Submodule
		expected  string
	}{
		{
			submodule: models.Submodule{
				Url: &urls.HttpUrl{"example.com", 443, "repository", true},
			},
			expected: "example.com:443/repository",
		},
		{
			submodule: models.Submodule{
				Url: &urls.HttpUrl{"another-example.com", 8080, "another-repository", true},
				Ref: "main",
			},
			expected: "another-example.com:8080/another-repository@main",
		},
	}

	for _, test := range tests {
		actual := test.submodule.RemoteIdentifier()
		if actual != test.expected {
			t.Errorf("Expected %s, but got %s", test.expected, actual)
		}
	}
}

func TestSubmoduleIdentifier(t *testing.T) {
	tests := []struct {
		submodule models.Submodule
		expected  string
	}{
		{
			submodule: models.Submodule{
				Path: "   ",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			expected: "example.com:443/repository@main>repository",
		},
		{
			submodule: models.Submodule{
				Path: "/src/module",
				Url:  &urls.HttpUrl{"example.com", 8080, "another-repository", true},
				Ref:  "dev",
			},
			expected: "example.com:8080/another-repository@dev>/src/module",
		},
		{
			submodule: models.Submodule{
				Path: "src/module",
				Url:  &urls.HttpUrl{"example.com", 8080, "another-repository", true},
				Ref:  "dev",
			},
			expected: "example.com:8080/another-repository@dev>src/module",
		},
	}

	for _, test := range tests {
		actual := test.submodule.Identifier()
		if actual != test.expected {
			t.Errorf("Expected %s, but got %s", test.expected, actual)
		}
	}
}

func TestSubmoduleString(t *testing.T) {
	tests := []struct {
		submodule models.Submodule
		expected  string
	}{
		{
			submodule: models.Submodule{
				Path: "/valid/path",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			expected: "Submodule example.com:443/repository@main>/valid/path",
		},
		{
			submodule: models.Submodule{
				Path: "",
				Url:  &urls.HttpUrl{"another-example.com", 8080, "another-repository", true},
			},
			expected: "Submodule another-example.com:8080/another-repository>another-repository",
		},
	}

	for _, test := range tests {
		actual := test.submodule.String()
		if actual != test.expected {
			t.Errorf("Expected %s, but got %s", test.expected, actual)
		}
	}
}

func TestSubmoduleValidate(t *testing.T) {
	tests := []struct {
		submodule models.Submodule
		valid     bool
	}{
		{
			submodule: models.Submodule{
				Path: "foo",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
			},
			valid: true,
		},
		{
			submodule: models.Submodule{
				Path: "foo",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			valid: true,
		},
		{
			submodule: models.Submodule{
				Path: "",
				Url:  &urls.HttpUrl{},
				Ref:  "main",
			},
			valid: false,
		},
		{
			submodule: models.Submodule{
				Path: "*foo",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			valid: false,
		},
		{
			submodule: models.Submodule{
				Path: "fo*o",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			valid: false,
		},
		{
			submodule: models.Submodule{
				Path: "!foo",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			valid: false,
		},
		{
			submodule: models.Submodule{
				Path: "fo!o",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			valid: false,
		},
		{
			submodule: models.Submodule{
				Path: "/valid/path",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "invalid ref",
			},
			valid: false,
		},
		{
			submodule: models.Submodule{
				Path: "valid/path",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "valid",
			},
			valid: true,
		},
		{
			submodule: models.Submodule{
				Path: "valid/path",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
			},
			valid: true,
		},
	}

	for index, test := range tests {
		err := test.submodule.Validate()
		if test.valid && err != nil {
			t.Errorf("Validation failed for case %d: %v", index+1, err)
		}
		if !test.valid && err == nil {
			t.Errorf("Validation passed for case %d", index+1)
		}
	}
}
