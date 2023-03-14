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
}

type SetDefaultsFunc func(any) error

var (
	// DefaultSetDefaultsFunc is the default SetDefaults implementation.
	// Delegates to [defaults.Set].
	DefaultSetDefaultsFunc SetDefaultsFunc = defaults.Set

	// SetDefaults sets the default values for the given task.
	SetDefaults = DefaultSetDefaultsFunc
)

// NewTask returns a new Task struct for the given map of data.
func NewTask(taskData map[string]any) (Task, error) { //nolint:ireturn // intentional
	taskType, ok := taskData["type"]
	if !ok {
		return nil, errors.New("undefined task type")
	}

	var task Task // these should all be pointers
	switch taskType {
	case "create", "generate":
		task = &CreateTask{}
	case "delete":
		task = &DeleteTask{}
	case "generator":
		task = &GeneratorTask{}
	default:
		return nil, fmt.Errorf("unknown task type: %v", taskType)
	}

	// Using a custom mapstructure setup so that we can leverage the
	// built in validation of our enum types (based on encoding.TextMarshaler).
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.TextUnmarshallerHookFunc(),
		Result:     task,
	}
	decoder, _ := mapstructure.NewDecoder(decoderConfig)

	// Set struct defaults, decode data map into the struct, and then validate
	if err := SetDefaults(task); err != nil {
		return nil, err
	}
	// if err := mapstructure.Decode(taskData, task); err != nil {
	if err := decoder.Decode(taskData); err != nil {
		return nil, err
	}
	if err := validate.Struct(task); err != nil {
		return nil, err
	}
	return task, nil
}
