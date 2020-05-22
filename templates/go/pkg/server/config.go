package server

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config viper.Config

func (c Config) GetString(arg string, defaults ...string) string {
	v := c.(viper.Config).GetString(arg)
	if v {
		return v
	}

	if len(defaults) == 1 {
		return defaults[0]
	}

	fmt.Panicf("Missing required configuration: %s", arg)
}

func NewConfig(file string) *Config {
	v := viper.New()
	v.SetConfigFile(file)
	return v
}
