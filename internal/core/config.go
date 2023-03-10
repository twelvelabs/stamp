package core

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

type Config struct {
	Debug     bool           `mapstructure:"debug"`
	Defaults  map[string]any `mapstructure:"defaults"   default:"{}"`
	StorePath string         `mapstructure:"store_path" default:"~/.stamp/packages"`
}

// NewDefaultConfig returns a new, default config.
func NewDefaultConfig() (*Config, error) {
	config := &Config{}
	return config, defaults.Set(config)
}

// NewConfig returns a new config for the file at path.
// If path is empty, uses one of:
//   - .stamp.yaml
//   - ~/.stamp/.stamp.yaml
func NewConfig(path string) (*Config, error) {
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.stamp/")
		viper.SetConfigName(".stamp")
		viper.SetConfigType("yaml")
	}
	viper.SetEnvPrefix("stamp")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	config, _ := NewDefaultConfig()
	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	if config.Debug {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		fmt.Fprintln(os.Stderr, "Store path:", config.StorePath)
	}

	return config, nil
}
