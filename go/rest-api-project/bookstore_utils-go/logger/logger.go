package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	envLogLevel  = "LOG_LEVEL"
	envLogOutput = "LOG_OUTPUT"
)

// Dont return a pointer to a struct when building a library
// Instead return a pointer to an interface which makes things easier to replace and test
type LoggerInterface interface {
	Info(string, ...zap.Field)
	Error(string, error, ...zap.Field)
	Printf(string, ...interface{}) // This method satisfies the Logger interface for elasticsearch
	Print(v ...interface{})        // This method satisfies the Logger interface for mysql
}

type loggerImpl struct {
	internalLogger *zap.Logger
}

var (
	logger LoggerInterface
)

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"},
		Level:       zap.NewAtomicLevelAt(getLevel()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "message",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	internalLogger, err := logConfig.Build()
	if err != nil {
		panic(err)
	}

	logger = &loggerImpl{internalLogger: internalLogger}
}

func getLevel() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(envLogLevel))) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func getOutput() string {
	out := strings.TrimSpace(os.Getenv(envLogOutput))
	if len(out) == 0 {
		return "stdout"
	}
	return out
}

func GetLogger() LoggerInterface {
	return logger
}

func (l loggerImpl) Info(msg string, tags ...zap.Field) {
	l.internalLogger.Info(msg, tags...)
	l.internalLogger.Sync() // sync is gonna flush the logs
}

func (l loggerImpl) Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.internalLogger.Error(msg, tags...)
	l.internalLogger.Sync() // sync is gonna flush the logs
}

// This method implements the Logger interface for elasticsearch
func (l loggerImpl) Printf(format string, v ...interface{}) {
	if len(v) == 0 {
		l.Info(format)
	} else {
		l.Info(fmt.Sprintf(format, v...))
	}
}

func (l loggerImpl) Print(v ...interface{}) {
	l.Info(fmt.Sprintf("%v", v))
}
