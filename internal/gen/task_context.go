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
	return &TaskContext{
		DryRun:   dryRun,
		IO:       ios,
		Logger:   iostreams.NewActionLogger(ios, ios.Formatter(), dryRun),
		Prompter: prompter,
		Store:    store,
	}
}
