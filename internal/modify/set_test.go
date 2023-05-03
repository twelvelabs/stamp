package modify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSet(t *testing.T) {
	s := NewSet(1, 2, 3)
	assert.Equal(t, 3, s.Len())
}

func TestSet_Contains(t *testing.T) {
	s := NewSet(3, 2, 1)

	assert.Equal(t, true, s.Contains(3))
	assert.Equal(t, false, s.Contains(99))
}

func TestSet_Add(t *testing.T) {
	s := NewSet()
	s.Add(1)
	s.Add(2)
	s.Add(1)

	assert.Equal(t, 2, s.Len())
}

func TestSet_ToSlice(t *testing.T) {
	s := NewSet(3, 2, 1)

	assert.Equal(t, []any{3, 2, 1}, s.ToSlice())
}
