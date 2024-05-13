package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"testing"
)

func TestSubmoduleCleanUp(t *testing.T) {
	submodule := models.Submodule{
		Path: "/err/../path",
		Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
		Ref:  "  main  ",
	}

	submodule.Clean()

	if submodule.Path != "/path" {
		t.Fatalf("expected cleaned up path to be 'path', got %s", submodule.Path)
	}

	if submodule.Ref != "main" {
		t.Fatalf("expected cleaned up ref to be main, got %s", submodule.Ref)
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

	for index, test := range tests {
		t.Run(fmt.Sprintf("TestSubmoduleRemoteIdentifier-%d", index+1), func(t *testing.T) {
			actual := test.submodule.RemoteIdentifier()
			if actual != test.expected {
				t.Fatalf("expected %s, but got %s", test.expected, actual)
			}
		})
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

	for index, test := range tests {
		t.Run(fmt.Sprintf("TestSubmoduleRemoteIdentifier-%d", index+1), func(t *testing.T) {
			actual := test.submodule.Identifier()
			if actual != test.expected {
				t.Fatalf("expected %s, but got %s", test.expected, actual)
			}
		})
	}
}

func TestSubmoduleString(t *testing.T) {
	tests := []struct {
		submodule models.Submodule
		expected  string
	}{
		{
			submodule: models.Submodule{
				Path: "/err/path",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			expected: "Submodule example.com:443/repository@main>/err/path",
		},
		{
			submodule: models.Submodule{
				Path: "",
				Url:  &urls.HttpUrl{"another-example.com", 8080, "another-repository", true},
			},
			expected: "Submodule another-example.com:8080/another-repository>another-repository",
		},
	}

	for index, test := range tests {
		t.Run(fmt.Sprintf("TestSubmoduleString-%d", index+1), func(t *testing.T) {
			actual := test.submodule.String()
			if actual != test.expected {
				t.Fatalf("Expected %s, but got %s", test.expected, actual)
			}
		})
	}
}

func TestSubmoduleValidate(t *testing.T) {
	tests := []struct {
		submodule models.Submodule
		err       bool
	}{
		{
			submodule: models.Submodule{
				Path: "foo",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
			},
			err: true,
		},
		{
			submodule: models.Submodule{
				Path: "foo",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			err: true,
		},
		{
			submodule: models.Submodule{
				Path: "",
				Url:  &urls.HttpUrl{},
				Ref:  "main",
			},
			err: false,
		},
		{
			submodule: models.Submodule{
				Path: "*foo",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			err: false,
		},
		{
			submodule: models.Submodule{
				Path: "fo*o",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			err: false,
		},
		{
			submodule: models.Submodule{
				Path: "!foo",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			err: false,
		},
		{
			submodule: models.Submodule{
				Path: "fo!o",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "main",
			},
			err: false,
		},
		{
			submodule: models.Submodule{
				Path: "/err/path",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "invalid ref",
			},
			err: false,
		},
		{
			submodule: models.Submodule{
				Path: "err/path",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
				Ref:  "err",
			},
			err: true,
		},
		{
			submodule: models.Submodule{
				Path: "err/path",
				Url:  &urls.HttpUrl{"example.com", 443, "repository", true},
			},
			err: true,
		},
	}

	for index, test := range tests {
		t.Run(fmt.Sprintf("TestSubmoduleValidate-%d", index+1), func(t *testing.T) {
			err := test.submodule.Validate()
			if test.err && err != nil {
				t.Fatalf("failed validation: %v", err)
			}
			if !test.err && err == nil {
				t.Fatalf("validation passed")
			}
		})
	}
}
