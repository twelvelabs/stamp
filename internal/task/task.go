package task

import (
	"errors"
	"fmt"

	//cspell: disable
	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/task/generate"
	"github.com/twelvelabs/stamp/internal/value"
	//cspell: enable
)

//go:generate moq -rm -out task_mock.go . Task

// Task is the interface a generator task.
type Task interface {
	Iterator(values map[string]any) []any
	Execute(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, dryRun bool) error
	IsDryRun() bool
	SetDryRun(value bool)
	ShouldExecute(values map[string]any) bool
}

type SetDefaultsFunc func(any) error

var (
	// DefaultSetDefaultsFunc is the default SetDefaults implementation.
	// Delegates to [defaults.Set].
	DefaultSetDefaultsFunc SetDefaultsFunc = defaults.Set

	// SetDefaults sets the default values for the given task.
	SetDefaults SetDefaultsFunc = DefaultSetDefaultsFunc
)

// NewTask returns a new Task struct for the given map of data.
func NewTask(taskData map[string]any) (Task, error) {
	taskType, ok := taskData["type"]
	if !ok {
		return nil, errors.New("undefined task type")
	}

	var task Task
	switch taskType {
	case "generate":
		task = &generate.Task{}
	default:
		return nil, fmt.Errorf("unknown task type: %v", taskType)
	}

	// Set struct defaults, decode data map into the struct, and then validate
	if err := SetDefaults(task); err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(taskData, task); err != nil {
		return nil, err
	}
	if err := value.ValidateStruct(task); err != nil {
		return nil, err
	}
	return task, nil
}
