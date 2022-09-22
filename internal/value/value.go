package value

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	//cspell:words pflag xstrings oneof
	//cspell:disable
	"github.com/Masterminds/sprig/v3"
	"github.com/creasty/defaults"
	"github.com/huandu/xstrings"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	//cspell:enable
)

var (
	ErrInvalidDataType = errors.New("invalid data type")

	// ensure Value implements each interface
	_ flag.Getter = &Value{}
	_ pflag.Value = &Value{}
)

// NewValue returns a new Value struct for the given map of data.
func NewValue(valueData map[string]any) (*Value, error) {
	val := &Value{}

	if err := defaults.Set(val); err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(valueData, val); err != nil {
		return nil, err
	}
	if err := ValidateStruct(val); err != nil {
		return nil, err
	}

	return val, nil
}

type Value struct {
	// Note: have to use `DataType` because `Type()` is a pflag.Value method.
	Name            string       `mapstructure:"name"                          validate:"required"`
	Help            string       `mapstructure:"help"`
	DataType        DataType     `mapstructure:"type"      default:"string"    validate:"required,oneof=bool int intSlice string"`
	Default         interface{}  `mapstructure:"default"`
	PromptConfig    PromptConfig `mapstructure:"prompt"    default:"on-unset"  validate:"required,oneof=always never on-empty on-unset"`
	InputMode       InputMode    `mapstructure:"mode"      default:"flag"      validate:"required,oneof=arg flag"`
	TransformRules  string       `mapstructure:"transform"`
	ValidationRules string       `mapstructure:"validate"`
	Options         []any        `mapstructure:"options"   default:"[]"`

	data   interface{}
	values *ValueSet
}

// The camel-cased name of value.
func (v *Value) CamelName() string {
	return xstrings.ToCamelCase(v.Name)
}

// The kebab-cased name of value.
func (v *Value) KebabName() string {
	return xstrings.ToKebabCase(v.Name)
}

// Get returns the rendered, casted value.
// Required to implement [flag.Getter] interface.
func (v *Value) Get() any {
	data, _ := v.get()
	return data
}

// IsBoolFlag returns true if the data type is `bool`.
// Required to implement [pflag.boolFlag] interface.
func (v *Value) IsBoolFlag() bool {
	return v.DataType == DataTypeBool
}

// IsArg returns true if InputMode is "arg".
func (v *Value) IsArg() bool {
	return v.InputMode == InputModeArg
}

// IsFlag returns true if InputMode is "flag".
func (v *Value) IsFlag() bool {
	return v.InputMode == InputModeFlag
}

// IsUnset returns true if a value has not been explicitly set.
func (v *Value) IsUnset() bool {
	return v.data == nil
}

// IsEmpty returns true if the
func (v *Value) IsEmpty() bool {
	rv := reflect.ValueOf(v.Get())
	switch rv.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(rv.String())) == 0
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return rv.Len() == 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return rv.IsNil()
	default:
		return !rv.IsValid() || reflect.DeepEqual(rv.Interface(), reflect.Zero(rv.Type()).Interface())
	}
}

// Key returns the materialized data key for the value.
func (v *Value) Key() string {
	// TODO: use inflection lib to properly "modularize"
	return v.CamelName()
}

// ShouldPrompt returns true if the user should be prompted for a value.
func (v *Value) ShouldPrompt() bool {
	switch v.PromptConfig {
	case PromptConfigAlways:
		return true
	case PromptConfigNever:
		return false
	case PromptConfigOnEmpty:
		return v.IsEmpty()
	case PromptConfigOnUnset:
		return v.IsUnset()
	default:
		return v.IsUnset()
	}
}

// Prompt prompts the user for a value.
func (v *Value) Prompt(prompter Prompter) error {
	var response interface{}
	var err error

	if !v.ShouldPrompt() {
		return nil
	}

	switch v.DataType {
	case DataTypeBool:
		response, err = prompter.Confirm(v.Name, cast.ToBool(v.Get()), v.Help, v.ValidationRules)
	case DataTypeInt, DataTypeString:
		if len(v.Options) > 0 {
			options := cast.ToStringSlice(v.Options)
			response, err = prompter.Select(v.Name, options, v.String(), v.Help, v.ValidationRules)
		} else {
			response, err = prompter.Input(v.Name, v.String(), v.Help, v.ValidationRules)
		}
	case DataTypeIntSlice:
		response, err = prompter.Input(v.Name, v.String(), v.Help, v.ValidationRules)
	default:
		return ErrInvalidDataType
	}
	if err != nil {
		return err // prompt error
	}

	err = v.set(response)
	if err != nil {
		return err // set error
	}
	return nil
}

