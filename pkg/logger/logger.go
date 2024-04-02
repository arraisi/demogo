package logger

import (
	"context"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"

	"github.com/pkg/errors"
)

var Log = defaultLogger()

func defaultLogger() Logger {
	logger, err := InitLogger("info")
	if err != nil {
		return nil
	}
	return logger
}

func SetLogger(log Logger) {
	Log = log
}

// build logrus logger
func InitLogger(lc string) (Logger, error) {
	logger, err := RegisterLogrus(lc)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return logger, nil
}

func New(ctx context.Context) EntryLogger {
	return Log.WithContext(ctx)
}

type Logger interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	WithFields(ctx context.Context, fields map[string]interface{}) EntryLogger
	WithContext(ctx context.Context) EntryLogger
	InfoTrace(span ddtrace.Span, format string, args ...interface{})
	WarnTrace(span ddtrace.Span, format string, args ...interface{})
	ErrorTrace(span ddtrace.Span, format string, args ...interface{})
}

type EntryLogger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}
