package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	level = zap.NewAtomicLevel()

	log = zap.New(zapcore.NewCore(
		NewEncoder(NewDefaultEncoderConfig()),
		zapcore.AddSync(os.Stderr),
		level,
	), zap.AddCaller(), zap.AddCallerSkip(1))
)

// NewLogger creates a new child logger with the provided namespace.
func NewLogger(ns string) *zap.Logger {
	encoder := NewEncoder(NewDefaultEncoderConfig())
	encoder.OpenNamespace(ns)
	return log.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stderr),
			level,
		)
	}), zap.AddCaller())
}

// NewLoggerWithLevel creates a new child logger with the provided namespace
// and level. Since this specifies a level, it overrides the global package
// level for this child logger only.
func NewLoggerWithLevel(ns string, lvl zapcore.Level) *zap.Logger {
	encoder := NewEncoder(NewDefaultEncoderConfig())
	encoder.OpenNamespace(ns)
	return log.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stderr),
			lvl,
		)
	}))
}

func Level() zapcore.Level {
	return level.Level()
}

func SetLevel(lvl zapcore.Level) {
	level.SetLevel(lvl)
}

func Debug(msg string, fields ...zapcore.Field) {
	log.Debug(msg, fields...)
}

func Debugf(format string, args ...interface{}) {
	log.Debug(fmt.Sprintf(format, args...))
}

func Info(msg string, fields ...zapcore.Field) {
	log.Info(msg, fields...)
}

func Infof(format string, args ...interface{}) {
	log.Info(fmt.Sprintf(format, args...))
}

func Warn(msg string, fields ...zapcore.Field) {
	log.Warn(msg, fields...)
}

func Warnf(format string, args ...interface{}) {
	log.Warn(fmt.Sprintf(format, args...))
}

func Error(msg string, fields ...zapcore.Field) {
	log.Error(msg, fields...)
}

func Errorf(format string, args ...interface{}) {
	log.Error(fmt.Sprintf(format, args...))
}

func Fatal(msg interface{}, fields ...zapcore.Field) {
	switch t := msg.(type) {
	case string:
		log.Fatal(t, fields...)
	case error:
		log.Fatal(t.Error(), fields...)
	default:
		log.Fatal(fmt.Sprintf("%+v", msg), fields...)
	}
}

func Fatalf(format string, args ...interface{}) {
	log.Fatal(fmt.Sprintf(format, args...))
}
