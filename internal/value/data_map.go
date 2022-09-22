package value

// DataMap is a shared data structure that each value uses when rendering.
// This allows for values to reference other values (via template expressions)
// in their defaults.
type DataMap map[string]any

// NewDataMap returns a new, empty DataMap.
func NewDataMap() DataMap {
	return make(DataMap)
}
