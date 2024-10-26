package stamp

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/twelvelabs/termite/conf"

	"github.com/twelvelabs/stamp/internal/fsutil"
)

type Config struct {
	Debug     bool           `yaml:"debug"       env:"STAMP_DEBUG"`
	Defaults  map[string]any `yaml:"defaults"    default:"{}"`
	DryRun    bool           `yaml:"dry_run"     env:"STAMP_DRY_RUN"`
	StorePath string         `yaml:"store_path"  env:"STAMP_STORE_PATH"  default:"~/.stamp/packages"`
}

// NewDefaultConfig returns a new, default config.
func NewDefaultConfig() (*Config, error) {
	config := &Config{}
	return config, defaults.Set(config)
}

// NewConfig returns a new config for the file at path.
// If path is empty, uses one of:
//   - .stamp.yaml
//   - ~/.stamp/config.yaml
func NewConfig(path string) (*Config, error) {
	var err error

	if path == "" {
		if fsutil.PathExists(".stamp.yaml") {
			path = ".stamp.yaml"
		} else {
			path = os.ExpandEnv("$HOME/.stamp/config.yaml")
		}
	}

	config, _ := NewDefaultConfig()
	if fsutil.PathExists(path) {
		config, err = conf.NewLoader(config, path).Load()
		if err != nil {
			return nil, fmt.Errorf("config load: %w", err)
		}
	}

	if config.Debug {
		fmt.Fprintln(os.Stderr, "Using config file:", path)
		fmt.Fprintln(os.Stderr, "Store path:", config.StorePath)
	}

	return config, nil
}
