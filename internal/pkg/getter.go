package pkg

import (
	"context"
	"fmt"
	"net/url"
	"os"

	getter "github.com/hashicorp/go-getter"
	cp "github.com/otiai10/copy" //cspell:disable-line
)

// Getter is a function that copies a package from `src` to `dst`.
type Getter func(ctx context.Context, src, dst string) error

// DefaultGetter uses hashicorp/go-getter to copy packages.
func DefaultGetter(ctx context.Context, src, dst string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to resolve working directory: %w", err)
	}

	// Override the default `file://` getter.
	getters := getter.Getters
	getters["file"] = NewCopyFileGetter()

	client := getter.Client{
		Ctx:             ctx,
		Src:             src,
		Dst:             dst,
		Pwd:             pwd,
		DisableSymlinks: true,
		Getters:         getters,
		Mode:            getter.ClientModeDir,
	}

	if err := client.Get(); err != nil {
		return fmt.Errorf("unable to install: %w", err)
	}

	return nil
}

// MockGetter delegates to the supplied handler function.
type MockGetter struct {
	Called  bool
	Ctx     context.Context //nolint:containedctx
	Src     string
	Dst     string
	handler Getter
}

// NewMockGetter returns a new MockGetter struct.
func NewMockGetter(handler Getter) *MockGetter {
	return &MockGetter{handler: handler}
}

// Get implements the Getter type and is what should be passed to the store.
// It logs the arguments and delegates to the handler.
func (m *MockGetter) Get(ctx context.Context, src, dst string) error {
	m.Called = true
	m.Ctx = ctx
	m.Src = src
	m.Dst = dst
	return m.handler(ctx, src, dst)
}

// CopyFileGetter is a wrapper around the default [getter.FileGetter]
// that copies local files and directories rather than symlinking.
// The default getter has a `Copy` flag that is _supposed_ to allow for this,
// but it's only respected when getting individual files (not directories).
type CopyFileGetter struct {
	getter.FileGetter
}

func NewCopyFileGetter() *CopyFileGetter {
	return &CopyFileGetter{
		FileGetter: getter.FileGetter{Copy: true},
	}
}

func (g *CopyFileGetter) Get(dst string, u *url.URL) error {
	src := u.Path
	if u.RawPath != "" {
		src = u.RawPath
	}

	// The source path must exist and be a directory to be usable.
	if fi, err := os.Stat(src); err != nil {
		return fmt.Errorf("source path error: %w", err)
	} else if !fi.IsDir() {
		return fmt.Errorf("source path must be a directory")
	}

	opt := cp.Options{
		OnDirExists: func(src, dst string) cp.DirExistsAction {
			return cp.Replace
		},
		OnSymlink: func(src string) cp.SymlinkAction {
			return cp.Skip
		},
	}

	err := cp.Copy(src, dst, opt)
	if err != nil {
		return err
	}

	return nil
}
