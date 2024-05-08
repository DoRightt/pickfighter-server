package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func Initialize(flag int) error {
	config := zap.NewDevelopmentConfig()

	config.Encoding = "json"
	config.OutputPaths = []string{"./logs/log.json"}
	config.EncoderConfig = zapcore.EncoderConfig{
		LevelKey:     "level",
		TimeKey:      "timestamp",
		CallerKey:    "caller",
		MessageKey:   "message",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	logLevel := int(zapcore.DebugLevel)

	config.Level = zap.NewAtomicLevelAt(zapcore.Level(logLevel))

	file, err := os.OpenFile("./logs/log.json", os.O_WRONLY|os.O_CREATE|flag, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
		return err
	}
	defer file.Close()

	l, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to run zap logger: %s", err)
		return err
	}

	logger = l.Sugar()

	return nil
}

func Get() *zap.SugaredLogger {
	return logger
}
