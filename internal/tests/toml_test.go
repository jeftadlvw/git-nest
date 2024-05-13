package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/internal"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"strings"
	"testing"
)

func TestPopulateNestConfigFromToml(t *testing.T) {
	inputString := `
[config]
allow_duplicate_origins = true

[[submodule]]
  path = "path-to-submodule"
  url = "https://example.com/url/to/repository"
  ref = "branch, tag or commit"

[[submodule]]
  path = "path-to-submodule"
  url = "https://example.com/url/to/repository"
  ref = "branch, tag or commit"
`
	nestConfig := models.NestConfig{}

	err := internal.PopulateNestConfigFromToml(&nestConfig, inputString, true)
	if err != nil {
		t.Fatalf("Error populating nest config from toml string: %v", err)
	}

	if nestConfig.Config.AllowDuplicateOrigins != true {
		t.Fatalf("Config.AllowDuplicateOrigins should be true")
	}

	if len(nestConfig.Submodules) != 2 {
		t.Fatalf("Submodules count mismatch (%d != %d)", len(nestConfig.Submodules), 2)
	}

	if nestConfig.Submodules[0].Path != "path-to-submodule" {
		t.Fatalf("submodule path does not match (%s != %s)", nestConfig.Submodules[0].Path, "path-to-submodule")
	}

	if nestConfig.Submodules[0].Url.String() != "https://example.com/url/to/repository" {
		t.Fatalf("submodule url does not match (%s != %s)", nestConfig.Submodules[0].Url.String(), "https://example.com/url/to/repository")
	}

	if nestConfig.Submodules[0].Ref != "branch, tag or commit" {
		t.Fatalf("submodule path does not match (%s != %s)", nestConfig.Submodules[0].Ref, "branch, tag or commit")
	}
}

func TestSubmoduleArrTomlStrFromNestConfig(t *testing.T) {

	nestConfig := models.NestConfig{}
	nestConfig.Submodules = append(nestConfig.Submodules, models.Submodule{})
	nestConfig.Submodules = append(nestConfig.Submodules, models.Submodule{"", &urls.HttpUrl{"example.com", 80, "/foo", false}, ""})

	expectedOutput := `[[submodule]]
  path = ""
  url = ""

[[submodule]]
  path = ""
  url = "http://example.com/foo"`

	getOutput := internal.SubmodulesToTomlConfig("  ", nestConfig.Submodules...)
	if strings.TrimSpace(getOutput) != strings.TrimSpace(expectedOutput) {
		t.Fatalf("\nExpected:\n>%s<\n\nActual:\n>%s<", expectedOutput, getOutput)
	}
}

func TestSubmoduleToTomlConfig(t *testing.T) {

	indent := "  "

	// empty submodule
	submodule := models.Submodule{}
	expectedOutput := `[[submodule]]
  path = ""
  url = ""`

	if output := internal.SubmoduleToTomlConfig(submodule, indent); output != expectedOutput {
		t.Fatalf("\nExpected:\n>%s<\n\nActual:\n>%s<", expectedOutput, output)
	}

	// set values
	submodule = models.Submodule{"example/path", &urls.HttpUrl{"example.com", 443, "", true}, "example-ref"}
	expectedOutput = `[[submodule]]
  path = "example/path"
  url = "https://example.com/"
  ref = "example-ref"`

	if output := internal.SubmoduleToTomlConfig(submodule, indent); output != expectedOutput {
		t.Fatalf("\nExpected:\n>%s<\n\nActual:\n>%s<", expectedOutput, output)
	}

	// double seperators
	submodule = models.Submodule{"example//path", &urls.HttpUrl{"example.com", 443, "", true}, "example-ref"}
	expectedOutput = `[[submodule]]
  path = "example/path"
  url = "https://example.com/"
  ref = "example-ref"`

	if output := internal.SubmoduleToTomlConfig(submodule, indent); output != expectedOutput {
		t.Fatalf("\nExpected:\n>%s<\n\nActual:\n>%s<", expectedOutput, output)
	}

	fmt.Println("2------")

	// windows path style
	submodule = models.Submodule{"example\\path", &urls.HttpUrl{"example.com", 443, "", true}, "example-ref"}
	expectedOutput = `[[submodule]]
  path = "example/path"
  url = "https://example.com/"
  ref = "example-ref"`

	if output := internal.SubmoduleToTomlConfig(submodule, indent); output != expectedOutput {
		t.Fatalf("\nExpected:\n>%s<\n\nActual:\n>%s<", expectedOutput, output)
	}

	// something messed up
	submodule = models.Submodule{"example\\\\path", &urls.HttpUrl{"example.com", 443, "", true}, "example-ref"}
	expectedOutput = `[[submodule]]
  path = "example/path"
  url = "https://example.com/"
  ref = "example-ref"`

	if output := internal.SubmoduleToTomlConfig(submodule, indent); output != expectedOutput {
		t.Fatalf("\nExpected:\n>%s<\n\nActual:\n>%s<", expectedOutput, output)
	}
}
