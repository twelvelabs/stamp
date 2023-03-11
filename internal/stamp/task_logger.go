package stamp

import (
	"fmt"
	"strings"

	"github.com/twelvelabs/termite/ui"
)

const (
	ActionWidth = 10
)

// NewTaskLogger returns a new TaskLogger.
func NewTaskLogger(u *ui.UserInterface, dryRun bool) *TaskLogger {
	return &TaskLogger{
		ui:     u,
		dryRun: dryRun,
	}
}

// TaskLogger logs formatted Task actions.
type TaskLogger struct {
	ui     *ui.UserInterface
	dryRun bool
}

// Info logs a line to os.StdErr prefixed with an info icon and action label.
//
// Example:
//
//	// Prints "• [    action]: hello, world\n"
//	Info("action", "hello, %s", "world")
func (l *TaskLogger) Info(action string, line string, args ...any) {
	icon := l.ui.InfoIcon()
	action = l.ui.Info(l.rightJustify(action))
	l.log(icon, action, line, args...)
}

// Success logs a line to os.StdErr prefixed with a success icon and action label.
//
// Example:
//
//	// Prints "✓ [    action]: hello, world\n"
//	Success("action", "hello, %s", "world")
func (l *TaskLogger) Success(action string, line string, args ...any) {
	icon := l.ui.SuccessIcon()
	action = l.ui.Success(l.rightJustify(action))
	l.log(icon, action, line, args...)
}

// Warning logs a line to os.StdErr prefixed with a warning icon and action label.
//
// Example:
//
//	// Prints "! [    action]: hello, world\n"
//	Warning("action", "hello, %s", "world")
func (l *TaskLogger) Warning(action string, line string, args ...any) {
	icon := l.ui.WarningIcon()
	action = l.ui.Warning(l.rightJustify(action))
	l.log(icon, action, line, args...)
}

// Failure logs a line to os.StdErr prefixed with a failure icon and action label.
//
// Example:
//
//	// Prints "✖ [    action]: hello, world\n"
//	Failure("action", "hello, %s", "world")
func (l *TaskLogger) Failure(action string, line string, args ...any) {
	icon := l.ui.FailureIcon()
	action = l.ui.Failure(l.rightJustify(action))
	l.log(icon, action, line, args...)
}

// logs the formatted icon, action, and line to StdErr.
func (l *TaskLogger) log(icon, action, line string, args ...any) {
	prefix := icon + " "
	if l.dryRun {
		prefix += "[DRY RUN]"
	}
	line = l.ensureNewline(prefix + "[" + action + "]: " + line)
	l.ui.Out(line, args...)
}

func (l *TaskLogger) ensureNewline(text string) string {
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	return text
}

func (l *TaskLogger) rightJustify(text string) string {
	return fmt.Sprintf("%*s", ActionWidth, text)
}
