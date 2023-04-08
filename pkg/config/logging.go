package config

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func InitLogger() {
	logrus.SetLevel(logrus.DebugLevel)
	//
	//	log, err := os.Create("/Users/bryce/output.log")
	//	if err != nil {
	//		logrus.Fatalf("unable to open log file: %s", err)
	//	}
	//	if os.Getenv("USE_STDOUT") != "1" {
	//		logrus.SetOutput(log)
	//	}
}

func LogHTTPRequest(r *http.Request) {
	// TODO(this)
}
