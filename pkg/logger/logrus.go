package logger

import (
	"context"
	"fmt"
	"github.com/arraisi/demogo/pkg/constant"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	logrusFieldRequestID string = "request_id"
)

type LoggerImpl struct {
	*logrus.Logger
}

func RegisterLogrus(lvl string) (Logger, error) {
	log := logrus.New()
	formatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	}
	log.SetFormatter(formatter)
	log.SetReportCaller(true)
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return nil, errors.Wrap(err, "Error init logrus")
	}
	log.SetLevel(level)
	return &LoggerImpl{
		Logger: log,
	}, nil
}

func (l *LoggerImpl) WithFields(ctx context.Context, fields map[string]interface{}) EntryLogger {
	entry := l.Logger.WithContext(ctx).
		WithField(logrusFieldRequestID, ctx.Value(constant.REQUEST_ID)).
		WithFields(fields)
	return entry
}

func (l *LoggerImpl) WithContext(ctx context.Context) EntryLogger {
	entry := l.Logger.WithContext(ctx).
		WithFields(
			logrus.Fields{
				logrusFieldRequestID: ctx.Value(constant.REQUEST_ID),
			})
	return entry
}

func (l *LoggerImpl) InfoTrace(span ddtrace.Span, format string, args ...interface{}) {
	l.Logf(logrus.InfoLevel, "%d %d msg: %s", span.Context().TraceID(), span.Context().SpanID(), fmt.Sprintf(format, args...))
}

func (l *LoggerImpl) WarnTrace(span ddtrace.Span, format string, args ...interface{}) {
	l.Logf(logrus.WarnLevel, "%d %d msg: %s", span.Context().TraceID(), span.Context().SpanID(), fmt.Sprintf(format, args...))
}

func (l *LoggerImpl) ErrorTrace(span ddtrace.Span, format string, args ...interface{}) {
	l.Logf(logrus.ErrorLevel, "%d %d msg: %s", span.Context().TraceID(), span.Context().SpanID(), fmt.Sprintf(format, args...))
}
