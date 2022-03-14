package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (lc LogClient) Debug(msg string, data interface{}) {
	var tags []zap.Field
	tags = append(tags, zap.Any("data", data))
	lc.Logger.Debug(msg, tags...)
	_ = lc.Logger.Sync()
}

func (lc LogClient) Error(msg string, err error) {
	var tags []zap.Field
	tags = append(tags, zap.NamedError("error", err))
	lc.Logger.Error(msg, tags...)
	_ = lc.Logger.Sync()
}

func (lc LogClient) Info(msg string) {
	lc.Logger.Info(msg)
	_ = lc.Logger.Sync()
}

func (lc LogClient) Warn(msg string) {
	lc.Logger.Warn(msg)
	_ = lc.Logger.Sync()
}

func (lc LogClient) Fatal(msg string) {
	lc.Logger.Fatal(msg)
	_ = lc.Logger.Sync()
}

func (lc LogClient) Panic(msg string) {
	lc.Logger.Panic(msg)
	_ = lc.Logger.Sync()
}

func stringToLevel(str string) zap.AtomicLevel {
	str = strings.ToLower(str)
	switch str {
	case "debug":
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "error":
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "info":
		return zap.NewAtomicLevel()
	default:
		return zap.NewAtomicLevelAt(zapcore.PanicLevel)
	}
}
