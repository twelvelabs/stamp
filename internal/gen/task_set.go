package gen

//cspell: words copystructure

import (
	//cspell: disable
	"github.com/mitchellh/copystructure"
	//cspell: enable
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
func (ts *TaskSet) Execute(ctx *TaskContext, values map[string]any) error {
	for _, t := range ts.All() {
		t.SetDryRun(ctx.DryRun)
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

				err = t.Execute(ctx, casted)
				if err != nil {
					return err
				}
			}
		} else {
			err := t.Execute(ctx, values)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
