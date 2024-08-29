package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type ZerologLogger struct {
	logger zerolog.Logger
}

func NewZerologLogger() *ZerologLogger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return &ZerologLogger{
		logger: logger,
	}
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
	z.logger.Error().Err(err).Fields(fields).Msg(msg)
}

// func (z *ZerologLogger) WithContext(ctx context.Context) *ZerologLogger {
// 	return &ZerologLogger{
// 		logger: z.logger.With().Logger(),
// 	}
// }
