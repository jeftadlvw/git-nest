package tests

import (
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/migrations/submodules"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/models/urls"
	"testing"
)

func TestUpdateUrlImplementsInterface(t *testing.T) {
	var _ interfaces.Migration = (*submodules.UpdateRef)(nil)
}

func TestUpdateUrl(t *testing.T) {
	exampleUrl, err := urls.HttpUrlFromString("https://example.com/")
	if err != nil {
		t.Fatalf("could not parse example url: %s", err)
	}

	tests := []struct {
		submodule *models.Submodule
		url       urls.HttpUrl
		err       bool
	}{
		{nil, urls.HttpUrl{}, true},
		{&models.Submodule{}, urls.HttpUrl{}, true},
		{&models.Submodule{}, exampleUrl, false},
	}

	for index, tc := range tests {
		err := submodules.UpdateUrl{
			Submodule: tc.submodule,
			Url:       tc.url,
		}.Migrate()

		if tc.err && err == nil {
			t.Fatalf("TestAppendSubmodule-%d expected error", index+1)
		}
		if !tc.err && err != nil {
			t.Fatalf("TestAppendSubmodule-%d unexpected error: %s", index+1, err)
		}
		if !tc.err && tc.submodule.Url.String() != exampleUrl.String() {
			t.Fatalf("TestAppendSubmodule-%d: context did not apply url", index+1)
		}
	}
}
