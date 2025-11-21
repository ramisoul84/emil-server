package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(args ...any)
	Infof(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	WithFields(fields map[string]any) Logger
	WithError(err error) Logger
}

type logger struct {
	entry *logrus.Entry
}

func New(environment, serviceName, version string) Logger {
	log := logrus.New()

	if environment == "development" {
		log.SetLevel(logrus.InfoLevel)
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			ForceColors:      true,
			PadLevelText:     true,
			QuoteEmptyFields: true,
		})
	} else {
		log.SetLevel(logrus.ErrorLevel)
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
		})
	}

	log.SetOutput(os.Stdout)

	entry := log.WithFields(logrus.Fields{
		"service": serviceName,
		"version": version,
	})

	return &logger{entry: entry}
}

func (l *logger) Info(args ...any) {
	l.entry.Info(args...)
}

func (l *logger) Infof(format string, args ...any) {
	l.entry.Infof(format, args...)
}

func (l *logger) Error(args ...any) {
	l.entry.Error(args...)
}

func (l *logger) Errorf(format string, args ...any) {
	l.entry.Errorf(format, args...)
}

func (l *logger) Warn(args ...any) {
	l.entry.Warn(args...)
}

func (l *logger) Warnf(format string, args ...any) {
	l.entry.Warnf(format, args...)
}

func (l *logger) WithFields(fields map[string]any) Logger {
	return &logger{
		entry: l.entry.WithFields(logrus.Fields(fields)),
	}
}

func (l *logger) WithError(err error) Logger {
	return &logger{entry: l.entry.WithError(err)}
}
