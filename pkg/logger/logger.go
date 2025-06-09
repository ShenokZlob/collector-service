package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "ts",
		LevelKey:   "level",
		NameKey:    "logger",
		CallerKey:  "caller",
		MessageKey: "msg",
		// StacktraceKey: "stacktrace",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.CapitalLevelEncoder, // INFO, WARN, ERROR…
		EncodeTime:   zapcore.ISO8601TimeEncoder,  // 2006-01-02T15:04:05.000Z0700
		EncodeCaller: zapcore.ShortCallerEncoder,  // file.go:42
	}

	cfg := zap.Config{
		Level:            level,
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build(zap.AddCaller()) //, zap.AddStacktrace(zap.ErrorLevel)
	if err != nil {
		return nil, err
	}

	return logger, nil
}
