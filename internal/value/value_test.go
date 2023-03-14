package value

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/ui"
)

func TestNewValue(t *testing.T) {
	tests := []struct {
		Name   string
		Data   map[string]any
		Output *Value
		Err    string
	}{
		{
			Name: "only requires a name and sets correct default values for other fields",
			Data: map[string]any{
				"key": "foo",
			},
			Output: &Value{
				Key:          "foo",
				DataType:     DataTypeString,
				PromptConfig: PromptConfigOnUnset,
				InputMode:    InputModeFlag,
				Options:      []any{},
			},
			Err: "",
		},
		{
			Name:   "returns an error when name is missing from the data map",
			Data:   map[string]any{},
			Output: nil,
			Err:    "Key is a required field",
		},
		{
			Name: "returns an error invalid data types are in the map",
			Data: map[string]any{
				"key": 123,
			},
			Output: nil,
			Err:    "'key' expected type 'string', got unconvertible type 'int'",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual, err := NewValue(test.Data)
			assert.Equal(t, test.Output, actual)
			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}

func TestValue_WithValueSet(t *testing.T) {
	vs := NewValueSet()
	v := &Value{}
	assert.Equal(t, vs, v.WithValueSet(vs).ValueSet())
}

func TestValue_DisplayName(t *testing.T) {
	tests := []struct {
		Desc   string
		Key    string
		Name   string
		Output string
	}{
		{
			Desc:   "should return a humanized version of key if no name provided",
			Key:    "FooBar",
			Output: "Foo Bar",
		},
		{
			Desc:   "should return name if provided",
			Key:    "FooBar",
			Name:   "Something else",
			Output: "Something else",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			value := &Value{
				Key:  tt.Key,
				Name: tt.Name,
			}
			assert.Equal(t, tt.Output, value.DisplayName())
		})
	}
}

func TestValue_FlagName(t *testing.T) {
	tests := []struct {
		Desc   string
		Key    string
		Flag   string
		Output string
	}{
		{
			Desc:   "should return a kebabcase version of key if no flag name provided",
			Key:    "FooBar",
			Output: "foo-bar",
		},
		{
			Desc:   "should return flag name if provided",
			Key:    "FooBar",
			Flag:   "something-else",
			Output: "something-else",
		},
		{
			Desc:   "should ensure flag name is kebabcase even if provided",
			Key:    "FooBar",
			Flag:   "Something Wrong",
			Output: "something-wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			value := &Value{
				Key:  tt.Key,
				Flag: tt.Flag,
			}
			assert.Equal(t, tt.Output, value.FlagName())
		})
	}
}

func TestValue_IsBoolFlag(t *testing.T) {
	assert.Equal(t, false, (&Value{DataType: DataTypeString}).IsBoolFlag())
	assert.Equal(t, true, (&Value{DataType: DataTypeBool}).IsBoolFlag())
}

func TestValue_Type(t *testing.T) {
	assert.Equal(t, "string", (&Value{DataType: DataTypeString}).Type())
	assert.Equal(t, "bool", (&Value{DataType: DataTypeBool}).Type())
}

func TestValue_IsEmpty(t *testing.T) {
	tests := []struct {
		Name    string
		Value   *Value
		IsEmpty bool
	}{
		{
			Name: "[bool] false is empty",
			Value: &Value{
				DataType: DataTypeBool,
				Default:  false,
			},
			IsEmpty: true,
		},
		{
			Name: "[bool] true is non-empty",
			Value: &Value{
				DataType: DataTypeBool,
				Default:  true,
			},
			IsEmpty: false,
		},

		{
			Name: "[int] 0 is empty",
			Value: &Value{
				DataType: DataTypeInt,
				Default:  0,
			},
			IsEmpty: true,
		},
		{
			Name: "[int] gt 0 is non-empty",
			Value: &Value{
				DataType: DataTypeInt,
				Default:  12,
			},
			IsEmpty: false,
		},

		{
			Name: "[intSlice] empty slice is empty",
			Value: &Value{
				DataType: DataTypeIntSlice,
				Default:  []int{},
			},
			IsEmpty: true,
		},
		{
			Name: "[intSlice] non-empty slice is non-empty",
			Value: &Value{
				DataType: DataTypeIntSlice,
				Default:  []int{1, 2, 3},
			},
			IsEmpty: false,
		},

		{
			Name: "[string] empty string is empty",
			Value: &Value{
				DataType: DataTypeString,
				Default:  "",
			},
			IsEmpty: true,
		},
		{
			Name: "[string] non-empty string is non-empty",
			Value: &Value{
				DataType: DataTypeString,
				Default:  "foo",
			},
			IsEmpty: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.IsEmpty, test.Value.IsEmpty())
		})
	}
}

