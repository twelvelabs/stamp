package value

type ValueSet struct {
	values  []*Value
	dataMap DataMap
}

// NewValueSet returns a new ValueSet.
func NewValueSet() *ValueSet {
	return &ValueSet{
		values:  []*Value{},
		dataMap: NewDataMap(),
	}
}

// GetValues returns all values in the set.
func (vs *ValueSet) GetValues() []*Value {
	return vs.values
}

// GetArgValues returns only the arg values.
func (vs *ValueSet) GetArgValues() []*Value {
	args, _ := vs.PartitionValues()
	return args
}

// GetFlagValues returns only the flag values.
func (vs *ValueSet) GetFlagValues() []*Value {
	_, flags := vs.PartitionValues()
	return flags
}

// PartitionValues partitions values into args and flags.
func (vs *ValueSet) PartitionValues() ([]*Value, []*Value) {
	args := []*Value{}
	flags := []*Value{}
	for _, v := range vs.GetValues() {
		if v.IsArg() {
			args = append(args, v)
		}
		if v.IsFlag() {
			flags = append(flags, v)
		}
	}
	return args, flags
}

// AddValue adds val to the set.
func (vs *ValueSet) AddValue(val *Value) *ValueSet {
	if val != nil {
		vs.values = append(vs.values, val.WithValueSet(vs))
		vs.SetData(val.Key(), val.Get())
	}
	return vs
}

// GetValue returns the value for key.
func (vs *ValueSet) GetValue(key string) *Value {
	for _, val := range vs.GetValues() {
		if val.Key() == key {
			return val
		}
	}
	return nil
}

// GetDataMap returns the materialized data map.
func (vs *ValueSet) GetDataMap() DataMap {
	return vs.dataMap
}

// SetDataMap replaces the materialized data map.
func (vs *ValueSet) SetDataMap(dataMap DataMap) *ValueSet {
	vs.dataMap = dataMap
	return vs
}

// HasData returns true if key exists in the materialized data map.
func (vs *ValueSet) HasData(key string) bool {
	_, ok := vs.dataMap[key]
	return ok
}

// GetData returns the materialized data for key.
func (vs *ValueSet) GetData(key string) any {
	return vs.dataMap[key]
}

// SetData sets the materialized data for key.
func (vs *ValueSet) SetData(key string, data any) *ValueSet {
	vs.dataMap[key] = data
	return vs
}

func (vs *ValueSet) SetArgs(args []string) ([]string, error) {
	for _, val := range vs.GetArgValues() {
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

func (vs *ValueSet) Prompt(prompter Prompter) error {
	for _, val := range vs.GetValues() {
		if err := val.Prompt(prompter); err != nil {
			return err
		}
	}
	return nil
}

func (vs *ValueSet) Validate() error {
	for _, val := range vs.GetValues() {
		if err := val.Validate(); err != nil {
			return err
		}
	}
	return nil
}
