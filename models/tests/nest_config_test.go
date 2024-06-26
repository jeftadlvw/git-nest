package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"testing"
)

func TestNestConfigValidate(t *testing.T) {
	cases := []struct {
		name   string
		config models.NestConfig
		err    bool
	}{
		// no submodules; err
		{
			name: "no submodules",
			config: models.NestConfig{
				Config:     models.Config{},
				Submodules: []models.Submodule{},
			},
			err: false,
		},

		// one submodule with path; err
		{
			name: "one submodule with path",
			config: models.NestConfig{
				Config: models.Config{},
				Submodules: []models.Submodule{
					{
						Path: "foo",
						Url:  &urls.HttpUrl{"example.com", 80, "path", false},
					},
				},
			},
			err: false,
		},

		// two submodules, including a ref; err
		{
			name: "two submodules, including a ref",
			config: models.NestConfig{
				Config: models.Config{},
				Submodules: []models.Submodule{
					{
						Path: "foo",
						Url:  &urls.HttpUrl{"example.com", 80, "path", false},
					},
					{
						Path: "bar",
						Url:  &urls.HttpUrl{"example.org", 80, "path", false},
						Ref:  "main",
					},
				},
			},
			err: false,
		},

		// duplicate use of directories; invalid
		{
			name: "duplicate use of directories",
			config: models.NestConfig{
				Config: models.Config{},
				Submodules: []models.Submodule{
					{
						Path: "foo",
						Url:  &urls.HttpUrl{"example.com", 80, "path", false},
					},
					{
						Path: "foo",
						Url:  &urls.HttpUrl{"example.org", 80, "path", false},
						Ref:  "main",
					},
				},
			},
			err: true,
		},

		// no module url; invalid
		{
			name: "no module url",
			config: models.NestConfig{
				Config: models.Config{},
				Submodules: []models.Submodule{
					{},
				},
			},
			err: true,
		},

		// duplicate submodule url paths in url without specified dirs; invalid
		{
			name: "duplicate submodule url paths in url without specified directories",
			config: models.NestConfig{
				Config: models.Config{},
				Submodules: []models.Submodule{
					{
						Url: &urls.HttpUrl{"example.com", 443, "path", true},
					},
					{
						Url: &urls.HttpUrl{"example.com", 443, "path", true},
					},
				},
			},
			err: true,
		},

		// duplicate submodule url paths with specified dirs; invalid
		{
			name: "duplicate submodule url paths with specified dirs",
			config: models.NestConfig{
				Config: models.Config{},
				Submodules: []models.Submodule{
					{
						Path: "foo",
						Url:  &urls.HttpUrl{"example.com", 443, "path", true},
					},
					{
						Path: "bar",
						Url:  &urls.HttpUrl{"example.com", 443, "path", true},
					},
				},
			},
			err: true,
		},

		// duplicate submodule url paths with specified dirs allowed by config; err
		{
			name: "duplicate submodule url paths with specified dirs allowed by config",
			config: models.NestConfig{
				Config: models.Config{AllowDuplicateOrigins: true},
				Submodules: []models.Submodule{
					{
						Path: "foo",
						Url:  &urls.HttpUrl{"example.com", 443, "path", true},
					},
					{
						Path: "bar",
						Url:  &urls.HttpUrl{"example.com", 443, "path", true},
					},
				},
			},
			err: false,
		},

		// duplicate submodule url paths + refs with specified dirs; invalid
		{
			name: "duplicate submodule url paths + refs with specified dirs",
			config: models.NestConfig{
				Config: models.Config{AllowDuplicateOrigins: true},
				Submodules: []models.Submodule{
					{
						Path: "foo",
						Url:  &urls.HttpUrl{"example.com", 443, "path", true},
						Ref:  "main",
					},
					{
						Path: "bar",
						Url:  &urls.HttpUrl{"example.com", 443, "path", true},
						Ref:  "main",
					},
				},
			},
			err: false,
		},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestConfigValidate-%d", index+1), func(t *testing.T) {
			err := tc.config.Validate()
			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if tc.err && err != nil {
				fmt.Printf("%s -> %s\n", tc.name, err)
			}
		})
	}
}
