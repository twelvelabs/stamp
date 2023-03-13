package stamp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewTaskMock(exe bool, exeErr error, iter []any) *TaskMock {
	return &TaskMock{
		IteratorFunc: func(values map[string]any) []any {
			return iter
		},
		ShouldExecuteFunc: func(values map[string]any) bool {
			return exe
		},
		ExecuteFunc: func(ctx *TaskContext, values map[string]any) error {
			return exeErr
		},
	}
}

func TestTaskSet_Add(t *testing.T) {
	task1 := &TaskMock{}
	task2 := &TaskMock{}
	ts := NewTaskSet()

	ts.Add(task1)
	assert.Equal(t, []Task{task1}, ts.All())
	ts.Add(task2)
	assert.Equal(t, []Task{task1, task2}, ts.All())
}

func TestTaskSet_AddsPathsWhenCallingExecute(t *testing.T) {
	task1 := NewTaskMock(true, nil, nil)

	ts := NewTaskSet().Add(task1)
	ts.SrcPath = "/path/to/src"
	ts.DstPath = "/path/to/dst"

	app := NewTestApp()
	ctx := NewTaskContext(app)
	values := map[string]any{}
	err := ts.Execute(ctx, values)

	assert.NoError(t, err)
	assert.Equal(t, map[string]any{}, values, "values should have been cloned before mutating")

	calls := task1.ExecuteCalls()
	assert.Equal(t, "/path/to/src", calls[0].Values["SrcPath"])
	assert.Equal(t, "/path/to/dst", calls[0].Values["DstPath"])

	// Now lets do it again, but with user supplied paths...
	task1 = NewTaskMock(true, nil, nil)
	ts = NewTaskSet().Add(task1)
	ts.SrcPath = "/path/to/src"
	ts.DstPath = "/path/to/dst"

	values = map[string]any{
		"SrcPath": "/custom/src/path",
		"DstPath": "/custom/dst/path",
	}
	err = ts.Execute(ctx, values)
	assert.NoError(t, err)

	calls = task1.ExecuteCalls()
	assert.Equal(t, "/path/to/src", calls[0].Values["SrcPath"], "src path can not be customized")
	assert.Equal(t, "/custom/dst/path", calls[0].Values["DstPath"], "dst path should be customized")
}

func TestTaskSet_OnlyExecutesTasksThatWantToBe(t *testing.T) {
	task1 := NewTaskMock(true, nil, nil)
	task2 := NewTaskMock(false, nil, nil)

	ts := NewTaskSet().Add(task1).Add(task2)

	app := NewTestApp()
	ctx := NewTaskContext(app)
	values := map[string]any{}
	err := ts.Execute(ctx, values)

	assert.NoError(t, err)
	assert.Len(t, task1.ExecuteCalls(), 1, "task1.Execute() should have been called.")
	assert.Len(t, task2.ExecuteCalls(), 0, "task2.Execute() should NOT have been called.")
}

func TestTaskSet_CanExecuteTasksMultipleTimes(t *testing.T) {
	indexes := []any{}
	items := []any{}

	task1 := NewTaskMock(true, nil, []any{"foo", "bar", "baz"})
	task1.ExecuteFunc = func(ctx *TaskContext, values map[string]any) error {
		indexes = append(indexes, values["_Index"])
		items = append(items, values["_Item"])
		return nil
	}

	ts := NewTaskSet().Add(task1)

	app := NewTestApp()
	ctx := NewTaskContext(app)
	values := map[string]any{}
	err := ts.Execute(ctx, values)

	assert.NoError(t, err)
	assert.Len(t, task1.ExecuteCalls(), 3, "task1.Execute() should have been called 3 times.")

	// Note: Can't look into task1.ExecuteCalls(), because each call contains
	// a reference/pointer to `values` - and thus by this point `_Index` and `_Items`
	// are set to the values from the last iteration.
	assert.Equal(t, 0, indexes[0])
	assert.Equal(t, "foo", items[0])

	assert.Equal(t, 1, indexes[1])
	assert.Equal(t, "bar", items[1])

	assert.Equal(t, 2, indexes[2])
	assert.Equal(t, "baz", items[2])
}

func TestTaskSet_HaltsExecutionAtTheFirstError(t *testing.T) {
	values := map[string]any{}

	task1 := NewTaskMock(true, nil, nil)
	task2 := NewTaskMock(true, errors.New("boom"), nil)
	task3 := NewTaskMock(true, nil, nil)

	ts := NewTaskSet()
	ts.Add(task1)
	ts.Add(task2)
	ts.Add(task3)

	app := NewTestApp()
	ctx := NewTaskContext(app)
	err := ts.Execute(ctx, values)

	assert.ErrorContains(t, err, "boom")
	assert.Len(t, task1.ExecuteCalls(), 1, "task1.Execute() should have been called.")
	assert.Len(t, task2.ExecuteCalls(), 1, "task2.Execute() should have been called.")
	assert.Len(t, task3.ExecuteCalls(), 0, "task3.Execute() should NOT have been called.")
}
