package core

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/twelvelabs/termite/ui"

	"github.com/twelvelabs/stamp/internal/fsutil"
	"github.com/twelvelabs/stamp/internal/gen"
	"github.com/twelvelabs/stamp/internal/prompt"
	"github.com/twelvelabs/stamp/internal/value"
)

type ctxKey string

var (
	ctxKeyApp ctxKey = "github.com/twelvelabs/stamp/internal/core.App"
)

type App struct {
	Config   *Config
	IO       *ui.IOStreams
	UI       *ui.UserInterface
	Prompter value.Prompter
	Store    *gen.Store

	ctx context.Context //nolint: containedctx
}

// Context returns the root [context.Context] for the app.
func (a *App) Context() context.Context {
	if a.ctx == nil {
		a.ctx = context.WithValue(context.Background(), ctxKeyApp, a)
	}
	return a.ctx
}

// AppForContext returns the app singleton stored in the given context.
func AppForContext(ctx context.Context) *App {
	return ctx.Value(ctxKeyApp).(*App)
}

func NewApp() (*App, error) {
	config, err := NewConfig("")
	if err != nil {
		return nil, err
	}

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
		UI:       ui.NewUserInterface(ios),
		Prompter: prompter,
		Store:    store,
	}

	return app, nil
}

func NewTestApp() *App {
	config, _ := NewDefaultConfig()
	ios := ui.NewTestIOStreams()
	prompter := &value.PrompterMock{}

	storePath, _ := filepath.Abs(filepath.Join("..", "gen", "testdata", "generators"))
	store := gen.NewStore(storePath)

	app := &App{
		Config:   config,
		IO:       ios,
		UI:       ui.NewUserInterface(ios),
		Prompter: prompter,
		Store:    store,
	}

	return app
}
