// logger/logger.go
package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Interface interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

type Logger struct {
	z *zap.SugaredLogger
}

var (
	logInstance *Logger
	once        sync.Once
)

// Init инициализирует логгер с заданным уровнем логирования (debug/info/warn/error/fatal)
func Init(level string) {
	once.Do(func() {
		zapLevel := parseLevel(level)

		cfg := zap.Config{
			Encoding:         "console", // или "json"
			Level:            zap.NewAtomicLevelAt(zapLevel),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			Development: true,
		}

		logger, err := cfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Logger init error: %v\n", err)
			os.Exit(1)
		}

		logInstance = &Logger{z: logger.Sugar()}
	})
}

// L возвращает глобальный логгер
func L() Interface {
	if logInstance == nil {
		Init("info") // дефолт
	}
	return logInstance
}

// реализация методов интерфейса

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.z.Debugf(msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.z.Infof(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.z.Warnf(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.z.Errorf(msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.z.Fatalf(msg, args...)
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
