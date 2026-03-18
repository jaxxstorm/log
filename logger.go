package log

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	base    *zap.Logger
	fields  []Field
	closeFn func() error
	exitFn  func(int)
}

type noopFatalHook struct{}

func New(config Config) (*Logger, error) {
	resolved, err := config.resolve()
	if err != nil {
		return nil, err
	}

	base := newBaseLogger(resolved.core())

	return &Logger{
		base:    base,
		closeFn: resolved.closeFn,
		exitFn:  os.Exit,
	}, nil
}

func newBaseLogger(core zapcore.Core) *zap.Logger {
	return zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.WithFatalHook(noopFatalHook{}),
	)
}

func (noopFatalHook) OnWrite(*zapcore.CheckedEntry, []zapcore.Field) {}

func (l *Logger) With(fields ...Field) *Logger {
	if l == nil {
		return nil
	}

	return &Logger{
		base:    l.base,
		fields:  mergeFields(l.fields, fields),
		closeFn: l.closeFn,
		exitFn:  l.exitFn,
	}
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.log(func(logger *zap.Logger, msg string, fields ...zap.Field) {
		logger.Debug(msg, fields...)
	}, msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.log(func(logger *zap.Logger, msg string, fields ...zap.Field) {
		logger.Info(msg, fields...)
	}, msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.log(func(logger *zap.Logger, msg string, fields ...zap.Field) {
		logger.Warn(msg, fields...)
	}, msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.log(func(logger *zap.Logger, msg string, fields ...zap.Field) {
		logger.Error(msg, fields...)
	}, msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	if l == nil || l.base == nil {
		return
	}

	merged := mergeFields(l.fields, fields)
	l.base.Fatal(msg, toZapFields(merged)...)
	l.finalizeFatal()
}

func (l *Logger) Sync() error {
	if l == nil || l.base == nil {
		return nil
	}

	return swallowSyncError(l.base.Sync())
}

func (l *Logger) Close() error {
	if l == nil {
		return nil
	}

	if err := l.Sync(); err != nil {
		return err
	}

	if l.closeFn == nil {
		return nil
	}

	return l.closeFn()
}

func (l *Logger) log(fn func(*zap.Logger, string, ...zap.Field), msg string, fields ...Field) {
	if l == nil || l.base == nil {
		return
	}

	merged := mergeFields(l.fields, fields)
	fn(l.base, msg, toZapFields(merged)...)
}

func (l *Logger) finalizeFatal() {
	_ = l.Sync()

	if l.closeFn != nil {
		_ = l.closeFn()
	}

	if l.exitFn != nil {
		l.exitFn(1)
	}
}

func mergeFields(base []Field, extras []Field) []Field {
	if len(base) == 0 && len(extras) == 0 {
		return nil
	}

	merged := make([]Field, 0, len(base)+len(extras))
	indexByKey := make(map[string]int, len(base)+len(extras))

	for _, field := range base {
		key := field.zap.Key
		if key == "" {
			merged = append(merged, field)
			continue
		}

		indexByKey[key] = len(merged)
		merged = append(merged, field)
	}

	for _, field := range extras {
		key := field.zap.Key
		if key == "" {
			merged = append(merged, field)
			continue
		}

		if idx, ok := indexByKey[key]; ok {
			merged[idx] = field
			continue
		}

		indexByKey[key] = len(merged)
		merged = append(merged, field)
	}

	return merged
}

func toZapFields(fields []Field) []zap.Field {
	if len(fields) == 0 {
		return nil
	}

	out := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		out = append(out, field.zap)
	}

	return out
}

func swallowSyncError(err error) error {
	if err == nil {
		return nil
	}

	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		return nil
	}

	return err
}
