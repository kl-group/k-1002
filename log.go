package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	GlobalLog = logrus.New()
)

func entryLog() {
	GlobalLog.SetFormatter(&logrus.TextFormatter{})
	GlobalLog.SetOutput(os.Stdout)
	GlobalLog.SetLevel(logrus.InfoLevel)
}
