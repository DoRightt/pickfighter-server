package logger

import (
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type FbLogger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Named(name string) *zap.SugaredLogger
}

var logger *zap.Logger

// Initialize initializes logger
func Initialize(logLevel zapcore.Level, logFilePath string) error {
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

	logger = zap.New(core, zap.AddCaller())

	return nil
}

// Get returns the logger instance
func Get() *zap.Logger {
	return logger
}

// Get returns the sugared logger instance
func GetSugared() *zap.SugaredLogger {
	return logger.Sugar()
}

// New creates a new instance of Zap logger with configuration based on application settings.
// It checks if the 'log_json' configuration option is set to true, and if so, configures the logger
// to output JSON format logs to the specified file path. Otherwise, it configures the logger for
// colored console output. The log level is determined by the 'log_level' configuration option, with
// a default of Debug level if the configuration is not set or invalid.
func New() *zap.Logger {
	config := zap.NewDevelopmentConfig()

	if logJson := viper.GetBool("log_json"); logJson {
		config.Encoding = "json"
		config.OutputPaths = []string{"logs/log.json"}
		config.EncoderConfig = zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "timestamp",
			CallerKey:    "caller",
			MessageKey:   "message",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		}
	} else {
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logLevel, ok := viper.Get("log_level").(int)
	if !ok {
		logLevel = int(zapcore.DebugLevel)
	}

	config.Level = zap.NewAtomicLevelAt(zapcore.Level(logLevel))
	l, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to run zap logger: %s", err)
	}

	return l
}

// NewSugared returns SugaredLogger version
func NewSugared() *zap.SugaredLogger {
	return New().Sugar()
}
