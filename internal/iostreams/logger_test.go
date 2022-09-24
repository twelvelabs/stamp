package iostreams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIconLogger(t *testing.T) {
	ios, _, _, _ := Test()
	formatter := ios.Formatter()

	logger := NewIconLogger(ios, formatter)

	assert.Equal(t, ios, logger.ios)
	assert.Equal(t, formatter, logger.formatter)
}

func TestIconLogger_Info(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		args     []any
		expected string
	}{
		{
			name:     "logs to stderr with an info icon",
			line:     "hello",
			args:     []any{},
			expected: "• hello",
		},
		{
			name:     "handles printf strings and args",
			line:     "hello %s",
			args:     []any{"world"},
			expected: "• hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, stderr := Test()

			logger := NewIconLogger(ios, ios.Formatter())
			logger.Info(tt.line, tt.args...)

			assert.Equal(t, tt.expected, stderr.String())
		})
	}
}

func TestIconLogger_Success(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		args     []any
		expected string
	}{
		{
			name:     "logs to stderr with a success icon",
			line:     "hello",
			args:     []any{},
			expected: "✓ hello",
		},
		{
			name:     "handles printf strings and args",
			line:     "hello %s",
			args:     []any{"world"},
			expected: "✓ hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, stderr := Test()

			logger := NewIconLogger(ios, ios.Formatter())
			logger.Success(tt.line, tt.args...)

			assert.Equal(t, tt.expected, stderr.String())
		})
	}
}

func TestIconLogger_Warning(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		args     []any
		expected string
	}{
		{
			name:     "logs to stderr with a warning icon",
			line:     "hello",
			args:     []any{},
			expected: "! hello",
		},
		{
			name:     "handles printf strings and args",
			line:     "hello %s",
			args:     []any{"world"},
			expected: "! hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, stderr := Test()

			logger := NewIconLogger(ios, ios.Formatter())
			logger.Warning(tt.line, tt.args...)

			assert.Equal(t, tt.expected, stderr.String())
		})
	}
}

func TestIconLogger_Failure(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		args     []any
		expected string
	}{
		{
			name:     "logs to stderr with a failure icon",
			line:     "hello",
			args:     []any{},
			expected: "✖ hello",
		},
		{
			name:     "handles printf strings and args",
			line:     "hello %s",
			args:     []any{"world"},
			expected: "✖ hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, stderr := Test()

			logger := NewIconLogger(ios, ios.Formatter())
			logger.Failure(tt.line, tt.args...)

			assert.Equal(t, tt.expected, stderr.String())
		})
	}
}
