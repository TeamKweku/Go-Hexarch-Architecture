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

func NewZerologLogger(isPrettyPrint bool) *ZerologLogger {
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
