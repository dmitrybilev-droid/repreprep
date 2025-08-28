package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
	return zerolog.New(output).With().Timestamp().Logger()
}