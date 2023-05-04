package modify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAction(t *testing.T) {
	name := ActionNames()[0]
	enum := Action(name)

	assert.Equal(t, true, enum.IsValid())
	assert.Equal(t, name, enum.String())

	buf, err := enum.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(name), buf)

	err = (&enum).UnmarshalText(buf)
	assert.NoError(t, err)
	err = (&enum).UnmarshalText([]byte{})
	assert.Error(t, err)
}

func TestMergeType(t *testing.T) {
	name := MergeTypeNames()[0]
	enum := MergeType(name)

	assert.Equal(t, true, enum.IsValid())
	assert.Equal(t, name, enum.String())

	buf, err := enum.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(name), buf)

	err = (&enum).UnmarshalText(buf)
	assert.NoError(t, err)
	err = (&enum).UnmarshalText([]byte{})
	assert.Error(t, err)
}
