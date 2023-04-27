package value

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/ui"
)

func TestNewValueSet(t *testing.T) {
	vs := NewValueSet()
	assert.NotNil(t, vs)
	assert.IsType(t, &ValueSet{}, vs)
}

func TestValueSet_AddAndGetValue(t *testing.T) {
	vs := NewValueSet()

	// Unknown keys return nil.
	assert.Nil(t, vs.Value("unknown"))

	// Adding nil is a noop.
	vs.Add(nil)
	assert.Equal(t, 0, vs.Len())

	// Calls to Add() are chainable.
	value1 := &Value{Key: "value1"}
	value2 := &Value{Key: "value2"}
	vs.Add(value1).Add(value2)
	assert.Equal(t, 2, vs.Len())
	assert.Equal(t, value1, vs.Value("value1"))
	assert.Equal(t, value2, vs.Value("value2"))

	// Duplicate values overwrite existing ones w/ the same key.
	dupe1 := &Value{Key: "value1", Default: "something"}
	vs.Add(dupe1)
	assert.Equal(t, 2, vs.Len())
	assert.Equal(t, dupe1, vs.Value("value1"))

	// Values can be prepended
	value0 := &Value{Key: "value0"}
	vs.Prepend(value0)
	assert.Equal(t, 3, vs.Len())

	// And `All()` returns values in the correct insertion order.
	values := vs.All()
	assert.Equal(t, value0, values[0]) // prepended.
	assert.Equal(t, dupe1, values[1])  // original insertion order despite overwrite.
	assert.Equal(t, value2, values[2])
}

func TestValueSet_ValuesMethods(t *testing.T) {
	none := []*Value{}
	args := []*Value{
		{Key: "arg1", InputMode: InputModeArg},
		{Key: "arg2", InputMode: InputModeArg},
		{Key: "arg3", InputMode: InputModeArg},
	}
	flags := []*Value{
		{Key: "flag1", InputMode: InputModeFlag},
		{Key: "flag2", InputMode: InputModeFlag},
		{Key: "flag3", InputMode: InputModeFlag},
	}
	all := append([]*Value{}, args...)
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

func TestValueSet_CacheMethods(t *testing.T) {
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

func TestValueSet_AddArgs(t *testing.T) {
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
					Key:       "Arg1",
					Default:   "Arg1",
					DataType:  DataTypeString,
					InputMode: InputModeArg,
				},
				{
					Key:       "Flag1",
					Default:   "Flag1",
					DataType:  DataTypeString,
					InputMode: InputModeFlag,
				},
				{
					Key:       "Arg2",
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
					Key:       "Arg1",
					Default:   "Arg1",
					DataType:  DataTypeString,
					InputMode: InputModeArg,
				},
				{
					Key:       "Flag1",
					Default:   "Flag1",
					DataType:  DataTypeString,
					InputMode: InputModeFlag,
				},
				{
					Key:       "Arg2",
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
					Key:       "Arg1",
					Default:   "Arg1",
					DataType:  DataTypeString,
					InputMode: InputModeArg,
				},
				{
					Key:       "Flag1",
					Default:   "Flag1",
					DataType:  DataTypeString,
					InputMode: InputModeFlag,
				},
				{
					Key:       "Arg2",
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
					Key:       "Arg1",
					Default:   "Arg1",
					DataType:  DataTypeString,
					InputMode: InputModeArg,
				},
				{
					Key:       "Arg2",
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
			Key:      "DstPath",
			DataType: DataTypeString,
			Default:  "~/src/untitled",
		}).
		Add(&Value{
			Key:      "ProjectSlug",
			DataType: DataTypeString,
			Default:  "{{ .DstPath | base }}", // depends on DstPath
		}).
		Add(&Value{
			Key:      "PackageName",
			DataType: DataTypeString,
			Default:  "{{ .ProjectSlug | underscore }}", // depends on ProjectSlug
		})

	assert.Equal(t, map[string]any{
		"DstPath":     "~/src/untitled",
		"ProjectSlug": "untitled",
		"PackageName": "untitled",
	}, vs.GetAll())

	// Setting DstPath should cause other keys to be re-evaluated
	_ = vs.Set("DstPath", "~/src/my-project")

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
	_ = vs.Set("Foo", "aaa")
	assert.Equal(t, "aaa", vs.Get("Foo"))
	assert.Equal(t, "aaa", vs.Cache().Get("Foo"))

	value := &Value{
		Key:      "Foo",
		DataType: DataTypeString,
		Default:  "bbb",
	}
	// Now it's added as a Value, so that should take precedence.
	vs.Add(value)
	assert.Equal(t, "bbb", value.Get())
	assert.Equal(t, "bbb", vs.Get("Foo"))
	assert.Equal(t, "bbb", vs.Cache().Get("Foo"))

	// And setting it, should set the underlying value (as well as the cache)
	_ = vs.Set("Foo", "ccc")
	assert.Equal(t, "ccc", value.Get())
	assert.Equal(t, "ccc", vs.Get("Foo"))
	assert.Equal(t, "ccc", vs.Cache().Get("Foo"))
}

func TestValueSet_GetAll(t *testing.T) {
	vs := NewValueSet().SetCache(DataMap{
		"NonValue": 123,
	})

	vs.Add(&Value{
		Key:      "ProjectName",
		DataType: DataTypeString,
		Default:  "Example",
	})
	vs.Add(&Value{
		Key:      "ProjectSlug",
		DataType: DataTypeString,
		Default:  "{{ .ProjectName | underscore }}",
	})

	_ = vs.Set("ProjectName", "My Proj")

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
					Key:             "v1",
					DataType:        DataTypeString,
					Default:         "foo",
					ValidationRules: "required",
				},
				{
					Key:             "v2",
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
					Key:             "v1",
					DataType:        DataTypeString,
					Default:         "foo",
					ValidationRules: "required",
				},
				{
					Key:             "v2",
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
	tests := []struct {
		name      string
		values    []*Value
		setup     func(p *ui.UserInterface)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "returns nil when all prompts succeed",
			values: []*Value{
				{
					Key:      "v1",
					DataType: DataTypeBool,
				},
				{
					Key:      "v2",
					DataType: DataTypeBool,
				},
			},
			setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchConfirm("V1"),
					ui.RespondBool(true),
				)
				p.RegisterStub(
					ui.MatchConfirm("V2"),
					ui.RespondBool(true),
				)
			},
			assertion: assert.NoError,
		},
		{
			name: "returns an error when any prompts fail",
			values: []*Value{
				{
					Key:      "v1",
					DataType: DataTypeBool,
				},
				{
					Key:      "v2",
					DataType: DataTypeBool,
				},
			},
			setup: func(p *ui.UserInterface) {
				p.RegisterStub(
					ui.MatchConfirm("V1"),
					ui.RespondBool(true),
				)
				p.RegisterStub(
					ui.MatchConfirm("V2"),
					ui.RespondError(errors.New("boom")),
				)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompter := ui.NewUserInterface(ui.NewTestIOStreams()).WithStubbing()
			defer prompter.VerifyStubs(t)
			if tt.setup != nil {
				tt.setup(prompter)
			}

			vs := NewValueSet()
			for _, v := range tt.values {
				vs.Add(v)
			}
			tt.assertion(t, vs.Prompt(prompter))
		})
	}
}
