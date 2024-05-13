package tests

import (
	"errors"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		config        models.Config
		expectedError error
	}{
		{models.Config{true, true}, nil},
		{models.Config{false, false}, nil},
	}
	for index, tc := range tests {
		t.Run(fmt.Sprintf("TestConfigValidate-%d", index+1), func(t *testing.T) {
			if got := tc.config.Validate(); !errors.Is(got, tc.expectedError) {
				t.Fatalf("invalid error %v, expected %v", got, tc.expectedError)
			}
		})
	}
}
