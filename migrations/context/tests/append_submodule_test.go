package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/migrations/context"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"slices"
	"testing"
)

func TestAppendSubmoduleImplementsInterface(t *testing.T) {
	var _ interfaces.Migration = (*context.AppendSubmodule)(nil)
}

func TestAppendSubmodule(t *testing.T) {
	exampleUrl, err := urls.HttpUrlFromString("https://example.com/")
	if err != nil {
		t.Fatalf("could not parse example url: %s", err)
	}

	tests := []struct {
		context   *models.NestContext
		submodule models.Submodule
		err       bool
	}{
		{nil, models.Submodule{}, true},
		{&models.NestContext{}, models.Submodule{}, true},
		{&models.NestContext{}, models.Submodule{Path: "foo"}, true},
		{&models.NestContext{}, models.Submodule{Path: "foo", Url: &exampleUrl}, false},
	}

	for index, tc := range tests {
		t.Run(fmt.Sprintf("TestAppendSubmodule-%d", index+1), func(t *testing.T) {
			err := context.AppendSubmodule{
				Context:   tc.context,
				Submodule: tc.submodule,
			}.Migrate()

			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !tc.err && !slices.Contains(tc.context.Config.Submodules, tc.submodule) {
				t.Fatalf("context did not contain submodule")
			}
		})
	}
}
