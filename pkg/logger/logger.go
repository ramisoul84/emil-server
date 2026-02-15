package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ramisoul84/emil-server/config"
	"github.com/rs/zerolog"
)

var (
	globalLogger Logger
	once         sync.Once
)

type Logger interface {
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Error() *zerolog.Event
	Warn() *zerolog.Event
	Fatal() *zerolog.Event
	WithFields(fields map[string]any) Logger
}

type logger struct {
	zerolog.Logger
}

func New(cfg *config.Config) Logger {
	var level zerolog.Level
	var output io.Writer

	zerolog.TimeFieldFormat = time.RFC3339Nano

	level = parseLevel(cfg.Logging.Level)

	if cfg.Logging.Output == "file" {
		file, err := os.OpenFile(cfg.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to create log file: %v, using stdout\n", err)
			output = os.Stdout
		} else {
			output = file
		}
	} else {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC822,
		}
	}

	zlog := zerolog.New(output).
		Level(level).
		With().
		//Str("service", cfg.App.Name).
		//Str("version", cfg.App.Version).
		Timestamp().
		Logger()

	return &logger{
		Logger: zlog,
	}
}

func InitGlobal(cfg *config.Config) {
	once.Do(func() {
		globalLogger = New(cfg)
	})
}

func Get() Logger {
	return globalLogger
}

func Debug() *zerolog.Event {
	return Get().Debug()
}

func Info() *zerolog.Event {
	return Get().Info()
}

func Warn() *zerolog.Event {
	return Get().Warn()
}

func Error() *zerolog.Event {
	return Get().Error()
}

func Fatal() *zerolog.Event {
	return Get().Fatal()
}

func (l *logger) WithFields(fields map[string]any) Logger {
	ctx := l.Logger.With()
	for key, value := range fields {
		ctx = ctx.Interface(key, value)
	}
	return &logger{
		Logger: ctx.Logger(),
	}
}

// parseLevel parses log level string to zerolog.Level
func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
