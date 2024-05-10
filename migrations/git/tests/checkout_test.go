package tests

import (
	"github.com/jeftadlvw/git-nest/interfaces"
	"github.com/jeftadlvw/git-nest/migrations/git"
	"testing"
)

func TestCheckoutImplementsInterface(t *testing.T) {
	var _ interfaces.Migration = (*git.Checkout)(nil)
}

func TestCheckout(t *testing.T) {
	// there is no real need to test this migration, as it's only a wrapper that returns a formatted error.
	// the wrapped functionality is tested at utils/tests/git_test.go.
}
