package value

//go:generate moq -rm -out prompter_mock.go . Prompter
type Prompter interface {
	// Prompt for a boolean yes/no value.
	Confirm(
		prompt string, defaultValue bool, help string, validationRules string,
	) (bool, error)
	// Prompt for single string value.
	Input(
		prompt string, defaultValue string, help string, validationRules string,
	) (string, error)
	// Prompt for a slice of string values w/ a fixed set of options.
	MultiSelect(
		prompt string, options []string, defaultValues []string, help string, validationRules string,
	) ([]string, error)
	// Prompt for single string value w/ a fixed set of options.
	Select(
		prompt string, options []string, defaultValue string, help string, validationRules string,
	) (string, error)
}

type ConfirmFunc func(
	prompt string, defaultValue bool, help string, validationRules string,
) (bool, error)
type InputFunc func(
	prompt string, defaultValue string, help string, validationRules string,
) (string, error)
type MultiSelectFunc func(
	prompt string, options []string, defaultValues []string, help string, validationRules string,
) ([]string, error)
type SelectFunc func(
	prompt string, options []string, defaultValue string, help string, validationRules string,
) (string, error)

// Creates a new ConfirmFunc that returns the given result and err.
func NewConfirmFunc(result bool, err error) ConfirmFunc {
	return func(prompt string, defaultValue bool, help string, validationRules string) (bool, error) {
		return result, err
	}
}

// Creates a new ConfirmFunc that returns the default value.
func NewNoopConfirmFunc() ConfirmFunc {
	return func(prompt string, defaultValue bool, help string, validationRules string) (bool, error) {
		return defaultValue, nil
	}
}

// Creates a new InputFunc that returns the given result and err.
func NewInputFunc(result string, err error) InputFunc {
	return func(prompt string, defaultValue string, help string, validationRules string) (string, error) {
		return result, err
	}
}

// Creates a new InputFunc that returns the default value.
func NewNoopInputFunc() InputFunc {
	return func(prompt string, defaultValue string, help string, validationRules string) (string, error) {
		return defaultValue, nil
	}
}

// Creates a new SelectFunc that returns the given result and err.
func NewMultiSelectFunc(result []string, err error) MultiSelectFunc {
	return func(
		prompt string, options []string, defaultValues []string, help string, validationRules string,
	) ([]string, error) {
		return result, err
	}
}

// Creates a new InputFunc that returns the default value.
func NewNoopMultiSelectFunc() MultiSelectFunc {
	return func(
		prompt string, options []string, defaultValues []string, help string, validationRules string,
	) ([]string, error) {
		return defaultValues, nil
	}
}

// Creates a new SelectFunc that returns the given result and err.
func NewSelectFunc(result string, err error) SelectFunc {
	return func(
		prompt string, options []string, defaultValue string, help string, validationRules string,
	) (string, error) {
		return result, err
	}
}

// Creates a new InputFunc that returns the default value.
func NewNoopSelectFunc() SelectFunc {
	return func(
		prompt string, options []string, defaultValue string, help string, validationRules string,
	) (string, error) {
		return defaultValue, nil
	}
}
