package value

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/gobuffalo/flect"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/swaggest/jsonschema-go"
	"github.com/twelvelabs/termite/render"
	"github.com/twelvelabs/termite/ui"
	"github.com/twelvelabs/termite/validate"
)

var (
	// ensure Value implements each interface.
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
	if err := validate.Struct(val); err != nil {
		return nil, err
	}

	return val, nil
}

type Value struct {
	// Note: have to use `DataType` because `Type()` is a pflag.Value method.
	Key             string       `mapstructure:"key"                           validate:"required"`
	Name            string       `mapstructure:"name"`
	Flag            string       `mapstructure:"flag"`
	Help            string       `mapstructure:"help"`
	DataType        DataType     `mapstructure:"type"      default:"string"    validate:"required,oneof=bool int intSlice string stringSlice"` //nolint:lll
	Default         any          `mapstructure:"default"`
	PromptConfig    PromptConfig `mapstructure:"prompt"    default:"on-unset"  validate:"required,oneof=always never on-empty on-unset"` //nolint:lll
	InputMode       InputMode    `mapstructure:"mode"      default:"flag"      validate:"required,oneof=arg flag hidden"`
	TransformRules  string       `mapstructure:"transform"`
	ValidationRules string       `mapstructure:"validate"`
	Options         []any        `mapstructure:"options"   nullable:"false"`
	If              string       `mapstructure:"if"        default:"true"`

	data   interface{}
	values *ValueSet
}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (Value) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithTitle("Value")
	schema.WithDescription("A generator input value.")

	schema.Properties["key"].TypeObject.
		WithTitle("Key").
		WithDescription(
			"The variable name for the value. " +
				"This is how you will refer to the value in template files.",
		).
		WithExamples("MyValue")

	schema.Properties["name"].TypeObject.
		WithTitle("Name").
		WithDescription(
			"The display name shown when prompting for the value. " +
				"Will default to a [humanized](https://pkg.go.dev/github.com/gobuffalo/flect#Humanize) " +
				"form of the [key](https://github.com/twelvelabs/stamp/tree/main/docs/value.md#key).",
		).
		WithExamples("Custom display name")

	schema.Properties["flag"].TypeObject.
		WithTitle("Flag").
		WithDescription(
			"The flag name for the value. " +
				"Will default to a [dash separated](https://pkg.go.dev/github.com/gobuffalo/flect#Dasherize) " +
				"form of the [key](https://github.com/twelvelabs/stamp/tree/main/docs/value.md#key).",
		).
		WithExamples("custom-flag-name")

	schema.Properties["help"].TypeObject.
		WithTitle("Help").
		WithDescription(
			"Help text describing the value. " +
				"Shown when prompting and when using the `--help` flag.",
		).
		WithExamples("You should enter a random value.")

	schema.Properties["default"].TypeObject.
		WithTitle("Default").
		WithDescription(
			"The value default. " +
				"Can refer to other values defined earlier in the list.",
		).
		WithExamples("{{ .OtherValue | underscore }}.txt")

	schema.Properties["transform"].TypeObject.
		WithTitle("Transform").
		WithDescription(
			"Optional, comma-separated list of " +
				"[transform](https://github.com/twelvelabs/stamp/tree/main/docs/transform.md) rules.",
		).
		WithExamples("trim,uppercase")

	schema.Properties["validate"].TypeObject.
		WithTitle("Validate").
		WithDescription(
			"Optional, comma-separated list of " +
				"[validation](https://github.com/go-playground/validator#baked-in-validations) rules.",
		).
		WithExamples("required,email")

	schema.Properties["options"].TypeObject.
		WithTitle("Options").
		WithDescription(
			"A fixed set of valid options for the value. " +
				"Will cause the value to be rendered as a single or " +
				"multi-select when prompted (depending on data type). " +
				"Attempts to assign a value not in this list will " +
				"raise a validation error.",
		).
		WithExamples([]string{
			"foo",
			"bar",
		})

	schema.Properties["if"].TypeObject.
		WithTitle("If").
		WithDescription(
			"Determines whether the value is enabled. "+
				"Can refer to other values defined earlier in the list "+
				"(allows for dynamic prompts).",
		).
		WithExamples(
			"{{ .UseDatabase }}",
			"{{ eq .Language \"python\" }}",
		)

	return nil
}

// DisplayName returns the human friendly display name.
func (v *Value) DisplayName() string {
	if v.Name != "" {
		return v.Name
	}
	return flect.Humanize(v.Key)
}