// Set sets the value.
// Returns an error if the input data can not be casted to the correct type.
// Required to implement the [pflag.Value] interface.
func (v *Value) Set(data string) error {
	return v.set(data)
}

// Required to implement the [pflag.Value] interface.
func (v *Value) String() string {
	switch v.DataType {
	case DataTypeIntSlice:
		return strings.Join(cast.ToStringSlice(v.Get()), ",")
	default:
		return cast.ToString(v.Get())
	}
}

// Required to implement the [pflag.Value] interface.
func (v *Value) Type() string {
	return v.DataType.String()
}

// Validate evaluates the configured validation rules.
func (v *Value) Validate() error {
	return v.validate(v.Get())
}

func (v *Value) ValueSet() *ValueSet {
	if v.values == nil {
		v.values = NewValueSet()
	}
	return v.values
}

// WithDataMap sets dm and returns the receiver.
func (v *Value) WithDataMap(dm DataMap) *Value {
	v.ValueSet().SetDataMap(dm)
	return v
}

// WithValueSet sets vs and returns the receiver.
func (v *Value) WithValueSet(vs *ValueSet) *Value {
	v.values = vs
	return v
}

func (v *Value) get() (any, error) {
	data := v.data
	if data == nil {
		data = v.Default
	}
	processed, err := v.process(data)
	if err != nil {
		return nil, err
	}
	return processed, nil
}

func (v *Value) set(data any) error {
	processed, err := v.process(data)
	if err != nil {
		return err
	}
	if err := v.validate(processed); err != nil {
		return err
	}
	v.data = processed
	v.ValueSet().SetData(v.Key(), processed)
	return nil
}

// Passes data through the render/cast/transform pipeline.
func (v *Value) process(data any) (any, error) {
	rendered, err := v.render(data)
	if err != nil {
		return nil, err
	}
	casted, err := v.cast(rendered)
	if err != nil {
		return nil, err
	}
	transformed, err := v.transform(casted)
	if err != nil {
		return nil, err
	}
	return transformed, nil
}

// Attempts to render the data as a [text/template].
// If data is not renderable, returns it as-is.
func (v *Value) render(data any) (any, error) {
	rawStr, ok := data.(string)
	if !ok {
		return data, nil
	}

	tmp, err := template.New("render").Funcs(sprig.FuncMap()).Parse(rawStr)
	if err != nil {
		return data, err
	}

	buf := &bytes.Buffer{}
	err = tmp.Execute(buf, v.ValueSet().GetDataMap())
	if err != nil {
		return data, err
	}

	return buf.String(), nil
}

// cast converts data to the values type.
func (v *Value) cast(data any) (any, error) {
	var casted any
	var err error

	csvParse := func(s string) []string {
		if s == "" {
			return []string{} // empty string; empty slice
		}
		segments := strings.Split(s, ",")
		for i, segment := range segments {
			segments[i] = strings.TrimSpace(segment)
		}
		return segments
	}

	// Incoming data for slice-types _may_ be comma separated strings
	// (depending on whether cast is being called by `get` or `set`).
	// Try to massage the data into something that can be handled by the cast func.
	coerceToSlice := func(data any) any {
		if str, ok := data.(string); ok {
			return csvParse(str)
		} else if data == nil {
			return []string{} // some of the cast functions can't handle nil
		}
		return data // :shrug:
	}

	switch v.DataType {
	case DataTypeBool:
		casted, err = cast.ToBoolE(data)
		if err != nil {
			// simplify the error message
			return casted, errors.New("unable to cast to bool")
		}
		return casted, nil
	case DataTypeInt:
		return cast.ToIntE(data)
	case DataTypeIntSlice:
		return cast.ToIntSliceE(coerceToSlice(data))
	case DataTypeString:
		return cast.ToStringE(data)
	default:
		return data, ErrInvalidDataType
	}
}

// Passes data through any configured transform rules.
func (v *Value) transform(data any) (any, error) {
	return Transform(v.Name, data, v.TransformRules)
}

// Passes data through any configured validation rules.
func (v *Value) validate(data any) error {
	rules := v.ValidationRules
	if len(v.Options) > 0 {
		// Ensure a validation rule for options (saves people from having to do so manually).
		// Appends a rule like: "oneof=foo bar baz" to the end of any existing rules.
		var rule string
		switch v.DataType {
		case DataTypeIntSlice:
			rule = "dive,"
		default:
			rule = ""
		}
		opts := cast.ToStringSlice(v.Options)
		rule += fmt.Sprintf("oneof=%s", strings.Join(opts, " "))
		if rules == "" {
			rules = rule
		} else {
			segments := strings.Split(rules, ",")
			segments = append(segments, rule)
			rules = strings.Join(segments, ",")
		}
	}
	return ValidateKeyVal(v.Name, data, rules)
}
