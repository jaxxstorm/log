package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/jaxxstorm/log"
	"github.com/jaxxstorm/log/example"
)

func main() {
	var (
		format = flag.String("format", "auto", "output format: auto, pretty, or json")
		level  = flag.String("level", "info", "log level: debug, info, warn, or error")
	)
	flag.Parse()

	if err := example.Run(example.Options{
		Output: os.Stdout,
		Level:  log.Level(*level),
		Format: log.Format(*format),
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
