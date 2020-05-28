package server

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func (c Config) GetString(arg string, defaults ...string) string {
	v := c.Viper.GetString(arg)
	if v != "" {
		return v
	}

	if len(defaults) == 1 {
		return defaults[0]
	}

	panic(fmt.Sprintf("Missing required configuration: %s", arg))
}

func NewConfig(file string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(file)
	err := v.ReadInConfig()
	return &Config{v}, err
}
