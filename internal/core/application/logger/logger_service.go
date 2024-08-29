package logger

import (
	"context"

	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/outbound"
)

type LoggerService struct {
	logger outbound.Logger
}

func NewLoggerService(logger outbound.Logger) *LoggerService {
	return &LoggerService{
		logger: logger,
	}
}

func (ls *LoggerService) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	ls.logger.Info(ctx, msg, fields)
}

func (ls *LoggerService) Error(
	ctx context.Context,
	err error,
	msg string,
	fields map[string]interface{},
) {
	ls.logger.Error(ctx, err, msg, fields)
}

//
// func (ls *LoggerService) WithContext(ctx context.Context) *LoggerService {
// 	return &LoggerService{logger: ls.logger.WithContext(ctx)}
// }
