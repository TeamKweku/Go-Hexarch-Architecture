package logger

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields map[string]interface{})
	Error(ctx context.Context, err error, msg string, fields map[string]interface{})
	Debug(ctx context.Context, msg string, fields map[string]interface{})
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	Fatal(ctx context.Context, msg string, fields map[string]interface{})
	WithContext(ctx context.Context) Logger
}

type ZerologLogger struct {
	logger zerolog.Logger
}

// verify at compile time the ZerologLogger implements Logger
var _ Logger = (*ZerologLogger)(nil)

func NewZerologLogger(isPrettyPrint bool) Logger {
	// Custom time format: "2023-05-25 14:30:00 EDT"
	customTimeFormat := "2006-01-02 15:04:05 MST"

	zerolog.TimeFieldFormat = customTimeFormat

	// render pretty printing conditional based on environment
	var logger zerolog.Logger
	if isPrettyPrint {
		logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: customTimeFormat,
		})
	} else {
		logger = log.Logger
	}

	// Adding caller information
	logger = logger.With().Caller().Logger()

	return &ZerologLogger{
		logger: logger,
	}
}

type contextKey string

const traceIDKey = contextKey("trace_id")

func (z *ZerologLogger) WithContext(ctx context.Context) Logger {
	logger := z.logger

	// Extract trace ID from context
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		// Create new logger with trace ID field
		logger = logger.With().
			Str("trace_id", traceID).
			Logger()
		return &ZerologLogger{logger: logger}
	}
	return &ZerologLogger{logger: logger}
}

func (z *ZerologLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	z.logger.Info().Fields(fields).Msg(msg)
}

func (z *ZerologLogger) Error(
	ctx context.Context,
	err error,
	msg string,
	fields map[string]interface{},
) {
	z.logger.Error().
		Err(err).
		Fields(fields).
		Str("stack_trace", getStackTrace()).
		Msg(msg)
}

func (z *ZerologLogger) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	z.logger.Debug().
		Fields(fields).
		Msg(msg)
}

func (z *ZerologLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	z.logger.Warn().
		Fields(fields).
		Msg(msg)
}

func (z *ZerologLogger) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	z.logger.Fatal().
		Fields(fields).
		Str("stack_trace", getStackTrace()).
		Msg(msg)
}

// Helper function for getting stack traces
func getStackTrace() string {
	// Capture the full stack trace
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	// Split the stack into lines
	lines := strings.Split(stack, "\n")

	// Skip the first line (contains "goroutine X [running]")
	if len(lines) > 0 {
		lines = lines[1:]
	}

	// Skip the frame for this function
	if len(lines) >= 2 {
		lines = lines[2:]
	}

	// Rebuild the stack trace
	return strings.Join(lines, "\n")
}
