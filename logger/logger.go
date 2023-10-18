package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

// Initialize initializes logger
func Initialize(logLevel string, logFilePath string) error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return err
	}

	config.Level = zap.NewAtomicLevelAt(level)

	// For log rotation
	logRotate := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     28,
		Compress:   true,
	}

	writeSyncer := zapcore.AddSync(logRotate)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		writeSyncer,
		level,
	)

	logger = zap.New(core, zap.AddCaller())

	return nil
}

// GetLogger returns the logger instance to use in other packages
func Get() *zap.Logger {
	return logger
}
