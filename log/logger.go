package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
)

type (
	// Logger Logger Interface
	Logger interface {
		Info(msg string, args ...interface{})
		Error(msg string, args ...interface{})
		RuntimeFatalError(msg string, args ...interface{})
		Debug(msg string, args ...interface{})
		Warn(msg string, args ...interface{})
	}
	Log struct {
		Provider *zap.Logger
	}
)

// NewLogger construct
func NewLogger() *Log {
	logger, _ := zap.NewProduction()
	return &Log{Provider: logger.WithOptions()}
}

func (l *Log) Info(msg string, args ...interface{}) {
	l.Provider.Info(msg, l.toZapField(args)...)
}

func (l *Log) Error(msg string, args ...interface{}) {
	l.Provider.Error(msg, l.toZapField(args)...)
}

func (l *Log) RuntimeFatalError(msg string, args ...interface{}) {
	l.Provider.Fatal(msg, l.toZapField(args)...)
}

func (l *Log) Warn(msg string, args ...interface{}) {
	l.Provider.Warn(msg, l.toZapField(args)...)
}

func (l *Log) Debug(msg string, args ...interface{}) {
	l.Provider.Debug(msg, l.toZapField(args)...)
}

func (l *Log) toZapField(args []interface{}) []zap.Field {
	var f []zap.Field
	for i, v := range args {
		f = append(f, zap.Field{
			Key:    strconv.Itoa(i),
			Type:   zapcore.StringType,
			String: fmt.Sprintln(v),
		})
	}
	return f
}
