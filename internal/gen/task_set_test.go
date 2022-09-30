package gen

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/value"
)

func NewTaskMock(exe bool, exeErr error, iter []any) *TaskMock {
	return &TaskMock{
		IteratorFunc: func(values map[string]any) []any {
			return iter
		},
		SetDryRunFunc: func(value bool) {},
		ShouldExecuteFunc: func(values map[string]any) bool {
			return exe
		},
		ExecuteFunc: func(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, dryRun bool) error {
			return exeErr
		},
	}
}

func TestTaskSetAdd(t *testing.T) {
	task1 := &TaskMock{}
	task2 := &TaskMock{}
	ts := NewTaskSet()

	ts.Add(task1)
	assert.Equal(t, []Task{task1}, ts.All())
	ts.Add(task2)
	assert.Equal(t, []Task{task1, task2}, ts.All())
}

func TestTaskSetOnlyExecutesTasksThatWantToBe(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	values := map[string]any{}

	task1 := NewTaskMock(true, nil, nil)
	task2 := NewTaskMock(false, nil, nil)

	ts := NewTaskSet()
	ts.Add(task1)
	ts.Add(task2)
	err := ts.Execute(values, ios, nil, false)

	assert.NoError(t, err)
	assert.Len(t, task1.ExecuteCalls(), 1, "task1.Execute() should have been called.")
	assert.Equal(t, values, task1.ExecuteCalls()[0].Values)
	assert.Len(t, task2.ExecuteCalls(), 0, "task2.Execute() should NOT have been called.")
}

func TestTaskSetCanExecuteTasksMultipleTimes(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	values := map[string]any{}

	iter := []any{"foo", "bar", "baz"}
	task1 := NewTaskMock(true, nil, iter)

	ts := NewTaskSet()
	ts.Add(task1)
	err := ts.Execute(values, ios, nil, false)

	assert.NoError(t, err)
	assert.Len(t, task1.ExecuteCalls(), 3, "task1.Execute() should have been called 3 times.")

	calls := task1.ExecuteCalls()
	assert.Equal(t, 0, calls[0].Values["_Index"])
	assert.Equal(t, "foo", calls[0].Values["_Item"])

	assert.Equal(t, 1, calls[1].Values["_Index"])
	assert.Equal(t, "bar", calls[1].Values["_Item"])

	assert.Equal(t, 2, calls[2].Values["_Index"])
	assert.Equal(t, "baz", calls[2].Values["_Item"])
}

func TestTaskSetHaltsExecutionAtTheFirstError(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	values := map[string]any{}

	task1 := NewTaskMock(true, nil, nil)
	task2 := NewTaskMock(true, errors.New("boom"), nil)
	task3 := NewTaskMock(true, nil, nil)

	ts := NewTaskSet()
	ts.Add(task1)
	ts.Add(task2)
	ts.Add(task3)
	err := ts.Execute(values, ios, nil, false)

	assert.ErrorContains(t, err, "boom")
	assert.Len(t, task1.ExecuteCalls(), 1, "task1.Execute() should have been called.")
	assert.Len(t, task2.ExecuteCalls(), 1, "task2.Execute() should have been called.")
	assert.Len(t, task3.ExecuteCalls(), 0, "task3.Execute() should NOT have been called.")
}
