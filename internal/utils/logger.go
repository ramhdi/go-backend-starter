package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogger(env string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
