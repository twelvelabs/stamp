package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataMapGetAndSet(t *testing.T) {
	dm := NewDataMap()
	assert.Equal(t, DataMap{}, dm)

	dm.Set("string-var", "hi")
	dm.Set("int-var", 1234)
	dm.Set("bool-var", true)

	// Get should return data or nil
	assert.Equal(t, "hi", dm.Get("string-var"))
	assert.Equal(t, 1234, dm.Get("int-var"))
	assert.Equal(t, true, dm.Get("bool-var"))
	assert.Equal(t, nil, dm.Get("not-found"))

	assert.Equal(t, DataMap{
		"string-var": "hi",
		"int-var":    1234,
		"bool-var":   true,
	}, dm)
}
