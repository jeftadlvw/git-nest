package tests

import (
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
		{&models.NestContext{}, models.Submodule{Path: "foo", Url: exampleUrl}, false},
	}

	for index, tc := range tests {
		err := context.AppendSubmodule{
			Context:   tc.context,
			Submodule: tc.submodule,
		}.Migrate()

		if tc.err && err == nil {
			t.Fatalf("TestAppendSubmodule-%d expected error", index+1)
		}
		if !tc.err && err != nil {
			t.Fatalf("TestAppendSubmodule-%d unexpected error: %s", index+1, err)
		}
		if !tc.err && !slices.Contains(tc.context.Config.Submodules, tc.submodule) {
			t.Fatalf("TestAppendSubmodule-%d: context did not contain submodule", index+1)
		}
	}
}
