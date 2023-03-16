package stamp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// dummy tests that exercise all the generated enum methods
// and ensure that coverage numbers don't take a hit.

func TestConflictConfig(t *testing.T) {
	name := ConflictConfigNames()[0]
	enum := ConflictConfig(name)

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

func TestMissingConfig(t *testing.T) {
	name := MissingConfigNames()[0]
	enum := MissingConfig(name)

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

func TestUpdateAction(t *testing.T) {
	name := UpdateActionNames()[0]
	enum := UpdateAction(name)

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
