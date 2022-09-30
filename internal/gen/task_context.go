package gen

import (
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/value"
)

// TaskContext holds configuration and dependencies used in Task.Execute().
type TaskContext struct {
	DryRun   bool
	IO       *iostreams.IOStreams
	Logger   *iostreams.ActionLogger
	Prompter value.Prompter
	Store    *Store
}

// NewTaskContext returns a configured TaskContext.
func NewTaskContext(ios *iostreams.IOStreams, prompter value.Prompter, store *Store, dryRun bool) *TaskContext {
	var logger *iostreams.ActionLogger
	// some tests need a context but not any of it's content, and so may pass in nil args.
	if ios != nil {
		logger = iostreams.NewActionLogger(ios, ios.Formatter(), dryRun)
	}
	return &TaskContext{
		DryRun:   dryRun,
		IO:       ios,
		Logger:   logger,
		Prompter: prompter,
		Store:    store,
	}
}
