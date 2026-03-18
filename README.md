# log

`log` is a small Go logging package I use for my own projects. It wraps Zap with
simple configuration, structured fields, and automatic pretty vs JSON output.

This repository is for personal use. I am not looking for outside pull requests.

## What it does

- Defaults to pretty output on a TTY and JSON when output is redirected.
- Supports `debug`, `info`, `warn`, and `error` levels.
- Adds caller information and timestamps automatically.
- Lets you attach persistent structured fields with `With(...)`.
- Can write to any `io.Writer` or directly to a file path.

## Install

```bash
go get github.com/jaxxstorm/log
```

## Basic usage

```go
package main

import (
	log "github.com/jaxxstorm/log"
)

func main() {
	logger, err := log.New(log.Config{})
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	logger.Info("service started", log.String("addr", ":8080"))
}
```

## Add context with `With`

```go
package main

import (
	log "github.com/jaxxstorm/log"
)

func main() {
	logger, err := log.New(log.Config{
		Level: log.DebugLevel,
	})
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	requestLogger := logger.With(
		log.String("component", "api"),
		log.String("request_id", "req-42"),
	)

	requestLogger.Info("request started", log.String("route", "/healthz"))
	requestLogger.Error("request failed", log.Int("status", 503))
}
```

## Control format and destination

```go
package main

import (
	"os"

	log "github.com/jaxxstorm/log"
)

func main() {
	logger, err := log.New(log.Config{
		Level:  log.InfoLevel,
		Format: log.JSONFormat,
		Output: os.Stdout,
	})
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	logger.Info("json log line", log.Bool("ok", true))
}
```

```go
logger, err := log.New(log.Config{
	OutputPath: "app.log",
})
if err != nil {
	panic(err)
}
defer logger.Close()
```

## Available options

`log.Config` supports:

- `Level`: `debug`, `info`, `warn`, or `error`
- `Format`: `auto`, `pretty`, or `json`
- `Output`: any `io.Writer`
- `OutputPath`: a file to append logs to

If `Format` is `auto` or left empty, the logger chooses pretty output for a
terminal and JSON otherwise.
