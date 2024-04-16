package tests

import (
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
	for _, tc := range tests {
		if got := tc.config.Validate(); got != tc.expectedError {
			t.Errorf("Config.Validate() returned invalid error %v, wanted %v", got, tc.expectedError)
		}
	}
}
