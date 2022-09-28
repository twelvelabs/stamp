package core

import (
	"fmt"
	"os"

	"github.com/creasty/defaults" //cspell:disable-line
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Debug     bool           `mapstructure:"debug"`
	Defaults  map[string]any `mapstructure:"defaults"   default:"{}"`
	StorePath string         `mapstructure:"store_path" default:"~/.stamp/packages"`
}

func NewConfig(in string) *Config {
	if in != "" {
		viper.SetConfigFile(in)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.stamp/")
		viper.SetConfigName(".stamp")
		viper.SetConfigType("yaml")
	}
	viper.SetEnvPrefix("stamp")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	config := &Config{}

	if err := defaults.Set(config); err != nil {
		cobra.CheckErr(err)
	}

	err := viper.Unmarshal(config)
	if err != nil {
		cobra.CheckErr(err)
	}

	if config.Debug {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		fmt.Fprintln(os.Stderr, "Store path:", config.StorePath)
	}

	return config
}
