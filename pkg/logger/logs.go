package logs

import "go.uber.org/zap"

func Debug(args ...interface{}) {
	zap.S().Debug(args)
}

func Debugf(msg string, args ...interface{}) {
	zap.S().Debugf(msg, args...)
}

func Debugw(msg string, args ...interface{}) {
	zap.S().Debugw(msg, args...)
}

func Info(args ...interface{}) {
	zap.S().Info(args)
}

func Infof(msg string, args ...interface{}) {
	zap.S().Infof(msg, args...)
}

func Infow(msg string, args ...interface{}) {
	zap.S().Infow(msg, args...)
}

func Warn(args ...interface{}) {
	zap.S().Warn(args)
}

func Warnf(msg string, args ...interface{}) {
	zap.S().Warnf(msg, args...)
}

func Warnw(msg string, args ...interface{}) {
	zap.S().Warnw(msg, args...)
}

func Error(args ...interface{}) {
	zap.S().Error(args)
}

func Errorf(msg string, args ...interface{}) {
	zap.S().Errorf(msg, args...)
}

func Errorw(msg string, args ...interface{}) {
	zap.S().Errorw(msg, args...)
}

func Fatal(args ...interface{}) {
	zap.S().Fatal(args)
}

func Fatalf(msg string, args ...interface{}) {
	zap.S().Fatalf(msg, args...)
}

func Fatalw(msg string, args ...interface{}) {
	zap.S().Fatalw(msg, args...)
}

func Panic(args ...interface{}) {
	zap.S().Panic(args)
}

func Panicf(msg string, args ...interface{}) {
	zap.S().Panicf(msg, args...)
}

func Panicw(msg string, args ...interface{}) {
	zap.S().Panicw(msg, args...)
}
