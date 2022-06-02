package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogClient struct {
	Logger *zap.Logger
}

func connectZap(level string) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "message",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	config := zap.Config{
		Level:            stringToLevel(level),
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}
	var err error
	client.Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
	defer client.Logger.Sync()
}
