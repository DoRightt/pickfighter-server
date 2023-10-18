package main

import (
	"errors"
	"projects/fb-server/logger"

	"go.uber.org/zap"
)

func main() {
	logLevel := "info"
	logFilePath := "logger/logs/log.json"

	if err := logger.Initialize(logLevel, logFilePath); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	mylog := logger.Get()

	mylog.Info("This is an info log message")
	mylog.Error("This is an error log message", zap.Error(errors.New("test")))
}
