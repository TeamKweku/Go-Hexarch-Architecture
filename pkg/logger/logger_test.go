package logger

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewZerologLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		isPrettyPrint bool
		wantConsole   bool
	}{
		{
			name:          "pretty print enabled",
			isPrettyPrint: true,
			wantConsole:   true,
		},
		{
			name:          "pretty print disabled",
			isPrettyPrint: false,
			wantConsole:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger := NewZerologLogger(tt.isPrettyPrint)
			require.NotNil(t, logger)
		})
	}
}

func TestZerologLogger_WithContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		ctx         context.Context
		wantTraceID interface{}
	}{
		{
			name:        "context with trace ID",
			ctx:         context.WithValue(context.Background(), traceIDKey, "test-trace-id"),
			wantTraceID: "test-trace-id",
		},
		{
			name:        "context without trace ID",
			ctx:         context.Background(),
			wantTraceID: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			baseLogger := zerolog.New(&buf)
			zLogger := &ZerologLogger{logger: baseLogger}
			contextLogger := zLogger.WithContext(tt.ctx)
			require.NotNil(t, contextLogger)

			contextLogger.Info(tt.ctx, "test message", nil)
			output := buf.String()

			if tt.wantTraceID != nil {
				assert.Contains(t, output, tt.wantTraceID)
			} else {
				assert.NotContains(t, output, "trace_id")
			}
		})
	}
}

func TestZerologLogger_LogLevels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		logFunc  func(Logger, context.Context)
		wantText string
		wantErr  bool
	}{
		{
			name: "info level",
			logFunc: func(l Logger, ctx context.Context) {
				l.Info(ctx, "info message", map[string]interface{}{"key": "value"})
			},
			wantText: `"level":"info"`,
		},
		{
			name: "error level",
			logFunc: func(l Logger, ctx context.Context) {
				l.Error(ctx, errors.New("test error"), "error message", map[string]interface{}{"key": "value"})
			},
			wantText: `"level":"error"`,
		},
		{
			name: "debug level",
			logFunc: func(l Logger, ctx context.Context) {
				l.Debug(ctx, "debug message", map[string]interface{}{"key": "value"})
			},
			wantText: `"level":"debug"`,
		},
		{
			name: "warn level",
			logFunc: func(l Logger, ctx context.Context) {
				l.Warn(ctx, "warn message", map[string]interface{}{"key": "value"})
			},
			wantText: `"level":"warn"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			baseLogger := zerolog.New(&buf)
			zLogger := &ZerologLogger{logger: baseLogger}
			ctx := context.Background()

			tt.logFunc(zLogger, ctx)
			output := buf.String()

			assert.Contains(t, output, tt.wantText)
			assert.Contains(t, output, `"key":"value"`)
		})
	}
}

func TestZerologLogger_Error(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	baseLogger := zerolog.New(&buf)
	zLogger := &ZerologLogger{logger: baseLogger}
	ctx := context.Background()
	testErr := errors.New("test error")

	zLogger.Error(ctx, testErr, "error message", map[string]interface{}{"key": "value"})
	output := buf.String()

	assert.Contains(t, output, `"level":"error"`)
	assert.Contains(t, output, `"error":"test error"`)
	assert.Contains(t, output, `"stack_trace"`)
	assert.Contains(t, output, `"key":"value"`)
}

func TestGetStackTrace(t *testing.T) {
	t.Parallel()

	// Function to help us generate a stack trace from a known location
	var trace string
	func() {
		trace = getStackTrace()
	}()

	// Basic validation
	if trace == "" {
		t.Fatal("Stack trace should not be empty")
	}

	// The trace should contain some expected elements
	t.Run("contains test file name", func(t *testing.T) {
		if !strings.Contains(trace, "logger_test.go") {
			t.Errorf("Stack trace should contain logger_test.go, got: %s", trace)
		}
	})

	t.Run("contains line numbers", func(t *testing.T) {
		if !strings.Contains(trace, ".go:") {
			t.Errorf("Stack trace should contain line numbers (.go:XX), got: %s", trace)
		}
	})

	t.Run("has multiple frames", func(t *testing.T) {
		frames := strings.Split(trace, "\n")
		if len(frames) < 2 {
			t.Error("Stack trace should have multiple frames")
		}
	})

	// Print the trace for debugging
	t.Logf("Stack trace:\n%s", trace)
}
