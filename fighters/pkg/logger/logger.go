package logger

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Init initializes global zap.Logger instance for the whole service
func Init(logLevel zapcore.Level, logFilePath string) error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	config.Level = zap.NewAtomicLevelAt(logLevel)

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
		logLevel,
	)

	l := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	if l == nil {
		return errors.New("error while creation zap logger")
	}

	namedLogger := l.Named("fighters-logger")

	zap.ReplaceGlobals(namedLogger)
	return nil
}
