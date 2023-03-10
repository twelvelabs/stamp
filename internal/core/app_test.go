package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	app, err := NewApp()
	assert.IsType(t, &App{}, app)
	assert.NotNil(t, app)
	assert.NoError(t, err)
}

func TestNewTestApp(t *testing.T) {
	app := NewTestApp()
	assert.IsType(t, &App{}, app)
	assert.NotNil(t, app)

	// ensure store is pointing to stamp/testdata/generators.
	generator, err := app.Store.Load("file")
	assert.NotNil(t, generator)
	assert.NoError(t, err)
}

func TestAppForContext(t *testing.T) {
	app := NewTestApp()
	ctx := app.Context()
	assert.Equal(t, app, AppForContext(ctx))
}
