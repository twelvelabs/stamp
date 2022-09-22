package value

//go:generate moq -rm -out prompter_mock.go . Prompter
type Prompter interface {
	// Prompt for single string value.
	Input(prompt string, help string, defaultValue string, validationRules string) (string, error)
	// Prompt for a boolean yes/no value.
	Confirm(prompt string, help string, defaultValue bool, validationRules string) (bool, error)
}

type ConfirmFunc func(prompt string, help string, defaultValue bool, validationRules string) (bool, error)
type InputFunc func(prompt string, help string, defaultValue string, validationRules string) (string, error)

// Creates a new ConfirmFunc that returns the given result and err.
func NewConfirmFunc(result bool, err error) ConfirmFunc {
	return func(prompt, help string, defaultValue bool, validationRules string) (bool, error) {
		return result, err
	}
}

// Creates a new ConfirmFunc that returns the default value.
func NewNoopConfirmFunc() ConfirmFunc {
	return func(prompt, help string, defaultValue bool, validationRules string) (bool, error) {
		return defaultValue, nil
	}
}

// Creates a new InputFunc that returns the given result and err.
func NewInputFunc(result string, err error) InputFunc {
	return func(prompt, help, defaultValue string, validationRules string) (string, error) {
		return result, err
	}
}

// Creates a new InputFunc that returns the default value.
func NewNoopInputFunc() InputFunc {
	return func(prompt, help, defaultValue string, validationRules string) (string, error) {
		return defaultValue, nil
	}
}
