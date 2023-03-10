package stamp

import (
	"github.com/twelvelabs/termite/ui"

	"github.com/twelvelabs/stamp/internal/value"
)

// TaskContext holds configuration and dependencies used in Task.Execute().
type TaskContext struct {
	DryRun   bool
	IO       *ui.IOStreams
	Logger   *TaskLogger
	Prompter value.Prompter
	Store    *Store
}

// NewTaskContext returns a configured TaskContext.
func NewTaskContext(app *App, dryRun bool) *TaskContext {
	logger := NewTaskLogger(app.IO, app.IO.Formatter(), dryRun)
	return &TaskContext{
		DryRun:   dryRun,
		IO:       app.IO,
		Logger:   logger,
		Prompter: app.Prompter,
		Store:    app.Store,
	}
}
