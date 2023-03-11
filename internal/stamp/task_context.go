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
	return &TaskContext{
		DryRun: dryRun,
		IO:     app.IO,
		UI:     app.UI,
		Logger: NewTaskLogger(app.UI, dryRun),
		Store:  app.Store,
	}
}
