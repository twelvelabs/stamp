package gen

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
	tasks   []Task
	SrcPath string
	DstPath string
}

// All returns all tasks in the set.
func (ts *TaskSet) All() []Task {
	return ts.tasks
}

// Add adds a task to the set.
func (ts *TaskSet) Add(t Task) *TaskSet {
	ts.tasks = append(ts.tasks, t)
	return ts
}

// Execute executes all tasks in order.
// Tasks that return false from `ShouldExecute()` are skipped.
// Tasks that return a slice from `Iterator()` will be executed once
// per item in the slice.
func (ts *TaskSet) Execute(ctx *TaskContext, data map[string]any) error {
	// deep-copy values
	copied, err := copystructure.Copy(data)
	if err != nil {
		return err
	}
	values := copied.(map[string]any)

	values["SrcPath"] = ts.SrcPath
	if _, ok := values["DstPath"]; !ok {
		values["DstPath"] = ts.DstPath
	}

	for _, t := range ts.All() {
		if !t.ShouldExecute(values) {
			continue
		}
		if iter := t.Iterator(values); iter != nil {
			for i, item := range iter {
				values["_Index"] = i
				values["_Item"] = item
				err := t.Execute(ctx, values)
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
