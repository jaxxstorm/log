package log

import (
	"time"

	"go.uber.org/zap"
)

type Field struct {
	zap zap.Field
}

func String(key, value string) Field {
	return Field{zap: zap.String(key, value)}
}

func Int(key string, value int) Field {
	return Field{zap: zap.Int(key, value)}
}

func Bool(key string, value bool) Field {
	return Field{zap: zap.Bool(key, value)}
}

func Duration(key string, value time.Duration) Field {
	return Field{zap: zap.Duration(key, value)}
}

func Error(err error) Field {
	return Field{zap: zap.Error(err)}
}

func Any(key string, value any) Field {
	return Field{zap: zap.Any(key, value)}
}