// FlagName returns the kebab-cased flag name.
func (v *Value) FlagName() string {
	if v.Flag != "" {
		return flect.Dasherize(v.Flag)
	}
	return flect.Dasherize(v.Key)
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

// IsHidden returns true if InputMode is "hidden".
func (v *Value) IsHidden() bool {
	return v.InputMode == InputModeHidden
}

// IsUnset returns true if a value has not been explicitly set.
func (v *Value) IsUnset() bool {
	return v.data == nil
}

// IsEnabled returns true when the `If` field evaluates to true.
// The user is only prompted for enabled values.
func (v *Value) IsEnabled() bool {
	ok, _ := render.String(v.If, v.ValueSet().Cache())
	return cast.ToBool(ok)
}

// IsEmpty returns true if the value is empty.
func (v *Value) IsEmpty() bool {
	rv := reflect.ValueOf(v.Get())
	switch rv.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(rv.String())) == 0
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return rv.Len() == 0
	// Uncomment if we ever support complex value types
	// case reflect.Ptr, reflect.Interface, reflect.Func:
	//	return rv.IsNil()
	default:
		return !rv.IsValid() || reflect.DeepEqual(rv.Interface(), reflect.Zero(rv.Type()).Interface())
	}
}

// ShouldPrompt returns true if the user should be prompted for a value.
func (v *Value) ShouldPrompt() bool {
	if !v.IsEnabled() {
		return false
	}
	switch v.PromptConfig {
	case PromptConfigAlways:
		return true
	case PromptConfigNever:
		return false
	case PromptConfigOnEmpty:
		return v.IsEmpty() && !v.IsHidden()
	case PromptConfigOnUnset:
		return v.IsUnset() && !v.IsHidden()
	default:
		return v.IsUnset() && !v.IsHidden()
	}
}

// Prompt prompts the user for a value.
func (v *Value) Prompt(prompter ui.Prompter) error {
	if !v.ShouldPrompt() {
		return nil
	}

	options := cast.ToStringSlice(v.Options)

	var response interface{}
	var err error
	switch v.DataType {
	case DataTypeBool:
		defVal := cast.ToBool(v.Get())
		response, err = prompter.Confirm(v.DisplayName(), defVal,
			ui.WithHelp(v.Help), ui.WithValidation(v.ValidationRules))
	case DataTypeInt, DataTypeString:
		if len(options) > 0 {
			response, err = prompter.Select(v.DisplayName(), options, v.String(),
				ui.WithHelp(v.Help), ui.WithValidation(v.ValidationRules))
		} else {
			response, err = prompter.Input(v.DisplayName(), v.String(),
				ui.WithHelp(v.Help), ui.WithValidation(v.ValidationRules))
		}
	case DataTypeIntSlice, DataTypeStringSlice:
		if len(options) > 0 {
			defVal := cast.ToStringSlice(v.Get())
			response, err = prompter.MultiSelect(v.DisplayName(), options, defVal,
				ui.WithHelp(v.Help), ui.WithValidation(v.ValidationRules))
		} else {
			response, err = prompter.Input(v.DisplayName(), v.String(),
				ui.WithHelp(v.Help), ui.WithValidation(v.ValidationRules))
		}
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
	case DataTypeIntSlice, DataTypeStringSlice:
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

// WithValueCache sets dm and returns the receiver.
// Should only be used in tests.
func (v *Value) WithValueCache(dm DataMap) *Value {
	v.ValueSet().SetCache(dm)
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
	// Updating the cache (even on get) to help maintain freshness.
	// There are some corner cases where could render stale data,
	// and this prevents _most_ of them :grimacing:.
	v.ValueSet().Cache().Set(v.Key, processed)
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
	v.ValueSet().Cache().Set(v.Key, processed)
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
	str, ok := data.(string)
	if !ok {
		return data, nil
	}
	rendered, err := render.String(str, v.ValueSet().Cache())
	if err != nil {
		return data, err
	}
	return rendered, nil
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
	case DataTypeStringSlice:
		return cast.ToStringSliceE(coerceToSlice(data))
	default:
		return data, ErrInvalidDataType
	}
}

// Passes data through any configured transform rules.
func (v *Value) transform(data any) (any, error) {
	return Transform(v.Key, data, v.TransformRules)
}

// Passes data through any configured validation rules.
func (v *Value) validate(data any) error {
	rules := v.ValidationRules
	if len(v.Options) > 0 {
		// Ensure a validation rule for options (saves people from having to do so manually).
		// Appends a rule like: "oneof=foo bar baz" to the end of any existing rules.
		var rule string
		switch v.DataType {
		case DataTypeIntSlice, DataTypeStringSlice:
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
	return validate.KeyVal(v.Key, data, rules)
}
