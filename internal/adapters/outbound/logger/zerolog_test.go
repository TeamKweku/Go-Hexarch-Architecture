package logger

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestZerologLogger(t *testing.T) {
	t.Parallel()

	isPrettyPrint := true

	logger := NewZerologLogger(isPrettyPrint)
	assert.NotNil(t, logger)
	assert.IsType(t, &ZerologLogger{}, logger)
}

func TestZerologLogger_Info(t *testing.T) {
	t.Parallel()

	// creat a buffer to captur th log output
	var buf bytes.Buffer
	testLogger := zerolog.New(&buf).With().Timestamp().Logger()

	zLogger := &ZerologLogger{
		logger: testLogger,
	}

	ctx := context.Background()
	msg := "test info message"
	fields := map[string]interface{}{"key": "value"}
	zLogger.Info(ctx, msg, fields)

	logOutput := buf.String()
	assert.Contains(t, logOutput, msg)
	assert.Contains(t, logOutput, "key")
	assert.Contains(t, logOutput, "value")
	assert.Contains(t, logOutput, "info")
}

func TestZerologLogger_Error(t *testing.T) {
	t.Parallel()

	// Create a buffer to capture the log output
	var buf bytes.Buffer
	testLogger := zerolog.New(&buf).With().Timestamp().Logger()

	zLogger := &ZerologLogger{
		logger: testLogger,
	}

	ctx := context.Background()
	err := errors.New("test error")
	msg := "test error message"
	fields := map[string]interface{}{"key": "value"}

	zLogger.Error(ctx, err, msg, fields)

	logOutput := buf.String()
	assert.Contains(t, logOutput, msg)
	assert.Contains(t, logOutput, "test error")
	assert.Contains(t, logOutput, "key")
	assert.Contains(t, logOutput, "value")
	assert.Contains(t, logOutput, "error")
}
