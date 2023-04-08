package config

import (
	"github.com/spf13/viper"
)

func ParseEnv() {
	viper.SetEnvPrefix("resy")
	viper.AutomaticEnv()
}
