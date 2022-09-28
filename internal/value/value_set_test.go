package value

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValueSet(t *testing.T) {
	vs := NewValueSet()
	assert.NotNil(t, vs)
	assert.IsType(t, &ValueSet{}, vs)
}

func TestValueSetAddAndGetValue(t *testing.T) {
	tests := []struct {
		Name  string
		Value *Value
		Key   string
		Err   string
	}{
		{
			Name:  "value is found",
			Value: &Value{Name: "foo bar"},
			Key:   "FooBar",
			Err:   "",
		},
		{
			Name:  "value not found",
			Value: nil,
			Key:   "FooBar",
			Err:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			vs := NewValueSet()
			// AddValue() should be chainable,
			// and adding nil should be a noop
			vs.Add(nil).Add(test.Value)
			assert.Equal(t, test.Value, vs.Value(test.Key))
		})
	}
}

func TestValueSetValuesMethods(t *testing.T) {
	none := []*Value{}
	args := []*Value{
		{Name: "arg1", InputMode: InputModeArg},
		{Name: "arg2", InputMode: InputModeArg},
		{Name: "arg3", InputMode: InputModeArg},
	}
	flags := []*Value{
		{Name: "flag1", InputMode: InputModeFlag},
		{Name: "flag2", InputMode: InputModeFlag},
		{Name: "flag3", InputMode: InputModeFlag},
	}
	all := append(none, args...)
	all = append(all, flags...)

	assert.Len(t, none, 0)
	assert.Len(t, args, 3)
	assert.Len(t, flags, 3)
	assert.Len(t, all, 6)

	tests := []struct {
		Name  string
		Input []*Value
		All   []*Value
		Args  []*Value
		Flags []*Value
	}{
		{
			Name:  "empty set",
			Input: none,
			All:   none,
			Args:  none,
			Flags: none,
		},
		{
			Name:  "mixed set",
			Input: all,
			All:   all,
			Args:  args,
			Flags: flags,
		},
		{
			Name:  "args only",
			Input: args,
			All:   args,
			Args:  args,
			Flags: none,
		},
		{
			Name:  "flags only",
			Input: flags,
			All:   flags,
			Args:  none,
			Flags: flags,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			vs := NewValueSet()
			for _, v := range test.All {
				vs.Add(v)
			}
			assert.Equal(t, test.All, vs.All())
			assert.Equal(t, test.Args, vs.Args())
			assert.Equal(t, test.Flags, vs.Flags())
			args, flags := vs.Partition()
			assert.Equal(t, test.Args, args)
			assert.Equal(t, test.Flags, flags)
		})
	}
}

func TestValueSetCacheMethods(t *testing.T) {
	vs := NewValueSet()
	assert.Equal(t, DataMap{}, vs.Cache())

	vs.Cache().Set("string-var", "hi")
	vs.Cache().Set("int-var", 1234)

	assert.Equal(t, "hi", vs.Cache().Get("string-var"))
	assert.Equal(t, 1234, vs.Cache().Get("int-var"))

	assert.Equal(t, DataMap{
		"string-var": "hi",
		"int-var":    1234,
	}, vs.Cache())

	vs.SetCache(DataMap{})
	assert.Equal(t, DataMap{}, vs.Cache())
}

