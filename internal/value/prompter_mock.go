// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package value

import (
	"sync"
)

// Ensure, that PrompterMock does implement Prompter.
// If this is not the case, regenerate this file with moq.
var _ Prompter = &PrompterMock{}

// PrompterMock is a mock implementation of Prompter.
//
//	func TestSomethingThatUsesPrompter(t *testing.T) {
//
//		// make and configure a mocked Prompter
//		mockedPrompter := &PrompterMock{
//			ConfirmFunc: func(prompt string, defaultValue bool, help string, validationRules string) (bool, error) {
//				panic("mock out the Confirm method")
//			},
//			InputFunc: func(prompt string, defaultValue string, help string, validationRules string) (string, error) {
//				panic("mock out the Input method")
//			},
//			MultiSelectFunc: func(prompt string, options []string, defaultValues []string, help string, validationRules string) ([]string, error) {
//				panic("mock out the MultiSelect method")
//			},
//			SelectFunc: func(prompt string, options []string, defaultValue string, help string, validationRules string) (string, error) {
//				panic("mock out the Select method")
//			},
//		}
//
//		// use mockedPrompter in code that requires Prompter
//		// and then make assertions.
//
//	}
type PrompterMock struct {
	// ConfirmFunc mocks the Confirm method.
	ConfirmFunc func(prompt string, defaultValue bool, help string, validationRules string) (bool, error)

	// InputFunc mocks the Input method.
	InputFunc func(prompt string, defaultValue string, help string, validationRules string) (string, error)

	// MultiSelectFunc mocks the MultiSelect method.
	MultiSelectFunc func(prompt string, options []string, defaultValues []string, help string, validationRules string) ([]string, error)

	// SelectFunc mocks the Select method.
	SelectFunc func(prompt string, options []string, defaultValue string, help string, validationRules string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// Confirm holds details about calls to the Confirm method.
		Confirm []struct {
			// Prompt is the prompt argument value.
			Prompt string
			// DefaultValue is the defaultValue argument value.
			DefaultValue bool
			// Help is the help argument value.
			Help string
			// ValidationRules is the validationRules argument value.
			ValidationRules string
		}
		// Input holds details about calls to the Input method.
		Input []struct {
			// Prompt is the prompt argument value.
			Prompt string
			// DefaultValue is the defaultValue argument value.
			DefaultValue string
			// Help is the help argument value.
			Help string
			// ValidationRules is the validationRules argument value.
			ValidationRules string
		}
		// MultiSelect holds details about calls to the MultiSelect method.
		MultiSelect []struct {
			// Prompt is the prompt argument value.
			Prompt string
			// Options is the options argument value.
			Options []string
			// DefaultValues is the defaultValues argument value.
			DefaultValues []string
			// Help is the help argument value.
			Help string
			// ValidationRules is the validationRules argument value.
			ValidationRules string
		}
		// Select holds details about calls to the Select method.
		Select []struct {
			// Prompt is the prompt argument value.
			Prompt string
			// Options is the options argument value.
			Options []string
			// DefaultValue is the defaultValue argument value.
			DefaultValue string
			// Help is the help argument value.
			Help string
			// ValidationRules is the validationRules argument value.
			ValidationRules string
		}
	}
	lockConfirm     sync.RWMutex
	lockInput       sync.RWMutex
	lockMultiSelect sync.RWMutex
	lockSelect      sync.RWMutex
}

// Confirm calls ConfirmFunc.
func (mock *PrompterMock) Confirm(prompt string, defaultValue bool, help string, validationRules string) (bool, error) {
	if mock.ConfirmFunc == nil {
		panic("PrompterMock.ConfirmFunc: method is nil but Prompter.Confirm was just called")
	}
	callInfo := struct {
		Prompt          string
		DefaultValue    bool
		Help            string
		ValidationRules string
	}{
		Prompt:          prompt,
		DefaultValue:    defaultValue,
		Help:            help,
		ValidationRules: validationRules,
	}
	mock.lockConfirm.Lock()
	mock.calls.Confirm = append(mock.calls.Confirm, callInfo)
	mock.lockConfirm.Unlock()
	return mock.ConfirmFunc(prompt, defaultValue, help, validationRules)
}

