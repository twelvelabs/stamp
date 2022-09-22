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
//			ConfirmFunc: func(prompt string, help string, defaultValue bool, validationRules string) (bool, error) {
//				panic("mock out the Confirm method")
//			},
//			InputFunc: func(prompt string, help string, defaultValue string, validationRules string) (string, error) {
//				panic("mock out the Input method")
//			},
//		}
//
//		// use mockedPrompter in code that requires Prompter
//		// and then make assertions.
//
//	}
type PrompterMock struct {
	// ConfirmFunc mocks the Confirm method.
	ConfirmFunc func(prompt string, help string, defaultValue bool, validationRules string) (bool, error)

	// InputFunc mocks the Input method.
	InputFunc func(prompt string, help string, defaultValue string, validationRules string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// Confirm holds details about calls to the Confirm method.
		Confirm []struct {
			// Prompt is the prompt argument value.
			Prompt string
			// Help is the help argument value.
			Help string
			// DefaultValue is the defaultValue argument value.
			DefaultValue bool
			// ValidationRules is the validationRules argument value.
			ValidationRules string
		}
		// Input holds details about calls to the Input method.
		Input []struct {
			// Prompt is the prompt argument value.
			Prompt string
			// Help is the help argument value.
			Help string
			// DefaultValue is the defaultValue argument value.
			DefaultValue string
			// ValidationRules is the validationRules argument value.
			ValidationRules string
		}
	}
	lockConfirm sync.RWMutex
	lockInput   sync.RWMutex
}

// Confirm calls ConfirmFunc.
func (mock *PrompterMock) Confirm(prompt string, help string, defaultValue bool, validationRules string) (bool, error) {
	if mock.ConfirmFunc == nil {
		panic("PrompterMock.ConfirmFunc: method is nil but Prompter.Confirm was just called")
	}
	callInfo := struct {
		Prompt          string
		Help            string
		DefaultValue    bool
		ValidationRules string
	}{
		Prompt:          prompt,
		Help:            help,
		DefaultValue:    defaultValue,
		ValidationRules: validationRules,
	}
	mock.lockConfirm.Lock()
	mock.calls.Confirm = append(mock.calls.Confirm, callInfo)
	mock.lockConfirm.Unlock()
	return mock.ConfirmFunc(prompt, help, defaultValue, validationRules)
}

// ConfirmCalls gets all the calls that were made to Confirm.
// Check the length with:
//
//	len(mockedPrompter.ConfirmCalls())
func (mock *PrompterMock) ConfirmCalls() []struct {
	Prompt          string
	Help            string
	DefaultValue    bool
	ValidationRules string
} {
	var calls []struct {
		Prompt          string
		Help            string
		DefaultValue    bool
		ValidationRules string
	}
	mock.lockConfirm.RLock()
	calls = mock.calls.Confirm
	mock.lockConfirm.RUnlock()
	return calls
}

// Input calls InputFunc.
func (mock *PrompterMock) Input(prompt string, help string, defaultValue string, validationRules string) (string, error) {
	if mock.InputFunc == nil {
		panic("PrompterMock.InputFunc: method is nil but Prompter.Input was just called")
	}
	callInfo := struct {
		Prompt          string
		Help            string
		DefaultValue    string
		ValidationRules string
	}{
		Prompt:          prompt,
		Help:            help,
		DefaultValue:    defaultValue,
		ValidationRules: validationRules,
	}
	mock.lockInput.Lock()
	mock.calls.Input = append(mock.calls.Input, callInfo)
	mock.lockInput.Unlock()
	return mock.InputFunc(prompt, help, defaultValue, validationRules)
}

// InputCalls gets all the calls that were made to Input.
// Check the length with:
//
//	len(mockedPrompter.InputCalls())
func (mock *PrompterMock) InputCalls() []struct {
	Prompt          string
	Help            string
	DefaultValue    string
	ValidationRules string
} {
	var calls []struct {
		Prompt          string
		Help            string
		DefaultValue    string
		ValidationRules string
	}
	mock.lockInput.RLock()
	calls = mock.calls.Input
	mock.lockInput.RUnlock()
	return calls
}
