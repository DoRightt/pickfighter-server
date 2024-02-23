package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestValidateServerArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected error
	}{
		{"ValidRoute", []string{"auth"}, nil},
		{"EmptyArgs", []string{}, errEmptyApiRoute},
		{"InvalidRoute", []string{"invalidRoute"}, fmt.Errorf("allowed routes are: %s", "auth, common")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := validateServerArgs(cmd, tc.args)

			assert.Equal(t, tc.expected, err)
		})
	}
}
