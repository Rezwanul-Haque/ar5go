package logger

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

const (
	LVLINFO  = "INFO"
	LVLERROR = "ERROR"
	LVLDEBUG = "DEBUG"
)

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"},
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	if log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

func Info(msg string, tags ...zap.Field) {
	log.Info(msg, tags...)
	_ = log.Sync()
}

func InfoAsJson(msg string, data ...interface{}) {
	prettyPrint(LVLINFO, msg, data)
}
func ErrorAsJson(msg string, data ...interface{}) {
	prettyPrint(LVLERROR, msg, data)
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.Error(msg, tags...)
	_ = log.Sync()
}

func prettyPrint(level string, msg string, data ...interface{}) {
	if r, err := json.MarshalIndent(&data, "", "  "); err == nil {
		fmt.Printf("[%v] %v %v: \n %v\n", level, time.Now(), msg, string(r))
	}
}
