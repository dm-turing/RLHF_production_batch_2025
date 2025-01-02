package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	logrus.WithField("token", "SECRET_TOKEN").Info("Processing request")
}