func TestValueSetAddArgs(t *testing.T) {
	tests := []struct {
		Name    string
		Values  []*Value
		Input   []string
		Output  []string
		DataMap DataMap
		Err     string
	}{
		{
			Name:    "no values no input",
			Values:  []*Value{},
			Input:   []string{},
			Output:  []string{},
			DataMap: NewDataMap(),
			Err:     "",
		},
		{
			Name: "values no input",
			Values: []*Value{
				{
					Name:      "Arg1",
					Default:   "Arg1",
					DataType:  DataTypeString,
					InputMode: InputModeArg,
				},
				{
					Name:      "Flag1",
					Default:   "Flag1",
					DataType:  DataTypeString,
					InputMode: InputModeFlag,
				},
				{
					Name:      "Arg2",
					Default:   "Arg2",
					DataType:  DataTypeString,
					InputMode: InputModeArg,
				},
			},
			Input:  []string{},
			Output: []string{},
			DataMap: DataMap{
				"Arg1":  "Arg1",
				"Arg2":  "Arg2",
				"Flag1": "Flag1",
			},
			Err: "",
		},
		{
			Name: "values with partial input",
			Values: []*Value{
				{
					Name:      "Arg1",
					Default:   "Arg1",
					DataType:  DataTypeString,
					InputMode: InputModeArg,
				},
				{
					Name:      "Flag1",
					Default:   "Flag1",
					DataType:  DataTypeString,
					InputMode: InputModeFlag,
				},
				{
					Name:      "Arg2",
					Default:   "Arg2",
					DataType:  DataTypeString,
					InputMode: InputModeArg,
				},
			},
			Input:  []string{"foo"},
			Output: []string{},
			DataMap: DataMap{
				"Arg1":  "foo",
				"Arg2":  "Arg2",
				"Flag1": "Flag1",
			},
			Err: "",
		},
		{
			Name: "values with excess input",
			Values: []*Value{
				{
					Name:      "Arg1",
					Default:   "Arg1",
					DataType:  DataTypeString,
					InputMode: InputModeArg,
				},
				{
					Name:      "Flag1",
					Default:   "Flag1",
					DataType:  DataTypeString,
					InputMode: InputModeFlag,
				},
				{
					Name:      "Arg2",
					Default:   "Arg2",
					DataType:  DataTypeString,
					InputMode: InputModeArg,
				},
			},
			Input:  []string{"foo", "bar", "baz"},
			Output: []string{"baz"},
			DataMap: DataMap{
				"Arg1":  "foo",
				"Arg2":  "bar",
				"Flag1": "Flag1",
			},
			Err: "",
		},
		{
			Name: "values with invalid input",
			Values: []*Value{
				{
					Name:      "Arg1",
					Default:   "Arg1",
					DataType:  DataTypeString,
					InputMode: InputModeArg,
				},
				{
					Name:      "Arg2",
					Default:   "Arg2",
					DataType:  DataTypeBool,
					InputMode: InputModeArg,
				},
			},
			Input:  []string{"foo", "bar"},
			Output: nil,
			DataMap: DataMap{
				"Arg1":  "foo",
				"Arg2":  "Arg2",
				"Flag1": "Flag1",
			},
			Err: "unable to cast to bool",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			vs := NewValueSet().SetCache(test.DataMap)
			for _, v := range test.Values {
				vs.Add(v)
			}
			remaining, err := vs.SetArgs(test.Input)

			assert.Equal(t, test.Output, remaining)
			assert.Equal(t, test.DataMap, vs.Cache())

			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}

func TestValueSet_CacheInvalidation(t *testing.T) {
	vs := NewValueSet().
		Add(&Value{
			Name:     "Dst Path",
			DataType: DataTypeString,
			Default:  "~/src/untitled",
		}).
		Add(&Value{
			Name:     "Project Slug",
			DataType: DataTypeString,
			Default:  "{{ .DstPath | base }}", // depends on DstPath
		}).
		Add(&Value{
			Name:     "Package Name",
			DataType: DataTypeString,
			Default:  "{{ .ProjectSlug | underscore }}", // depends on ProjectSlug
		})

	assert.Equal(t, map[string]any{
		"DstPath":     "~/src/untitled",
		"ProjectSlug": "untitled",
		"PackageName": "untitled",
	}, vs.GetAll())

	// Setting DstPath should cause other keys to be re-evaluated
	vs.Set("DstPath", "~/src/my-project")

	assert.Equal(t, map[string]any{
		"DstPath":     "~/src/my-project",
		"ProjectSlug": "my-project",
		"PackageName": "my_project",
	}, vs.GetAll())
}

func TestValueSet_GetAndSet(t *testing.T) {
	dm := DataMap{
		"NonValue": 123,
	}
	vs := NewValueSet().SetCache(dm)

	// Should only know about the cache
	assert.Equal(t, 123, vs.Get("NonValue"))
	assert.Equal(t, nil, vs.Get("Foo"))

	// "foo" isn't a Value, so setting it sets it in the cache
	vs.Set("Foo", "aaa")
	assert.Equal(t, "aaa", vs.Get("Foo"))
	assert.Equal(t, "aaa", vs.Cache().Get("Foo"))

	value := &Value{
		Name:     "Foo",
		DataType: DataTypeString,
		Default:  "bbb",
	}
	// Now it's added as a Value, so that should take precedence.
	vs.Add(value)
	assert.Equal(t, "bbb", value.Get())
	assert.Equal(t, "bbb", vs.Get("Foo"))
	assert.Equal(t, "bbb", vs.Cache().Get("Foo"))

	// And setting it, should set the underlying value (as well as the cache)
	vs.Set("Foo", "ccc")
	assert.Equal(t, "ccc", value.Get())
	assert.Equal(t, "ccc", vs.Get("Foo"))
	assert.Equal(t, "ccc", vs.Cache().Get("Foo"))
}

func TestValueSet_GetAll(t *testing.T) {
	vs := NewValueSet().SetCache(DataMap{
		"NonValue": 123,
	})

	vs.Add(&Value{
		Name:     "Project Name",
		DataType: DataTypeString,
		Default:  "Example",
	})
	vs.Add(&Value{
		Name:     "Project Slug",
		DataType: DataTypeString,
		Default:  "{{ .ProjectName | underscore }}",
	})

	vs.Set("ProjectName", "My Proj")

	assert.Equal(t, map[string]any{
		"NonValue":    123,
		"ProjectName": "My Proj",
		"ProjectSlug": "my_proj",
	}, vs.GetAll())
}

func TestValueSet_Validate(t *testing.T) {
	tests := []struct {
		name      string
		values    []*Value
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "returns nil when all values are valid",
			values: []*Value{
				{
					DataType:        DataTypeString,
					Default:         "foo",
					ValidationRules: "required",
				},
				{
					DataType:        DataTypeString,
					Default:         "bar",
					ValidationRules: "required",
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "returns an error when any values are invalid",
			values: []*Value{
				{
					DataType:        DataTypeString,
					Default:         "foo",
					ValidationRules: "required",
				},
				{
					DataType:        DataTypeString,
					Default:         "",
					ValidationRules: "required",
				},
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vs := NewValueSet()
			for _, v := range tt.values {
				vs.Add(v)
			}
			tt.assertion(t, vs.Validate())
		})
	}
}

func TestValueSet_Prompt(t *testing.T) {
	callCount := 0
	tests := []struct {
		name      string
		values    []*Value
		prompter  *PrompterMock
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "returns nil when all prompts succeed",
			values: []*Value{
				{
					DataType: DataTypeBool,
				},
				{
					DataType: DataTypeBool,
				},
			},
			prompter: &PrompterMock{
				ConfirmFunc: NewConfirmFunc(true, nil), // always succeeds
			},
			assertion: assert.NoError,
		},
		{
			name: "returns an error when any prompts fail",
			values: []*Value{
				{
					DataType: DataTypeBool,
				},
				{
					DataType: DataTypeBool,
				},
			},
			prompter: &PrompterMock{
				// fails on second prompt
				ConfirmFunc: func(prompt string, defaultValue bool, help, validationRules string) (bool, error) {
					callCount++
					if callCount >= 2 {
						return false, errors.New("boom")
					}
					return true, nil
				},
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vs := NewValueSet()
			for _, v := range tt.values {
				vs.Add(v)
			}
			tt.assertion(t, vs.Prompt(tt.prompter))
			assert.Equal(t, len(vs.All()), len(tt.prompter.ConfirmCalls()))
		})
	}
}
