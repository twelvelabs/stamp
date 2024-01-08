package stamp

import (
	"errors"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/twelvelabs/termite/validate"
)

//go:generate moq -rm -out task_mock.go . Task

// Task is the interface a generator task.
type Task interface {
	// Iterator returns a slice of values if the task should be run more than once.
	// Configured via the `each` attribute (default nil).
	Iterator(values map[string]any) []any

	// Execute executes the task.
	Execute(context *TaskContext, values map[string]any) error

	// ShouldExecute returns true if the task should be run.
	// Configured via the `if` attribute (default `true`).
	ShouldExecute(values map[string]any) bool

	// TypeKey returns the short name used to identify the task in generator metadata.
	TypeKey() string
}

// AllTasks returns a slice of all known tasks.
// Whenever a new task type is created, it should be added here.
func AllTasks() []Task {
	tasks := []Task{
		&CreateTask{},
		&UpdateTask{},
		&DeleteTask{},
		&GeneratorTask{},
	}
	// Ensure defaults are set.
	// Needed for TypeKey() implementation.
	for _, t := range tasks {
		defaults.MustSet(t)
	}
	return tasks
}

// NewTask returns a new Task struct for the given map of data.
func NewTask(taskData map[string]any) (Task, error) {
	taskType, ok := taskData["type"]
	if !ok {
		return nil, errors.New("undefined task type")
	}

	var task Task
	for _, t := range AllTasks() {
		if t.TypeKey() == taskType {
			task = t
			break
		}
	}
	if task == nil {
		return nil, fmt.Errorf("unknown task type: %v", taskType)
	}

	// Using a custom mapstructure setup so that we can leverage the
	// built in validation of our enum types (based on encoding.TextMarshaler).
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.TextUnmarshallerHookFunc(),
		Result:     task,
	}
	decoder, _ := mapstructure.NewDecoder(decoderConfig)

	// if err := mapstructure.Decode(taskData, task); err != nil {
	if err := decoder.Decode(taskData); err != nil {
		return nil, err
	}
	if err := validate.Struct(task); err != nil {
		return nil, err
	}
	return task, nil
}
