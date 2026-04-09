package clidemo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		errMsg      string
	}{
		{
			name:        "empty string",
			input:       "",
			shouldError: true,
			errMsg:      "input cannot be empty",
		},
		{
			name:        "only spaces",
			input:       "   ",
			shouldError: true,
			errMsg:      "input cannot be empty",
		},
		{
			name:        "only tabs",
			input:       "\t\t\t",
			shouldError: true,
			errMsg:      "input cannot be empty",
		},
		{
			name:        "only newlines",
			input:       "\n\n\n",
			shouldError: true,
			errMsg:      "input cannot be empty",
		},
		{
			name:        "mixed whitespace",
			input:       " \t\n \t\n ",
			shouldError: true,
			errMsg:      "input cannot be empty",
		},
		{
			name:        "unicode whitespace",
			input:       "\u0020\u00A0\u2000\u2001",
			shouldError: true,
			errMsg:      "input cannot be empty",
		},
		{
			name:        "valid single word",
			input:       "hello",
			shouldError: false,
		},
		{
			name:        "valid with leading spaces",
			input:       "  valid",
			shouldError: false,
		},
		{
			name:        "valid with trailing spaces",
			input:       "valid  ",
			shouldError: false,
		},
		{
			name:        "valid with surrounding spaces",
			input:       "  valid  ",
			shouldError: false,
		},
		{
			name:        "valid multi-word",
			input:       "Learn Go",
			shouldError: false,
		},
		{
			name:        "valid with tabs",
			input:       "\tvalid\t",
			shouldError: false,
		},
		{
			name:        "valid unicode",
			input:       "你好世界",
			shouldError: false,
		},
		{
			name:        "valid emoji",
			input:       "🎉",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInput(tt.input)

			if tt.shouldError {
				assert.Error(t, err, "validateInput(%q) should return error", tt.input)
				assert.EqualError(t, err, tt.errMsg, "error message mismatch")
			} else {
				assert.NoError(t, err, "validateInput(%q) should not return error", tt.input)
			}
		})
	}
}
