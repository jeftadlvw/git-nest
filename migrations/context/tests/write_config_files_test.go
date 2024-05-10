package tests

import (
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/migrations/context"
	"testing"
)

func TestWriteConfigFilesImplementsInterface(t *testing.T) {
	var _ interfaces.Migration = (*context.WriteConfigFiles)(nil)
}

func TestWriteConfigFiles(t *testing.T) {
	// there is no real need to test this migration, as it's only a wrapper that returns a formatted error.
	// the wrapped functionality is tested at internal/tests/write_config_test.go.
}