func TestValue_GetAndSet(t *testing.T) { //nolint: maintidx
	RegisterTransformer("explode", func(a any) (any, error) {
		return "lol", errors.New("boom")
	})

	tests := []struct {
		Name   string
		Value  *Value
		Input  string
		Output any
		String string
		Err    string
	}{
		{
			Name: "raises when type is empty",
			Value: (&Value{
				DataType: "",
			}),
			Input:  "something",
			Output: nil,
			Err:    "not a valid DataType",
		},
		{
			Name: "raises when type is invalid",
			Value: (&Value{
				DataType: "wat",
			}),
			Input:  "something",
			Output: nil,
			Err:    "not a valid DataType",
		},

		// BOOL
		{
			Name: "[bool] accepts valid input",
			Value: (&Value{
				DataType: "bool",
				Default:  true,
			}),
			Input:  "false",
			Output: false,
			String: "false",
			Err:    "",
		},
		{
			Name: "[bool] returns default value if no input given",
			Value: (&Value{
				DataType: "bool",
				Default:  "{{.WannaDance}}",
			}).WithValueCache(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Input:  "",
			Output: true,
			String: "true",
			Err:    "",
		},
		{
			Name: "[bool] errors on invalid value",
			Value: (&Value{
				DataType: "bool",
				Default:  true,
			}),
			Input:  "wat",
			Output: true,
			String: "true",
			Err:    "unable to cast to bool",
		},

		// INT
		{
			Name: "[int] accepts valid input",
			Value: (&Value{
				DataType: "int",
				Default:  10,
			}),
			Input:  "25",
			Output: 25,
			String: "25",
			Err:    "",
		},
		{
			Name: "[int] returns default value if no input given",
			Value: (&Value{
				DataType: "int",
				Default:  "{{ add .Year 1 }}",
			}).WithValueCache(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Input:  "",
			Output: 1978,
			String: "1978",
			Err:    "",
		},
		{
			Name: "[int] errors when value not in options",
			Value: (&Value{
				DataType: "int",
				Default:  1,
				Options:  []any{1, 2, 3},
			}),
			Input:  "25",
			Output: 1,
			String: "1",
			Err:    "must be one of [1 2 3]",
		},
		{
			Name: "[int] errors on invalid value",
			Value: (&Value{
				DataType: "int",
				Default:  10,
			}),
			Input:  "wat",
			Output: 10,
			String: "10",
			Err:    "unable to cast",
		},

		{
			Name: "[intSlice] accepts single value input",
			Value: (&Value{
				DataType: "intSlice",
				Default:  []int{1, 2, 3},
			}),
			Input:  "4",
			Output: []int{4},
			String: "4",
			Err:    "",
		},
		{
			Name: "[intSlice] accepts comma separated input",
			Value: (&Value{
				DataType: "intSlice",
				Default:  []int{1, 2, 3},
			}),
			Input:  "4, 5, 6",
			Output: []int{4, 5, 6},
			String: "4,5,6",
			Err:    "",
		},
		{
			Name: "[intSlice] converts nil values to empty slice",
			Value: (&Value{
				DataType: "intSlice",
				Default:  nil,
			}),
			Input:  "",
			Output: []int{},
			String: "",
			Err:    "",
		},
		{
			Name: "[intSlice] returns default value if no input given",
			Value: (&Value{
				DataType: "intSlice",
				Default:  "{{ add .Year 1 }},{{ add .Year 2 }}",
			}).WithValueCache(DataMap{
				"Year": 1977,
			}),
			Input:  "",
			Output: []int{1978, 1979},
			String: "1978,1979",
			Err:    "",
		},
		{
			Name: "[intSlice] errors when value not in options",
			Value: (&Value{
				DataType: "intSlice",
				Default:  []int{1, 2},
				Options:  []any{1, 2, 3},
			}),
			Input:  "3, 4",
			Output: []int{1, 2},
			String: "1,2",
			Err:    "must be one of [1 2 3]",
		},
		{
			Name: "[intSlice] errors on invalid value",
			Value: (&Value{
				DataType: "intSlice",
				Default:  []int{1, 2},
			}),
			Input:  "foo, bar",
			Output: []int{1, 2},
			String: "1,2",
			Err:    "unable to cast",
		},

		// STRING
		{
			Name: "[string] accepts valid input",
			Value: (&Value{
				DataType: "string",
				Default:  "foo",
			}),
			Input:  "bar",
			Output: "bar",
			String: "bar",
			Err:    "",
		},
		{
			Name: "[string] errors when value not in options",
			Value: (&Value{
				DataType: "string",
				Default:  "foo",
				Options:  []any{"foo", "bar", "baz"},
			}),
			Input:  "nope",
			Output: "foo",
			String: "foo",
			Err:    "must be one of [foo bar baz]",
		},
		{
			Name: "[string] option validation does not interfere with other rules",
			Value: (&Value{
				DataType:        "string",
				Default:         "foo",
				Options:         []any{"foo", "bar", "baz"},
				ValidationRules: "alpha",
			}),
			Input:  "123",
			Output: "foo",
			String: "foo",
			Err:    "can only contain alphabetic characters",
		},
		{
			Name: "[string] renders default values",
			Value: (&Value{
				DataType: "string",
				Default:  "{{.First}} {{.Last}}",
			}).WithValueCache(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Input:  "",
			Output: "Joey Ramone",
			String: "Joey Ramone",
			Err:    "",
		},
		{
			Name: "[string] renders submitted values",
			Value: (&Value{
				DataType: "string",
				Default:  "{{.First}} {{.Last}}",
			}).WithValueCache(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Input:  "Dee Dee {{.Last}}", //nolint: dupword
			Output: "Dee Dee Ramone",    //nolint: dupword
			String: "Dee Dee Ramone",    //nolint: dupword
			Err:    "",
		},
		{
			Name: "[string] returns an empty value on template parse error",
			Value: (&Value{
				DataType: "string",
				Default:  "{{.}",
			}),
			Input:  "",
			Output: nil,
			String: "",
			Err:    "",
		},
		{
			Name: "[string] returns an empty value on template render error",
			Value: (&Value{
				DataType: "string",
				Default:  `{{ fail "boom" }}`,
			}),
			Input:  "",
			Output: nil,
			String: "",
			Err:    "",
		},
		{
			Name: "[string] transforms submitted values",
			Value: (&Value{
				DataType:       "string",
				Default:        "",
				TransformRules: "trim, dasherize, uppercase",
			}),
			Input:  "     hello world     ",
			Output: "HELLO-WORLD",
			String: "HELLO-WORLD",
			Err:    "",
		},
		{
			Name: "[string] transforms default values",
			Value: (&Value{
				DataType:       "string",
				Default:        "{{.First}} {{.Last}}",
				TransformRules: "uppercase",
			}).WithValueCache(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Input:  "",
			Output: "JOEY RAMONE",
			String: "JOEY RAMONE",
			Err:    "",
		},
		{
			Name: "[string] returns an empty value on transform error",
			Value: (&Value{
				DataType:       "string",
				Default:        "",
				TransformRules: "explode",
			}),
			Input:  "howdy",
			Output: nil,
			String: "",
			Err:    "boom",
		},

		{
			Name: "[stringSlice] accepts single value input",
			Value: (&Value{
				DataType: "stringSlice",
				Default:  []string{"foo", "bar"},
			}),
			Input:  "baz",
			Output: []string{"baz"},
			String: "baz",
			Err:    "",
		},
		{
			Name: "[stringSlice] accepts comma separated input",
			Value: (&Value{
				DataType: "stringSlice",
				Default:  []string{"foo"},
			}),
			Input:  "bar, baz",
			Output: []string{"bar", "baz"},
			String: "bar,baz",
			Err:    "",
		},
		{
			Name: "[stringSlice] returns default value if no input given",
			Value: (&Value{
				DataType: "stringSlice",
				Default:  "{{ .First }},{{ .Last }}",
			}).WithValueCache(DataMap{
				"First": "Joey",
				"Last":  "Ramone",
			}),
			Input:  "",
			Output: []string{"Joey", "Ramone"},
			String: "Joey,Ramone",
			Err:    "",
		},
		{
			Name: "[stringSlice] errors when value not in options",
			Value: (&Value{
				DataType: "stringSlice",
				Default:  []string{"foo"},
				Options:  []any{"foo", "bar"},
			}),
			Input:  "bar, baz",
			Output: []string{"foo"},
			String: "foo",
			Err:    "must be one of [foo bar]",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var err error

			value := test.Value
			if test.Input != "" {
				err = value.Set(test.Input)
			}

			assert.Equal(t, test.Output, value.Get())
			assert.Equal(t, test.String, value.String())

			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}

func TestValue_Prompt(t *testing.T) { //nolint: maintidx
	tests := []struct {
		Name   string
		Value  *Value
		Setup  func(p *ui.UserInterface)
		Output interface{}
		Err    string
	}{
		{
			Name: "invalid type",
			Value: (&Value{
				DataType: "not-a-type",
			}),
			Output: nil,
			Err:    "not a valid DataType",
		},
		{
			Name: "[bool] true",
			Value: (&Value{
				DataType: "bool",
				Default:  false,
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchConfirm(""),
					ui.RespondBool(true),
				)
			},
			Output: true,
			Err:    "",
		},
		{
			Name: "[bool] default",
			Value: (&Value{
				DataType: "bool",
				Default:  "{{.WannaDance}}",
			}).WithValueCache(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchConfirm(""),
					ui.RespondDefault(),
				)
			},
			Output: true,
			Err:    "",
		},
		{
			Name: "[bool] error",
			Value: (&Value{
				DataType: "bool",
				Default:  false,
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchConfirm(""),
					ui.RespondError(errors.New("boom")),
				)
			},
			Output: false,
			Err:    "boom",
		},
		{
			Name: "[bool] prompt disabled",
			Value: (&Value{
				DataType:     "bool",
				Default:      false,
				PromptConfig: PromptConfigNever,
			}),
			Output: false,
			Err:    "",
		},

		{
			Name: "[int]",
			Value: (&Value{
				DataType: "int",
				Default:  1,
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchInput(""),
					ui.RespondString("12"),
				)
			},
			Output: 12,
			Err:    "",
		},
		{
			Name: "[int] with options",
			Value: (&Value{
				DataType: "int",
				Default:  1,
				Options:  []any{1, 2, 3},
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchSelect(""),
					ui.RespondString("3"),
				)
			},
			Output: 3,
			Err:    "",
		},
		{
			Name: "[int] default",
			Value: (&Value{
				DataType: "int",
				Default:  "{{ add .Year 1 }}",
			}).WithValueCache(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchInput(""),
					ui.RespondDefault(),
				)
			},
			Output: 1978,
			Err:    "",
		},
		{
			Name: "[int] invalid",
			Value: (&Value{
				DataType: "int",
				Default:  1,
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchInput(""),
					ui.RespondString("not an int"),
				)
			},
			Output: nil,
			Err:    "unable to cast",
		},

		{
			Name: "[intSlice]",
			Value: (&Value{
				DataType: "intSlice",
				Default:  "",
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchInput(""),
					ui.RespondString("1,2,3"),
				)
			},
			Output: []int{1, 2, 3},
			Err:    "",
		},
		{
			Name: "[intSlice] with options",
			Value: (&Value{
				DataType: "intSlice",
				Default:  "",
				Options:  []any{1, 2, 3},
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchMultiSelect(""),
					ui.RespondStringSlice([]string{"1", "2"}),
				)
			},
			Output: []int{1, 2},
			Err:    "",
		},
		{
			Name: "[intSlice] with options default",
			Value: (&Value{
				DataType: "intSlice",
				Default:  []int{1, 2},
				Options:  []any{1, 2, 3},
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchMultiSelect(""),
					ui.RespondDefault(),
				)
			},
			Output: []int{1, 2},
			Err:    "",
		},

		{
			Name: "[string]",
			Value: (&Value{
				DataType: "string",
				Default:  "foo",
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchInput(""),
					ui.RespondString("bar"),
				)
			},
			Output: "bar",
			Err:    "",
		},
		{
			Name: "[string] default",
			Value: (&Value{
				DataType: "string",
				Default:  "{{.First}} {{.Last}}",
			}).WithValueCache(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchInput(""),
					ui.RespondDefault(),
				)
			},
			Output: "Joey Ramone",
			Err:    "",
		},
		{
			Name: "[string] with options",
			Value: (&Value{
				DataType: "string",
				Default:  "foo",
				Options:  []any{"foo", "bar", "baz"},
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchSelect(""),
					ui.RespondString("baz"),
				)
			},
			Output: "baz",
			Err:    "",
		},
		{
			Name: "[string] with options default",
			Value: (&Value{
				DataType: "string",
				Default:  "foo",
				Options:  []any{"foo", "bar", "baz"},
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchSelect(""),
					ui.RespondDefault(),
				)
			},
			Output: "foo",
			Err:    "",
		},

		{
			Name: "[stringSlice]",
			Value: (&Value{
				DataType: "stringSlice",
				Default:  "",
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchInput(""),
					ui.RespondString("foo,bar"),
				)
			},
			Output: []string{"foo", "bar"},
			Err:    "",
		},
		{
			Name: "[stringSlice] with options",
			Value: (&Value{
				DataType: "stringSlice",
				Default:  "",
				Options:  []any{"foo", "bar", "baz"},
			}),
			Setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchMultiSelect(""),
					ui.RespondStringSlice([]string{"foo", "bar"}),
				)
			},
			Output: []string{"foo", "bar"},
			Err:    "",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			prompter := ui.NewUserInterface(ui.NewTestIOStreams()).WithStubbing()
			defer prompter.VerifyStubs(t)

			if test.Setup != nil {
				test.Setup(prompter)
			}

			err := test.Value.Prompt(prompter)
			if test.Err == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.Output, test.Value.Get(), "Get() should match output")
			} else {
				assert.ErrorContains(t, err, test.Err)
				assert.Equal(t, test.Value.Default, test.Value.Get(), "Get() should match default")
			}
		})
	}
}

