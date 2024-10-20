package stamp

import (
	"github.com/mitchellh/copystructure"
)

// NewTaskSet returns a new TaskSet.
func NewTaskSet() *TaskSet {
	return &TaskSet{
		tasks: []Task{},
	}
}

// TaskSet is an ordered set of tasks that can be executed sequentially.
type TaskSet struct {
	tasks     []Task
	SrcPath   string
	DstPath   string
	Generator *Generator
}

// All returns all tasks in the set.
func (ts *TaskSet) All() []Task {
	return ts.tasks
}

// Add adds a task to the set.
func (ts *TaskSet) Add(t Task) *TaskSet {
	if pt, ok := t.(*PluginTask); ok {
		pt.Generator = ts.Generator
		t = pt
		// fmt.Printf("PT: %#v \n", pt)
	}
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
		if iter := t.Iterator(values); iter != nil { //nolint: nestif
			for i, item := range iter {
				values["_Index"] = i
				values["_Item"] = item
				if t.ShouldExecute(values) {
					err := t.Execute(ctx, values)
					if err != nil {
						return err
					}
				}
			}
		} else if t.ShouldExecute(values) {
			err := t.Execute(ctx, values)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
