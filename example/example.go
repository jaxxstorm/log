// Package example demonstrates the public logging API.
//
// The logger defaults to pretty output when writing to an interactive terminal
// and to JSON when output is redirected or written to a non-terminal writer.
package example

import (
	"io"

	log "github.com/jaxxstorm/log"
)

type Options struct {
	Output io.Writer
	Level  log.Level
	Format log.Format
}

func Run(options Options) error {
	logger, err := log.New(log.Config{
		Output: options.Output,
		Level:  options.Level,
		Format: options.Format,
	})
	if err != nil {
		return err
	}
	defer logger.Close()

	requestLogger := logger.With(
		log.String("component", "example"),
		log.String("request_id", "req-42"),
	)

	requestLogger.Info("starting request", log.String("route", "/healthz"))
	requestLogger.Error("request failed", log.Int("status", 503))

	return nil
}