// ConfirmCalls gets all the calls that were made to Confirm.
// Check the length with:
//
//	len(mockedPrompter.ConfirmCalls())
func (mock *PrompterMock) ConfirmCalls() []struct {
	Prompt          string
	DefaultValue    bool
	Help            string
	ValidationRules string
} {
	var calls []struct {
		Prompt          string
		DefaultValue    bool
		Help            string
		ValidationRules string
	}
	mock.lockConfirm.RLock()
	calls = mock.calls.Confirm
	mock.lockConfirm.RUnlock()
	return calls
}

// Input calls InputFunc.
func (mock *PrompterMock) Input(prompt string, defaultValue string, help string, validationRules string) (string, error) {
	if mock.InputFunc == nil {
		panic("PrompterMock.InputFunc: method is nil but Prompter.Input was just called")
	}
	callInfo := struct {
		Prompt          string
		DefaultValue    string
		Help            string
		ValidationRules string
	}{
		Prompt:          prompt,
		DefaultValue:    defaultValue,
		Help:            help,
		ValidationRules: validationRules,
	}
	mock.lockInput.Lock()
	mock.calls.Input = append(mock.calls.Input, callInfo)
	mock.lockInput.Unlock()
	return mock.InputFunc(prompt, defaultValue, help, validationRules)
}

// InputCalls gets all the calls that were made to Input.
// Check the length with:
//
//	len(mockedPrompter.InputCalls())
func (mock *PrompterMock) InputCalls() []struct {
	Prompt          string
	DefaultValue    string
	Help            string
	ValidationRules string
} {
	var calls []struct {
		Prompt          string
		DefaultValue    string
		Help            string
		ValidationRules string
	}
	mock.lockInput.RLock()
	calls = mock.calls.Input
	mock.lockInput.RUnlock()
	return calls
}

// MultiSelect calls MultiSelectFunc.
func (mock *PrompterMock) MultiSelect(prompt string, options []string, defaultValues []string, help string, validationRules string) ([]string, error) {
	if mock.MultiSelectFunc == nil {
		panic("PrompterMock.MultiSelectFunc: method is nil but Prompter.MultiSelect was just called")
	}
	callInfo := struct {
		Prompt          string
		Options         []string
		DefaultValues   []string
		Help            string
		ValidationRules string
	}{
		Prompt:          prompt,
		Options:         options,
		DefaultValues:   defaultValues,
		Help:            help,
		ValidationRules: validationRules,
	}
	mock.lockMultiSelect.Lock()
	mock.calls.MultiSelect = append(mock.calls.MultiSelect, callInfo)
	mock.lockMultiSelect.Unlock()
	return mock.MultiSelectFunc(prompt, options, defaultValues, help, validationRules)
}

// MultiSelectCalls gets all the calls that were made to MultiSelect.
// Check the length with:
//
//	len(mockedPrompter.MultiSelectCalls())
func (mock *PrompterMock) MultiSelectCalls() []struct {
	Prompt          string
	Options         []string
	DefaultValues   []string
	Help            string
	ValidationRules string
} {
	var calls []struct {
		Prompt          string
		Options         []string
		DefaultValues   []string
		Help            string
		ValidationRules string
	}
	mock.lockMultiSelect.RLock()
	calls = mock.calls.MultiSelect
	mock.lockMultiSelect.RUnlock()
	return calls
}

// Select calls SelectFunc.
func (mock *PrompterMock) Select(prompt string, options []string, defaultValue string, help string, validationRules string) (string, error) {
	if mock.SelectFunc == nil {
		panic("PrompterMock.SelectFunc: method is nil but Prompter.Select was just called")
	}
	callInfo := struct {
		Prompt          string
		Options         []string
		DefaultValue    string
		Help            string
		ValidationRules string
	}{
		Prompt:          prompt,
		Options:         options,
		DefaultValue:    defaultValue,
		Help:            help,
		ValidationRules: validationRules,
	}
	mock.lockSelect.Lock()
	mock.calls.Select = append(mock.calls.Select, callInfo)
	mock.lockSelect.Unlock()
	return mock.SelectFunc(prompt, options, defaultValue, help, validationRules)
}

// SelectCalls gets all the calls that were made to Select.
// Check the length with:
//
//	len(mockedPrompter.SelectCalls())
func (mock *PrompterMock) SelectCalls() []struct {
	Prompt          string
	Options         []string
	DefaultValue    string
	Help            string
	ValidationRules string
} {
	var calls []struct {
		Prompt          string
		Options         []string
		DefaultValue    string
		Help            string
		ValidationRules string
	}
	mock.lockSelect.RLock()
	calls = mock.calls.Select
	mock.lockSelect.RUnlock()
	return calls
}
