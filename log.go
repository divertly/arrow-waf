package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

func configureLogging(env string, level string) (zerolog.Logger, error) {
	var log zerolog.Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if G.Config.Log.Format == "text" && env != "prod" {
		log = zerolog.New(zerolog.ConsoleWriter{
			Out:             os.Stdout,
			FormatLevel:     formatLevel,
			FormatFieldName: formatFieldName,
			TimeFormat:      time.RFC3339,
			NoColor:         !G.Config.Log.Color,
		}).With().Timestamp().Logger()
		log.Info().Str(this()).Msg("Using console logger")
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		writer := diode.NewWriter(os.Stderr, 1000, 10*time.Millisecond, func(missed int) {
			fmt.Printf("Logger dropped %d messages", missed)
		})
		log = zerolog.New(writer).With().Timestamp().Logger()
		log.Info().Str(this()).Msg("Using diode logger")
	}
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Error().Err(err)
		return log, err
	}
	zerolog.SetGlobalLevel(lvl)
	log = log.Level(lvl)
	return log, nil
}

func formatLevel(i interface{}) string {
	col := color.Bold //nolint:ineffassign
	switch i.(string) {
	case "trace":
		col = color.FgBlue
	case "debug":
		col = color.FgCyan
	case "info":
		col = color.FgGreen
	case "warn":
		col = color.FgYellow
	case "error":
		col = color.FgRed
	case "fatal", "panic":
		col = color.FgMagenta
	default:
		col = color.Bold
	}
	s := color.New(col).SprintFunc()
	return fmt.Sprintf("| %-8s |", s(strings.ToUpper(i.(string))))
}

func formatFieldName(i interface{}) string {
	return fmt.Sprintf("%s=", i)
}
