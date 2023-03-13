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
func NewTaskContext(app *App) *TaskContext {
	return &TaskContext{
		DryRun: app.Config.DryRun,
		IO:     app.IO,
		UI:     app.UI,
		Logger: NewTaskLogger(app.UI, app.Config.DryRun),
		Store:  app.Store,
	}
}
