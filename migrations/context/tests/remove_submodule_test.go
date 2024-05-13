package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/migrations/context"
	"github.com/jeftadlvw/git-nest/models"
	"testing"
)

func TestRemoveSubmoduleImplementsInterface(t *testing.T) {
	var _ interfaces.Migration = (*context.RemoveSubmodule)(nil)
}

func TestRemoveSubmodule(t *testing.T) {

	tests := []struct {
		existingSubmoduleCount int
		submoduleIndex         int
		err                    bool
	}{
		{0, 0, true},
		{0, 1, true},
		{1, 0, false},
		{1, 1, true},
		{1, -1, true},
		{2, 1, false},
		{3, 1, false},
	}

	for index, tc := range tests {
		t.Run(fmt.Sprintf("TestRemoveSubmodule-%d", index+1), func(t *testing.T) {
			mockContext := models.NestContext{}

			for range tc.existingSubmoduleCount {
				mockContext.Config.Submodules = append(mockContext.Config.Submodules, models.Submodule{})
			}

			err := context.RemoveSubmodule{
				Context:        &mockContext,
				SubmoduleIndex: tc.submoduleIndex,
			}.Migrate()

			if tc.err && err == nil {
				t.Fatalf("no error, but expected one")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !tc.err && len(mockContext.Config.Submodules) != tc.existingSubmoduleCount-1 {
				t.Fatalf("submodule was not removed")
			}
		})
	}
}
