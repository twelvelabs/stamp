package value

import (
	"github.com/twelvelabs/termite/ui"
)

// ValueSet is a unique set of Values (identified by Value.Key).
type ValueSet struct { //nolint:revive
	keys    []string
	values  map[string]*Value
	dataMap DataMap
}

// NewValueSet returns a new ValueSet.
func NewValueSet() *ValueSet {
	return &ValueSet{
		values:  map[string]*Value{},
		dataMap: NewDataMap(),
	}
}

// Len returns the number of values in the set.
func (vs *ValueSet) Len() int {
	return len(vs.values)
}

// All returns all values in the set.
func (vs *ValueSet) All() []*Value {
	values := []*Value{}
	for _, key := range vs.keys {
		values = append(values, vs.values[key])
	}
	return values
}

// Args returns only the arg values.
func (vs *ValueSet) Args() []*Value {
	args, _ := vs.Partition()
	return args
}

// Flags returns only the flag values.
func (vs *ValueSet) Flags() []*Value {
	_, flags := vs.Partition()
	return flags
}

// Partition partitions values into args and flags.
func (vs *ValueSet) Partition() ([]*Value, []*Value) {
	args := []*Value{}
	flags := []*Value{}
	for _, v := range vs.All() {
		if v.IsArg() {
			args = append(args, v)
		}
		if v.IsFlag() {
			flags = append(flags, v)
		}
	}
	return args, flags
}

// Add appends a value to the set.
// Values are identified by Value.Key and duplicates are overwritten.
func (vs *ValueSet) Add(value *Value) *ValueSet {
	return vs.add(value, false)
}

// Prepend prepends a value to the set.
// Values are identified by Value.Key and duplicates are overwritten.
func (vs *ValueSet) Prepend(value *Value) *ValueSet {
	return vs.add(value, true)
}

func (vs *ValueSet) add(value *Value, prepend bool) *ValueSet {
	if value != nil {
		if _, found := vs.values[value.Key]; !found {
			if prepend {
				vs.keys = append([]string{value.Key}, vs.keys...)
			} else {
				vs.keys = append(vs.keys, value.Key)
			}
		}
		vs.values[value.Key] = value.WithValueSet(vs)
		vs.Cache().Set(value.Key, value.Get())
	}
	return vs
}

// Value returns the value for key.
func (vs *ValueSet) Value(key string) *Value {
	for _, val := range vs.All() {
		if val.Key == key {
			return val
		}
	}
	return nil
}

// Cache returns the materialized data map.
func (vs *ValueSet) Cache() DataMap {
	return vs.dataMap
}

// SetCache replaces the materialized data map.
func (vs *ValueSet) SetCache(dataMap DataMap) *ValueSet {
	vs.dataMap = dataMap
	return vs
}

// Get returns the value data for key.
// If key is not found, then returns the data for key
// from the cache.
func (vs *ValueSet) Get(key string) any {
	if v := vs.Value(key); v != nil {
		return v.Get()
	}
	return vs.Cache().Get(key)
}

// GetAll returns all data in the set.
func (vs *ValueSet) GetAll() map[string]any {
	data := map[string]any{}

	// Start w/ the cached values (so we get non-value data)
	for k, v := range vs.Cache() {
		data[k] = v
	}

	// Then do an explicit get on each value.
	// Doing this so because some values may have opted out of prompting,
	// and if so then may have default values that need to be rendered w/
	// the latest set of data.
	for _, val := range vs.All() {
		data[val.Key] = val.Get()
	}

	return data
}

// Set sets the value data for key.
// If key is not found, then sets the data in the cache
// so that it can be used by other values
// (see SrcPath and DstPath in stamp.Generator).
func (vs *ValueSet) Set(key string, value any) error {
	if v := vs.Value(key); v != nil {
		return v.set(value)
	}
	vs.Cache().Set(key, value)
	return nil
}

// SetArgs attempts to set all positional values with args.
// If len(args) > len(ValueSet.Args()), then the remaining
// items in args are returned.
func (vs *ValueSet) SetArgs(args []string) ([]string, error) {
	for _, val := range vs.Args() {
		if len(args) == 0 {
			continue
		}

		// shift off the first arg
		var a string
		a, args = args[0], args[1:]

		// attempt to set
		err := val.Set(a)
		if err != nil {
			return nil, err
		}
	}
	return args, nil
}

// Prompt calls Value.Prompt() for each value in the set.
// Returns the first error received.
func (vs *ValueSet) Prompt(prompter ui.Prompter) error {
	_ = vs.GetAll() // workaround for cache invalidation issue
	for _, val := range vs.All() {
		if err := val.Prompt(prompter); err != nil {
			return err
		}
	}
	return nil
}

// Validate calls Value.Validate() for each value in the set.
// Returns the first error received.
func (vs *ValueSet) Validate() error {
	for _, val := range vs.All() {
		if err := val.Validate(); err != nil {
			return err
		}
	}
	return nil
}
