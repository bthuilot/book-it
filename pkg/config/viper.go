package config

import (
	"github.com/spf13/viper"
)

func Parse() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}
