package tests

import (
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/migrations/fs"
	"github.com/jeftadlvw/git-nest/models"
	"testing"
)

func TestDeleteDirectoryImplementsInterface(t *testing.T) {
	var _ interfaces.Migration = (*fs.DeleteDirectory)(nil)
}

func TestDeleteDirectory(t *testing.T) {

	tests := []struct {
		path string
		err  bool
	}{
		{"", true},
		{".", true},
		{"/", true},
		{t.TempDir(), false},
	}

	for index, tc := range tests {
		err := fs.DeleteDirectory{
			Path:   models.Path(tc.path),
			DryRun: true,
		}.Migrate()

		if tc.err && err == nil {
			t.Fatalf("TestAppendSubmodule-%d expected error", index+1)
		}
		if !tc.err && err != nil {
			t.Fatalf("TestAppendSubmodule-%d unexpected error: %s", index+1, err)
		}
	}
}
