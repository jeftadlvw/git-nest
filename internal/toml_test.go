package internal

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"strings"
	"testing"
)

func TestPopulateNestConfigFromToml(t *testing.T) {
	inputString := `
[config]

[[submodule]]
  path = "path-to-submodule"
  url = "url/to/repository"
  ref = "branch, tag or commit"
`
	getNestConfig := models.NestConfig{}

	err := PopulateNestConfigFromToml(&getNestConfig, inputString)
	if err != nil {
		t.Fatalf("Error getting NestConfigFromString: %v", err)
	}

	fmt.Printf("%+v\n", getNestConfig)
}

func TestSubmoduleArrTomlStrFromNestConfig(t *testing.T) {

	nestConfig := models.NestConfig{}
	nestConfig.Submodules = append(nestConfig.Submodules, models.Submodule{})
	nestConfig.Submodules = append(nestConfig.Submodules, models.Submodule{})

	expectedOutput := `[[submodule]]
  path = ""
  url = "http://:0"

[[submodule]]
  path = ""
  url = "http://:0"`

	getOutput := SubmodulesToTomlConfig("  ", nestConfig.Submodules...)
	if strings.TrimSpace(getOutput) != strings.TrimSpace(expectedOutput) {
		t.Fatalf("\nExpected:\n%s\nActual:\n%s", expectedOutput, getOutput)
	}
}
