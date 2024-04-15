package internal

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jeftadlvw/git-nest/models"
	"strings"
)

func PopulateNestConfigFromToml(nestConfig *models.NestConfig, s string) error {
	_, err := toml.Decode(s, &nestConfig)
	if err != nil {
		return err
	}

	return nil
}

func SubmodulesToTomlConfig(indent string, submodules ...models.Submodule) string {
	var sb strings.Builder
	for _, submodule := range submodules {
		sb.WriteString(SubmoduleToTomlConfig(submodule, indent))
		sb.WriteString("\n")
	}

	return strings.TrimSpace(sb.String())
}

func SubmoduleToTomlConfig(s models.Submodule, indent string) string {
	var sb strings.Builder

	sb.WriteString("[[submodule]]")
	sb.WriteString("\n")

	sb.WriteString(formatTomlKeyValue("path", s.Path.String(), indent))
	sb.WriteString(formatTomlKeyValue("url", s.Url.String(), indent))

	if s.Ref != "" {
		sb.WriteString(formatTomlKeyValue("ref", s.Ref, indent))
	}

	return sb.String()
}

func formatTomlKeyValue(k string, v string, indent string) string {
	return fmt.Sprintf("%s%s = \"%s\"\n", indent, k, v)
}
