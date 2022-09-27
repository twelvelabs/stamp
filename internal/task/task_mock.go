// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package task

import (
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/value"
	"sync"
)

// Ensure, that TaskMock does implement Task.
// If this is not the case, regenerate this file with moq.
var _ Task = &TaskMock{}

// TaskMock is a mock implementation of Task.
//
//	func TestSomethingThatUsesTask(t *testing.T) {
//
//		// make and configure a mocked Task
//		mockedTask := &TaskMock{
//			ExecuteFunc: func(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, dryRun bool) error {
//				panic("mock out the Execute method")
//			},
//			IteratorFunc: func(values map[string]any) []any {
//				panic("mock out the Iterator method")
//			},
//			SetDryRunFunc: func(valueMoqParam bool)  {
//				panic("mock out the SetDryRun method")
//			},
//			ShouldExecuteFunc: func(values map[string]any) bool {
//				panic("mock out the ShouldExecute method")
//			},
//		}
//
//		// use mockedTask in code that requires Task
//		// and then make assertions.
//
//	}
type TaskMock struct {
	// ExecuteFunc mocks the Execute method.
	ExecuteFunc func(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, dryRun bool) error

	// IteratorFunc mocks the Iterator method.
	IteratorFunc func(values map[string]any) []any

	// SetDryRunFunc mocks the SetDryRun method.
	SetDryRunFunc func(valueMoqParam bool)

	// ShouldExecuteFunc mocks the ShouldExecute method.
	ShouldExecuteFunc func(values map[string]any) bool

	// calls tracks calls to the methods.
	calls struct {
		// Execute holds details about calls to the Execute method.
		Execute []struct {
			// Values is the values argument value.
			Values map[string]any
			// Ios is the ios argument value.
			Ios *iostreams.IOStreams
			// Prompter is the prompter argument value.
			Prompter value.Prompter
			// DryRun is the dryRun argument value.
			DryRun bool
		}
		// Iterator holds details about calls to the Iterator method.
		Iterator []struct {
			// Values is the values argument value.
			Values map[string]any
		}
		// SetDryRun holds details about calls to the SetDryRun method.
		SetDryRun []struct {
			// ValueMoqParam is the valueMoqParam argument value.
			ValueMoqParam bool
		}
		// ShouldExecute holds details about calls to the ShouldExecute method.
		ShouldExecute []struct {
			// Values is the values argument value.
			Values map[string]any
		}
	}
	lockExecute       sync.RWMutex
	lockIterator      sync.RWMutex
	lockSetDryRun     sync.RWMutex
	lockShouldExecute sync.RWMutex
}

// Execute calls ExecuteFunc.
func (mock *TaskMock) Execute(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, dryRun bool) error {
	if mock.ExecuteFunc == nil {
		panic("TaskMock.ExecuteFunc: method is nil but Task.Execute was just called")
	}
	callInfo := struct {
		Values   map[string]any
		Ios      *iostreams.IOStreams
		Prompter value.Prompter
		DryRun   bool
	}{
		Values:   values,
		Ios:      ios,
		Prompter: prompter,
		DryRun:   dryRun,
	}
	mock.lockExecute.Lock()
	mock.calls.Execute = append(mock.calls.Execute, callInfo)
	mock.lockExecute.Unlock()
	return mock.ExecuteFunc(values, ios, prompter, dryRun)
}

// ExecuteCalls gets all the calls that were made to Execute.
// Check the length with:
//
//	len(mockedTask.ExecuteCalls())
func (mock *TaskMock) ExecuteCalls() []struct {
	Values   map[string]any
	Ios      *iostreams.IOStreams
	Prompter value.Prompter
	DryRun   bool
} {
	var calls []struct {
		Values   map[string]any
		Ios      *iostreams.IOStreams
		Prompter value.Prompter
		DryRun   bool
	}
	mock.lockExecute.RLock()
	calls = mock.calls.Execute
	mock.lockExecute.RUnlock()
	return calls
}

// Iterator calls IteratorFunc.
func (mock *TaskMock) Iterator(values map[string]any) []any {
	if mock.IteratorFunc == nil {
		panic("TaskMock.IteratorFunc: method is nil but Task.Iterator was just called")
	}
	callInfo := struct {
		Values map[string]any
	}{
		Values: values,
	}
	mock.lockIterator.Lock()
	mock.calls.Iterator = append(mock.calls.Iterator, callInfo)
	mock.lockIterator.Unlock()
	return mock.IteratorFunc(values)
}

// IteratorCalls gets all the calls that were made to Iterator.
// Check the length with:
//
//	len(mockedTask.IteratorCalls())
func (mock *TaskMock) IteratorCalls() []struct {
	Values map[string]any
} {
	var calls []struct {
		Values map[string]any
	}
	mock.lockIterator.RLock()
	calls = mock.calls.Iterator
	mock.lockIterator.RUnlock()
	return calls
}

// SetDryRun calls SetDryRunFunc.
func (mock *TaskMock) SetDryRun(valueMoqParam bool) {
	if mock.SetDryRunFunc == nil {
		panic("TaskMock.SetDryRunFunc: method is nil but Task.SetDryRun was just called")
	}
	callInfo := struct {
		ValueMoqParam bool
	}{
		ValueMoqParam: valueMoqParam,
	}
	mock.lockSetDryRun.Lock()
	mock.calls.SetDryRun = append(mock.calls.SetDryRun, callInfo)
	mock.lockSetDryRun.Unlock()
	mock.SetDryRunFunc(valueMoqParam)
}

// SetDryRunCalls gets all the calls that were made to SetDryRun.
// Check the length with:
//
//	len(mockedTask.SetDryRunCalls())
func (mock *TaskMock) SetDryRunCalls() []struct {
	ValueMoqParam bool
} {
	var calls []struct {
		ValueMoqParam bool
	}
	mock.lockSetDryRun.RLock()
	calls = mock.calls.SetDryRun
	mock.lockSetDryRun.RUnlock()
	return calls
}

// ShouldExecute calls ShouldExecuteFunc.
func (mock *TaskMock) ShouldExecute(values map[string]any) bool {
	if mock.ShouldExecuteFunc == nil {
		panic("TaskMock.ShouldExecuteFunc: method is nil but Task.ShouldExecute was just called")
	}
	callInfo := struct {
		Values map[string]any
	}{
		Values: values,
	}
	mock.lockShouldExecute.Lock()
	mock.calls.ShouldExecute = append(mock.calls.ShouldExecute, callInfo)
	mock.lockShouldExecute.Unlock()
	return mock.ShouldExecuteFunc(values)
}

// ShouldExecuteCalls gets all the calls that were made to ShouldExecute.
// Check the length with:
//
//	len(mockedTask.ShouldExecuteCalls())
func (mock *TaskMock) ShouldExecuteCalls() []struct {
	Values map[string]any
} {
	var calls []struct {
		Values map[string]any
	}
	mock.lockShouldExecute.RLock()
	calls = mock.calls.ShouldExecute
	mock.lockShouldExecute.RUnlock()
	return calls
}
