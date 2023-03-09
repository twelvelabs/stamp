package core

import (
	"fmt"

	"github.com/twelvelabs/termite/ui"

	"github.com/twelvelabs/stamp/internal/fsutil"
	"github.com/twelvelabs/stamp/internal/gen"
	"github.com/twelvelabs/stamp/internal/prompt"
	"github.com/twelvelabs/stamp/internal/value"
)

type App struct {
	Config   *Config
	IO       *ui.IOStreams
	Prompter value.Prompter
	Store    *gen.Store
}

func NewApp() (*App, error) {
	config := NewConfig("")
	ios := ui.NewIOStreams()
	prompter := prompt.NewSurveyPrompter(ios.In, ios.Out, ios.Err)

	storePath, err := fsutil.NormalizePath(config.StorePath)
	if err != nil {
		return nil, fmt.Errorf("startup error: %w", err)
	}
	err = fsutil.EnsureDirWritable(storePath)
	if err != nil {
		return nil, fmt.Errorf("startup error: %w", err)
	}
	store := gen.NewStore(storePath)

	app := &App{
		Config:   config,
		IO:       ios,
		Prompter: prompter,
		Store:    store,
	}

	return app, nil
}