func TestValue_ShouldPrompt(t *testing.T) {
	tests := []struct {
		Name         string
		Value        *Value
		Input        any
		ShouldPrompt bool
	}{
		{
			Name: "[never] never prompts - even if no value",
			Value: (&Value{
				DataType:     DataTypeString,
				PromptConfig: PromptConfigNever,
			}),
			ShouldPrompt: false,
		},
		{
			Name: "[always] always prompts - even if value present",
			Value: (&Value{
				DataType:     DataTypeString,
				PromptConfig: PromptConfigAlways,
			}),
			Input:        "bar",
			ShouldPrompt: true,
		},
		{
			Name: "[on-empty] only prompts when value is empty",
			Value: (&Value{
				DataType:     DataTypeString,
				PromptConfig: PromptConfigOnEmpty,
			}),
			Input:        "",
			ShouldPrompt: true,
		},
		{
			Name: "[on-empty] only prompts when value is not explicitly set by user",
			Value: (&Value{
				DataType:     DataTypeString,
				Default:      "foo",
				PromptConfig: PromptConfigOnUnset,
			}),
			ShouldPrompt: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.Input != nil {
				_ = test.Value.Set(test.Input.(string))
			}
			assert.Equal(t, test.ShouldPrompt, test.Value.ShouldPrompt())
		})
	}
}
