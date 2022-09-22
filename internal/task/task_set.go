package task

//cspell: words copystructure

import (
	//cspell: disable
	"github.com/mitchellh/copystructure"
	//cspell: enable

	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/value"
)

// NewTaskSet returns a new TaskSet.
func NewTaskSet() *TaskSet {
	return &TaskSet{
		tasks: []Task{},
	}
}

// TaskSet is an ordered set of tasks that can be executed sequentially.
type TaskSet struct {
	tasks []Task
}

// All returns all tasks in the set.
func (ts *TaskSet) All() []Task {
	return ts.tasks
}

// Add adds a task to the set.
func (ts *TaskSet) Add(t Task) {
	ts.tasks = append(ts.tasks, t)
}

// Execute executes all tasks in order.
// Tasks that return false from `ShouldExecute()` are skipped.
// Tasks that return a slice from `Iterator()` will be executed once
// per item in the slice.
func (ts *TaskSet) Execute(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, dryRun bool) error {
	for _, t := range ts.All() {
		if !t.ShouldExecute(values) {
			continue
		}
		if iter := t.Iterator(values); iter != nil {
			for i, item := range iter {
				// deep-copy values
				copied, err := copystructure.Copy(values)
				if err != nil {
					return err
				}

				casted := copied.(map[string]any)
				casted["_Index"] = i
				casted["_Item"] = item

				err = t.Execute(casted, ios, prompter, dryRun)
				if err != nil {
					return err
				}
			}
		} else {
			err := t.Execute(values, ios, prompter, dryRun)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
