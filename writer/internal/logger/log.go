package myLog

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log *MyLogger

func init() {
	Log = initLogger()
}

type MyLogger struct {
	Lg zerolog.Logger
}

func initLogger() *MyLogger {
	return &MyLogger{
		Lg: zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.UnixDate,
		}),
	}
}
func (l *MyLogger) Infof(mes string, v ...interface{}) {
	if len(v) == 0 {
		l.Lg.Info().Msgf(mes)
		return
	}
	l.Lg.Info().Msgf(mes, v)
}

func (l *MyLogger) Debugf(mes string, v ...interface{}) {
	if len(v) == 0 {
		l.Lg.Debug().Msgf(mes)
		return
	}
	l.Lg.Debug().Msgf(mes, v)
}

func (l *MyLogger) Errorf(mes string, v ...interface{}) {
	if len(v) == 0 {
		l.Lg.Error().Msgf(mes)
		return
	}
	l.Lg.Error().Msgf(mes, v)
}

func (l *MyLogger) Warnf(mes string, v ...interface{}) {
	if len(v) == 0 {
		l.Lg.Warn().Msgf(mes)
		return
	}
	l.Lg.Warn().Msgf(mes, v)
}

func (l *MyLogger) Fatalf(mes string, v ...interface{}) {
	if len(v) == 0 {
		l.Lg.Fatal().Msgf(mes)
		return
	}
	l.Lg.Fatal().Msgf(mes, v)
}
