package stamp

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/twelvelabs/termite/ui"

	"github.com/twelvelabs/stamp/internal/fsutil"
)

type ctxKey string

var (
	ctxKeyApp ctxKey = "github.com/twelvelabs/stamp/internal/stamp.App"
)

type App struct {
	Config *Config
	IO     *ui.IOStreams
	UI     *ui.UserInterface
	Store  *Store
	Meta   *AppMeta

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

func NewApp(meta *AppMeta) (*App, error) {
	config, err := NewConfig("")
	if err != nil {
		return nil, err
	}

	ios := ui.NewIOStreams()

	storePath, err := fsutil.NormalizePath(config.StorePath)
	if err != nil {
		return nil, fmt.Errorf("startup error: %w", err)
	}
	err = fsutil.EnsureDirWritable(storePath)
	if err != nil {
		return nil, fmt.Errorf("startup error: %w", err)
	}
	store := NewStore(storePath)
	err = store.Init()
	if err != nil {
		return nil, err
	}

	app := &App{
		Config: config,
		IO:     ios,
		UI:     ui.NewUserInterface(ios),
		Store:  store,
		Meta:   meta,
	}

	return app, nil
}

func NewTestApp() *App {
	meta := NewAppMeta("test", "", "0")
	config, _ := NewDefaultConfig()
	ios := ui.NewTestIOStreams()

	storePath, _ := filepath.Abs(filepath.Join("testdata", "generators"))
	store := NewStore(storePath)

	app := &App{
		Config: config,
		IO:     ios,
		UI:     ui.NewUserInterface(ios).WithStubbing(),
		Store:  store,
		Meta:   meta,
	}

	return app
}

func NewAppMeta(version, commit, date string) *AppMeta {
	buildTime, _ := time.Parse(time.RFC3339, date)

	meta := &AppMeta{
		BuildCommit: commit,
		BuildTime:   buildTime,
		Version:     version,
		GOOS:        runtime.GOOS,
		GOARCH:      runtime.GOARCH,
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		meta.BuildGoVersion = info.GoVersion
		meta.BuildVersion = info.Main.Version
		meta.BuildChecksum = info.Main.Sum
	}

	return meta
}

type AppMeta struct {
	BuildCommit    string
	BuildTime      time.Time
	BuildGoVersion string
	BuildVersion   string
	BuildChecksum  string
	Version        string
	GOOS           string
	GOARCH         string
}
