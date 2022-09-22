// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package generate

import (
	"fmt"
	"strings"
)

const (
	// ConflictKeep is a Conflict of type keep.
	ConflictKeep Conflict = "keep"
	// ConflictReplace is a Conflict of type replace.
	ConflictReplace Conflict = "replace"
	// ConflictPrompt is a Conflict of type prompt.
	ConflictPrompt Conflict = "prompt"
)

var _ConflictNames = []string{
	string(ConflictKeep),
	string(ConflictReplace),
	string(ConflictPrompt),
}

// ConflictNames returns a list of possible string values of Conflict.
func ConflictNames() []string {
	tmp := make([]string, len(_ConflictNames))
	copy(tmp, _ConflictNames)
	return tmp
}

// String implements the Stringer interface.
func (x Conflict) String() string {
	return string(x)
}

// String implements the Stringer interface.
func (x Conflict) IsValid() bool {
	_, err := ParseConflict(string(x))
	return err == nil
}

var _ConflictValue = map[string]Conflict{
	"keep":    ConflictKeep,
	"replace": ConflictReplace,
	"prompt":  ConflictPrompt,
}

// ParseConflict attempts to convert a string to a Conflict.
func ParseConflict(name string) (Conflict, error) {
	if x, ok := _ConflictValue[name]; ok {
		return x, nil
	}
	return Conflict(""), fmt.Errorf("%s is not a valid Conflict, try [%s]", name, strings.Join(_ConflictNames, ", "))
}

// MarshalText implements the text marshaller method.
func (x Conflict) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Conflict) UnmarshalText(text []byte) error {
	tmp, err := ParseConflict(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
