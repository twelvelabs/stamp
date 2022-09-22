package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				"name": "foo",
			},
			Output: &Value{
				Name:         "foo",
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
			Err:    "Name is a required field",
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

func TestValueGetAndSet(t *testing.T) {
	tests := []struct {
		Name   string
		Value  *Value
		Input  string
		Output any
		Err    string
	}{
		{
			Name: "empty type",
			Value: (&Value{
				DataType: "",
			}),
			Input:  "something",
			Output: nil,
			Err:    "invalid data type",
		},
		{
			Name: "invalid type",
			Value: (&Value{
				DataType: "wat",
			}),
			Input:  "something",
			Output: nil,
			Err:    "invalid data type",
		},

		// BOOL
		{
			Name: "bool",
			Value: (&Value{
				DataType: "bool",
				Default:  true,
			}),
			Input:  "false",
			Output: false,
			Err:    "",
		},
		{
			Name: "bool default",
			Value: (&Value{
				DataType: "bool",
				Default:  "{{.WannaDance}}",
			}).WithDataMap(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Input:  "",
			Output: true,
			Err:    "",
		},
		{
			Name: "bool invalid",
			Value: (&Value{
				DataType: "bool",
				Default:  true,
			}),
			Input:  "wat",
			Output: true,
			Err:    "unable to cast to bool",
		},

		// INT
		{
			Name: "int",
			Value: (&Value{
				DataType: "int",
				Default:  10,
			}),
			Input:  "25",
			Output: 25,
			Err:    "",
		},
		{
			Name: "int default",
			Value: (&Value{
				DataType: "int",
				Default:  "{{ add .Year 1 }}",
			}).WithDataMap(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Input:  "",
			Output: 1978,
			Err:    "",
		},
		{
			Name: "int invalid",
			Value: (&Value{
				DataType: "int",
				Default:  10,
			}),
			Input:  "wat",
			Output: 10,
			Err:    "unable to cast",
		},

		// STRING
		{
			Name: "string",
			Value: (&Value{
				DataType: "string",
				Default:  "foo",
			}),
			Input:  "bar",
			Output: "bar",
			Err:    "",
		},
		{
			Name: "string default",
			Value: (&Value{
				DataType: "string",
				Default:  "{{.First}} {{.Last}}",
			}).WithDataMap(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Input:  "",
			Output: "Joey Ramone",
			Err:    "",
		},
		{
			Name: "string rendered",
			Value: (&Value{
				DataType: "string",
				Default:  "{{.First}} {{.Last}}",
			}).WithDataMap(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Input:  "Dee Dee {{.Last}}",
			Output: "Dee Dee Ramone",
			Err:    "",
		},
		{
			Name: "string transformed",
			Value: (&Value{
				DataType:       "string",
				Default:        "",
				TransformRules: "trim, kebabcase, uppercase",
			}),
			Input:  "     hello world     ",
			Output: "HELLO-WORLD",
			Err:    "",
		},
		{
			Name: "string default transformed",
			Value: (&Value{
				DataType:       "string",
				Default:        "{{.First}} {{.Last}}",
				TransformRules: "uppercase",
			}).WithDataMap(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Input:  "",
			Output: "JOEY RAMONE",
			Err:    "",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var err error

			value := test.Value
			if test.Input != "" {
				err = value.Set(test.Input)
			}
			result := value.Get()

			assert.Equal(t, test.Output, result)
			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}

func TestPrompt(t *testing.T) {
	tests := []struct {
		Name     string
		Value    *Value
		Prompter Prompter
		Output   interface{}
		Err      string
	}{
		{
			Name: "bool true",
			Value: (&Value{
				DataType: "bool",
				Default:  false,
			}),
			Prompter: &PrompterMock{
				ConfirmFunc: NewConfirmFunc(true, nil),
			},
			Output: true,
			Err:    "",
		},
		{
			Name: "bool default",
			Value: (&Value{
				DataType: "bool",
				Default:  "{{.WannaDance}}",
			}).WithDataMap(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Prompter: &PrompterMock{
				ConfirmFunc: NewNoopConfirmFunc(),
			},
			Output: true,
			Err:    "",
		},
		{
			Name: "int",
			Value: (&Value{
				DataType: "int",
				Default:  1,
			}),
			Prompter: &PrompterMock{
				InputFunc: NewInputFunc("12", nil),
			},
			Output: 12,
			Err:    "",
		},
		{
			Name: "int default",
			Value: (&Value{
				DataType: "int",
				Default:  "{{ add .Year 1 }}",
			}).WithDataMap(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Prompter: &PrompterMock{
				InputFunc: NewNoopInputFunc(),
			},
			Output: 1978,
			Err:    "",
		},
		{
			Name: "int invalid",
			Value: (&Value{
				DataType: "int",
				Default:  1,
			}),
			Prompter: &PrompterMock{
				InputFunc: NewInputFunc("not an int", nil),
			},
			Output: nil,
			Err:    "unable to cast",
		},
		{
			Name: "string",
			Value: (&Value{
				DataType: "string",
				Default:  "foo",
			}),
			Prompter: &PrompterMock{
				InputFunc: NewInputFunc("bar", nil),
			},
			Output: "bar",
			Err:    "",
		},
		{
			Name: "string default",
			Value: (&Value{
				DataType: "string",
				Default:  "{{.First}} {{.Last}}",
			}).WithDataMap(DataMap{
				"First":      "Joey",
				"Last":       "Ramone",
				"Year":       1977,
				"WannaDance": true,
			}),
			Prompter: &PrompterMock{
				InputFunc: NewNoopInputFunc(),
			},
			Output: "Joey Ramone",
			Err:    "",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			value := test.Value
			err := value.Prompt(test.Prompter)

			if test.Err == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.Output, value.Get(), "Get() should match output")
			} else {
				assert.ErrorContains(t, err, test.Err)
				assert.Equal(t, value.Default, value.Get(), "Get() should match default")
			}
		})
	}
}
