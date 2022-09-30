package iostreams

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystem(t *testing.T) {
	ios := System()
	assert.NotNil(t, ios.In)
	assert.NotNil(t, ios.Out)
	assert.NotNil(t, ios.Err)
}

func TestIOStreamImplementations(t *testing.T) {
	s := &systemIOStream{File: os.Stdin}
	_ = s.String()

	m := &mockIOStream{Buffer: &bytes.Buffer{}, fd: 1}
	assert.Equal(t, 1, int(m.Fd()))
}
