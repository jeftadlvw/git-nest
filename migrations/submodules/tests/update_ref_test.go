package tests

import (
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/migrations/submodules"
	"github.com/jeftadlvw/git-nest/models"
	"strings"
	"testing"
)

func TestUpdateRefImplementsInterface(t *testing.T) {
	var _ interfaces.Migration = (*submodules.UpdateRef)(nil)
}

func TestUpdateRef(t *testing.T) {
	tests := []struct {
		submodule *models.Submodule
		ref       string
		err       bool
	}{
		{nil, "", true},
		{&models.Submodule{}, "", true},
		{&models.Submodule{}, "foo", false},
		{&models.Submodule{}, "   foo  \n  ", false},
	}

	for index, tc := range tests {
		err := submodules.UpdateRef{
			Submodule: tc.submodule,
			Ref:       tc.ref,
		}.Migrate()

		if tc.err && err == nil {
			t.Fatalf("TestAppendSubmodule-%d expected error", index+1)
		}
		if !tc.err && err != nil {
			t.Fatalf("TestAppendSubmodule-%d unexpected error: %s", index+1, err)
		}
		if !tc.err && tc.submodule != nil {
			if tc.submodule.Ref != strings.TrimSpace(tc.ref) {
				t.Fatalf("TestAppendSubmodule-%d: context did not set ref", index+1)
			}
		}
	}
}
