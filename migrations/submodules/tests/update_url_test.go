package tests

import (
	"fmt"
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
		t.Run(fmt.Sprintf("TestUpdateUrl-%d", index+1), func(t *testing.T) {
			t.Parallel()
			err := submodules.UpdateUrl{
				Submodule: tc.submodule,
				Url:       tc.url,
			}.Migrate()

			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !tc.err && tc.submodule.Url.String() != exampleUrl.String() {
				t.Fatalf("context did not apply url")
			}
		})
	}
}
