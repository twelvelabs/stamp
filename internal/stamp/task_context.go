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
func NewTaskContext(ios *ui.IOStreams, prompter value.Prompter, store *Store, dryRun bool) *TaskContext {
	var logger *TaskLogger
	// some tests need a context but not any of it's content, and so may pass in nil args.
	if ios != nil {
		logger = NewTaskLogger(ios, ios.Formatter(), dryRun)
	}
	return &TaskContext{
		DryRun:   dryRun,
		IO:       ios,
		Logger:   logger,
		Prompter: prompter,
		Store:    store,
	}
}
