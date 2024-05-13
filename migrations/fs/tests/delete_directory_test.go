package tests

import (
	"fmt"
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
		t.Run(fmt.Sprintf("TestDeleteDirectory-%d", index+1), func(t *testing.T) {
			err := fs.DeleteDirectory{
				Path:   models.Path(tc.path),
				DryRun: true,
			}.Migrate()

			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}
