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
			vs.AddValue(nil).AddValue(test.Value)
			assert.Equal(t, test.Value, vs.GetValue(test.Key))
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
				vs.AddValue(v)
			}
			assert.Equal(t, test.All, vs.GetValues())
			assert.Equal(t, test.Args, vs.GetArgValues())
			assert.Equal(t, test.Flags, vs.GetFlagValues())
			args, flags := vs.PartitionValues()
			assert.Equal(t, test.Args, args)
			assert.Equal(t, test.Flags, flags)
		})
	}
}

func TestValueSetDataMapMethods(t *testing.T) {
	vs := NewValueSet()

	// Data map should be empty
	assert.Equal(t, DataMap{}, vs.GetDataMap())
	// Keys should show as unset
	assert.Equal(t, false, vs.HasData("string-var"))

	// SetData should be chainable
	vs.SetData("string-var", "hi").
		SetData("int-var", 1234).
		SetData("bool-var", true)

	// GetData should return data or nil
	assert.Equal(t, "hi", vs.GetData("string-var"))
	assert.Equal(t, 1234, vs.GetData("int-var"))
	assert.Equal(t, true, vs.GetData("bool-var"))
	assert.Equal(t, nil, vs.GetData("not-found"))

	// Keys should show as set
	assert.Equal(t, true, vs.HasData("string-var"))

	// Data map should be populated
	assert.Equal(t, DataMap{
		"string-var": "hi",
		"int-var":    1234,
		"bool-var":   true,
	}, vs.GetDataMap())

	// Data map should be reset
	vs.SetDataMap(DataMap{})
	assert.Equal(t, DataMap{}, vs.GetDataMap())

	// Keys should show as unset again
	assert.Equal(t, false, vs.HasData("string-var"))
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
			vs := NewValueSet().SetDataMap(test.DataMap)
			for _, v := range test.Values {
				vs.AddValue(v)
			}
			remaining, err := vs.SetArgs(test.Input)

			assert.Equal(t, test.Output, remaining)
			assert.Equal(t, test.DataMap, vs.GetDataMap())

			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
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
				vs.AddValue(v)
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
				vs.AddValue(v)
			}
			tt.assertion(t, vs.Prompt(tt.prompter))
			assert.Equal(t, len(vs.GetValues()), len(tt.prompter.ConfirmCalls()))
		})
	}
}
