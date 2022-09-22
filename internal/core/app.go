package core

import (
	"fmt"

	"github.com/twelvelabs/stamp/internal/fsutil"
	"github.com/twelvelabs/stamp/internal/gen"
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/prompt"
	"github.com/twelvelabs/stamp/internal/value"
)

type App struct {
	Config   *Config
	FsUtil   *fsutil.FsUtil
	IO       *iostreams.IOStreams
	Prompter value.Prompter
	Store    *gen.Store
}

func NewApp() (*App, error) {
	config := NewConfig("")
	fsUtil := fsutil.NewFsUtil()
	ios := iostreams.System()
	prompter := prompt.NewSurveyPrompter(ios.In, ios.Out, ios.Err)

	storePath, err := fsUtil.NormalizePath(config.StorePath)
	if err != nil {
		return nil, fmt.Errorf("startup error: %w", err)
	}
	err = fsUtil.EnsureDirWritable(storePath)
	if err != nil {
		return nil, fmt.Errorf("startup error: %w", err)
	}

	store := gen.NewStore(storePath)
	app := &App{
		Config:   config,
		FsUtil:   fsUtil,
		IO:       ios,
		Prompter: prompter,
		Store:    store,
	}

	return app, nil
}
