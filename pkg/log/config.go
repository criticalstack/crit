package log

import (
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"go.uber.org/zap/zapcore"
)

// EpochTimeEncoder serializes a time.Time to an integer representing the
// number of seconds since the Unix epoch.
func EpochTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	sec := t.Unix()
	enc.AppendInt64(sec)
}

// ColorEpochTimeEncoder serializes a time.Time to an integer representing the
// number of seconds since the Unix epoch and adds color.
func ColorEpochTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	sec := t.Unix()
	enc.AppendString(color.HiWhiteString(strconv.FormatInt(sec, 10)))
}

// CapitalFullNameEncoder serializes the logger name to an all-caps string. For
// example, gossip is serialized to "GOSSIP".
func CapitalFullNameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	loggerName = strings.ToUpper(loggerName)
	enc.AppendString(loggerName)
}

// CapitalColorFullNameEncoder serializes the logger name to an all-caps string
// and adds color. For example, gossip is serialized to "GOSSIP" and colored
// white.
func CapitalColorFullNameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	loggerName = color.HiWhiteString(strings.ToUpper(loggerName))
	enc.AppendString(loggerName)
}

func NewDefaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     ColorEpochTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     CapitalColorFullNameEncoder,
	}
}
