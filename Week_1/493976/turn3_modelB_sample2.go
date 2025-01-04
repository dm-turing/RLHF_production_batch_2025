package main

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func main() {
	// Set up Logger
	log := logrus.New()
	// Disable color output (enable if supported by the terminal)
	log.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	writer := getLogWriter()
	hooks := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		},
		&logrus.JSONFormatter{},
	)
	log.AddHook(hooks)
	log.SetLevel(logrus.InfoLevel)

	// Write logs asynchronously
	for i := 0; i < 100000; i++ {
		log.Info("Hello, World!")
	}
}
