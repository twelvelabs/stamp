package gen

import (
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/value"
)

type TaskContext struct {
	DryRun   bool
	IO       *iostreams.IOStreams
	Logger   iostreams.Logger
	Prompter value.Prompter
	Store    *Store
}

func NewTaskContext(ios *iostreams.IOStreams, prompter value.Prompter, store *Store) *TaskContext {
	// TODO: set logger
	return &TaskContext{
		IO:       ios,
		Prompter: prompter,
		Store:    store,
	}
}
