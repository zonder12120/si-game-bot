package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/zonder12120/tg-quiz/internal/config"
)

func NewLogger(config config.Logging) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, err
	}

	zerolog.SetGlobalLevel(level)

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &logger, nil
}
