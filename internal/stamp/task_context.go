package stamp

import (
	"github.com/twelvelabs/termite/ui"
)

// TaskContext holds configuration and dependencies used in Task.Execute().
type TaskContext struct {
	DryRun bool
	IO     *ui.IOStreams
	UI     *ui.UserInterface
	Logger *TaskLogger
	Store  *Store
}

// NewTaskContext returns a configured TaskContext.
func NewTaskContext(app *App, dryRun bool) *TaskContext {
	logger := NewTaskLogger(app.IO, app.IO.Formatter(), dryRun)
	return &TaskContext{
		DryRun: dryRun,
		IO:     app.IO,
		UI:     app.UI,
		Logger: logger,
		Store:  app.Store,
	}
}
