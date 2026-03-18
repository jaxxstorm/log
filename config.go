package log

import (
	"errors"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"

	"github.com/jaxxstorm/log/internal/pretty"
)

type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	FatalLevel Level = "fatal"
)

type Format string

const (
	AutoFormat   Format = "auto"
	JSONFormat   Format = "json"
	PrettyFormat Format = "pretty"
)

type Config struct {
	Level      Level
	Format     Format
	Output     io.Writer
	OutputPath string
}

type resolvedConfig struct {
	level      zapcore.Level
	format     Format
	writer     zapcore.WriteSyncer
	closeFn    func() error
	encoderCfg zapcore.EncoderConfig
}

func (c Config) resolve() (resolvedConfig, error) {
	level, err := c.level()
	if err != nil {
		return resolvedConfig{}, err
	}

	writer, isTTY, closeFn, err := c.writer()
	if err != nil {
		return resolvedConfig{}, err
	}

	format, err := c.formatFor(isTTY)
	if err != nil {
		if closeFn != nil {
			_ = closeFn()
		}
		return resolvedConfig{}, err
	}

	return resolvedConfig{
		level:   level,
		format:  format,
		writer:  writer,
		closeFn: closeFn,
		encoderCfg: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}, nil
}

func (c Config) level() (zapcore.Level, error) {
	switch c.Level {
	case "", InfoLevel:
		return zapcore.InfoLevel, nil
	case DebugLevel:
		return zapcore.DebugLevel, nil
	case WarnLevel:
		return zapcore.WarnLevel, nil
	case ErrorLevel:
		return zapcore.ErrorLevel, nil
	case FatalLevel:
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unsupported level %q", c.Level)
	}
}

func (c Config) formatFor(isTTY bool) (Format, error) {
	switch c.Format {
	case "", AutoFormat, JSONFormat, PrettyFormat:
	default:
		return "", fmt.Errorf("unsupported format %q", c.Format)
	}

	if c.Format == "" || c.Format == AutoFormat {
		if isTTY {
			return PrettyFormat, nil
		}

		return JSONFormat, nil
	}

	return c.Format, nil
}

func (c Config) writer() (zapcore.WriteSyncer, bool, func() error, error) {
	if c.Output != nil && c.OutputPath != "" {
		return nil, false, nil, errors.New("output and output_path cannot be set together")
	}

	if c.OutputPath != "" {
		file, err := os.OpenFile(c.OutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, false, nil, fmt.Errorf("open output path: %w", err)
		}

		return zapcore.AddSync(file), false, file.Close, nil
	}

	writer := c.Output
	if writer == nil {
		writer = os.Stdout
	}

	return zapcore.AddSync(writer), writerIsTTY(writer), nil, nil
}

func writerIsTTY(writer io.Writer) bool {
	type fdWriter interface {
		Fd() uintptr
	}

	fileWriter, ok := writer.(fdWriter)
	if !ok {
		return false
	}

	return term.IsTerminal(int(fileWriter.Fd()))
}

func (r resolvedConfig) core() zapcore.Core {
	level := zap.NewAtomicLevelAt(r.level)

	if r.format == PrettyFormat {
		return pretty.NewCore(r.writer, level)
	}

	return zapcore.NewCore(zapcore.NewJSONEncoder(r.encoderCfg), r.writer, level)
}
