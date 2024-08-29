package outbound

import (
	"context"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields map[string]interface{})
	Error(ctx context.Context, err error, msg string, fields map[string]interface{})
	// WithContext(ctx context.Context)
}
