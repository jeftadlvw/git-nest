package internal

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jeftadlvw/git-nest/models"
	"strings"
)

/*
PopulateNestConfigFromToml populates a models.NestConfig from a configuration in TOML's markup language.
*/
func PopulateNestConfigFromToml(nestConfig *models.NestConfig, s string, strict bool) error {
	md, err := toml.Decode(s, &nestConfig)
	if err != nil {
		return err
	}

	undecoded := md.Undecoded()
	if len(undecoded) != 0 && strict {
		return fmt.Errorf("nest config contains undecoded keys: %q", undecoded)
	}

	return nil
}

/*
SubmoduleToTomlConfig returns a configuration string in TOML's markup language for a single models.Submodule.
*/
func SubmoduleToTomlConfig(s models.Submodule, indent string) string {
	var sb strings.Builder

	sb.WriteString("[[submodule]]")
	sb.WriteString("\n")

	sb.WriteString(formatTomlKeyValue("path", s.Path.UnixString(), indent))

	urlStr := ""
	if s.Url != nil {
		urlStr = s.Url.String()
	}

	sb.WriteString(formatTomlKeyValue("url", urlStr, indent))

	if s.Ref != "" {
		sb.WriteString(formatTomlKeyValue("ref", s.Ref, indent))
	}

	return strings.TrimSpace(sb.String())
}

/*
SubmodulesToTomlConfig returns a configuration string in TOML's markup language for more structs of type models.Submodule.
*/
func SubmodulesToTomlConfig(indent string, submodules ...models.Submodule) string {
	var sb strings.Builder
	for _, submodule := range submodules {
		sb.WriteString(SubmoduleToTomlConfig(submodule, indent))
		sb.WriteString("\n\n")
	}

	return strings.TrimSpace(sb.String())
}

/*
formatTomlKeyValue formats a key and value in TOML's markup language.
*/
func formatTomlKeyValue(k string, v string, indent string) string {
	return fmt.Sprintf("%s%s = \"%s\"\n", indent, k, v)
}
