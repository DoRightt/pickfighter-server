package logger

import (
	"log"
	"os"

	"github.com/spf13/viper"
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

// Get returns the logger instance to use in other packages
func Get() *zap.Logger {
	return logger
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
		config.OutputPaths = []string{"logger/logs/log.json"}
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

// ScraperLogger creates a new instance of Zap logger specifically configured for the scraper module.
// It configures the logger to output logs in JSON format to the file "logger/logs/scraper-log.json".
// The log level is set to Debug. The returned SugaredLogger can be used for logging in the scraper module.
// If the log file already exists, its contents will be truncated; otherwise, a new file will be created.
// An error is returned if there is a failure in creating or configuring the logger.
func ScraperLogger() (*zap.SugaredLogger, error) {
	config := zap.NewDevelopmentConfig()

	config.Encoding = "json"
	config.OutputPaths = []string{"logger/logs/scraper-log.json"}
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

	file, err := os.OpenFile("logger/logs/scraper-log.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
		return nil, err
	}
	defer file.Close()

	l, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to run zap logger: %s", err)
		return nil, err
	}

	return l.Sugar(), nil
}

// NewSugared returns SugaredLogger version
func NewSugared() *zap.SugaredLogger {
	return New().Sugar()
}
