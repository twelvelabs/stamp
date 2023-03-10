package stamp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/ui"
)

func TestTaskLogger_Info(t *testing.T) {
	tests := []struct {
		name     string
		dryRun   bool
		action   string
		line     string
		args     []any
		expected string
	}{
		{
			name:     "logs to stderr with an info icon",
			action:   "test",
			line:     "hello",
			args:     []any{},
			expected: "• [      test]: hello\n",
		},
		{
			name:     "handles printf strings and args",
			action:   "test",
			line:     "hello %s",
			args:     []any{"world"},
			expected: "• [      test]: hello world\n",
		},
		{
			name:     "only adds trailing newline if needed",
			action:   "test",
			line:     "hello\n",
			args:     []any{},
			expected: "• [      test]: hello\n",
		},
		{
			name:     "adds dry run prefix",
			dryRun:   true,
			action:   "test",
			line:     "hello\n",
			args:     []any{},
			expected: "• [DRY RUN][      test]: hello\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios := ui.NewTestIOStreams()

			logger := NewTaskLogger(ios, ios.Formatter(), tt.dryRun)
			logger.Info(tt.action, tt.line, tt.args...)

			assert.Equal(t, tt.expected, ios.Err.String())
		})
	}
}

func TestTaskLogger_Success(t *testing.T) {
	tests := []struct {
		name     string
		dryRun   bool
		action   string
		line     string
		args     []any
		expected string
	}{
		{
			name:     "logs to stderr with an info icon",
			action:   "test",
			line:     "hello",
			args:     []any{},
			expected: "✓ [      test]: hello\n",
		},
		{
			name:     "handles printf strings and args",
			action:   "test",
			line:     "hello %s",
			args:     []any{"world"},
			expected: "✓ [      test]: hello world\n",
		},
		{
			name:     "only adds trailing newline if needed",
			action:   "test",
			line:     "hello\n",
			args:     []any{},
			expected: "✓ [      test]: hello\n",
		},
		{
			name:     "adds dry run prefix",
			dryRun:   true,
			action:   "test",
			line:     "hello\n",
			args:     []any{},
			expected: "✓ [DRY RUN][      test]: hello\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios := ui.NewTestIOStreams()

			logger := NewTaskLogger(ios, ios.Formatter(), tt.dryRun)
			logger.Success(tt.action, tt.line, tt.args...)

			assert.Equal(t, tt.expected, ios.Err.String())
		})
	}
}

func TestTaskLogger_Warning(t *testing.T) {
	tests := []struct {
		name     string
		dryRun   bool
		action   string
		line     string
		args     []any
		expected string
	}{
		{
			name:     "logs to stderr with an info icon",
			action:   "test",
			line:     "hello",
			args:     []any{},
			expected: "! [      test]: hello\n",
		},
		{
			name:     "handles printf strings and args",
			action:   "test",
			line:     "hello %s",
			args:     []any{"world"},
			expected: "! [      test]: hello world\n",
		},
		{
			name:     "only adds trailing newline if needed",
			action:   "test",
			line:     "hello\n",
			args:     []any{},
			expected: "! [      test]: hello\n",
		},
		{
			name:     "adds dry run prefix",
			dryRun:   true,
			action:   "test",
			line:     "hello\n",
			args:     []any{},
			expected: "! [DRY RUN][      test]: hello\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios := ui.NewTestIOStreams()

			logger := NewTaskLogger(ios, ios.Formatter(), tt.dryRun)
			logger.Warning(tt.action, tt.line, tt.args...)

			assert.Equal(t, tt.expected, ios.Err.String())
		})
	}
}

func TestTaskLogger_Failure(t *testing.T) {
	tests := []struct {
		name     string
		dryRun   bool
		action   string
		line     string
		args     []any
		expected string
	}{
		{
			name:     "logs to stderr with an info icon",
			action:   "test",
			line:     "hello",
			args:     []any{},
			expected: "✖ [      test]: hello\n",
		},
		{
			name:     "handles printf strings and args",
			action:   "test",
			line:     "hello %s",
			args:     []any{"world"},
			expected: "✖ [      test]: hello world\n",
		},
		{
			name:     "only adds trailing newline if needed",
			action:   "test",
			line:     "hello\n",
			args:     []any{},
			expected: "✖ [      test]: hello\n",
		},
		{
			name:     "adds dry run prefix",
			dryRun:   true,
			action:   "test",
			line:     "hello\n",
			args:     []any{},
			expected: "✖ [DRY RUN][      test]: hello\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios := ui.NewTestIOStreams()

			logger := NewTaskLogger(ios, ios.Formatter(), tt.dryRun)
			logger.Failure(tt.action, tt.line, tt.args...)

			assert.Equal(t, tt.expected, ios.Err.String())
		})
	}
}
