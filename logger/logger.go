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

// GetLogger returns the logger instance to use in other packages
func Get() *zap.Logger {
	return logger
}

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

func NewSugared() *zap.SugaredLogger {
	return New().Sugar()
}
