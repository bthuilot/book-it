package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func InitLogger() {
	logrus.SetLevel(logrus.DebugLevel)
	if viper.IsSet("log_file") {
		logFile := viper.GetString("log_file")
		log, err := os.Create(logFile)
		if err != nil {
			logrus.Fatalf("unable to open log file: %s", err)
		}
		logrus.SetOutput(log)
	}
}
