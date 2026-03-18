package pretty

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap/zapcore"
)

type Core struct {
	writer zapcore.WriteSyncer
	level  zapcore.LevelEnabler
	fields []zapcore.Field
	styles styles
}

type styles struct {
	timestamp lipgloss.Style
	caller    lipgloss.Style
	message   lipgloss.Style
	key       lipgloss.Style
	debug     lipgloss.Style
	info      lipgloss.Style
	warn      lipgloss.Style
	err       lipgloss.Style
	fatal     lipgloss.Style
}

func NewCore(writer zapcore.WriteSyncer, level zapcore.LevelEnabler) *Core {
	return &Core{
		writer: writer,
		level:  level,
		styles: styles{
			timestamp: lipgloss.NewStyle().Faint(true),
			caller:    lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			message:   lipgloss.NewStyle().Bold(true),
			key:       lipgloss.NewStyle().Foreground(lipgloss.Color("12")),
			debug:     lipgloss.NewStyle().Foreground(lipgloss.Color("8")).SetString("DEBUG"),
			info:      lipgloss.NewStyle().Foreground(lipgloss.Color("10")).SetString("INFO"),
			warn:      lipgloss.NewStyle().Foreground(lipgloss.Color("11")).SetString("WARN"),
			err:       lipgloss.NewStyle().Foreground(lipgloss.Color("9")).SetString("ERROR"),
			fatal:     lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).SetString("FATAL"),
		},
	}
}

func (c *Core) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level)
}

func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	clone := *c
	clone.fields = append(append([]zapcore.Field{}, c.fields...), fields...)
	return &clone
}

func (c *Core) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !c.Enabled(entry.Level) {
		return checked
	}

	return checked.AddCore(entry, c)
}

func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	encoded := zapcore.NewMapObjectEncoder()

	for _, field := range c.fields {
		field.AddTo(encoded)
	}

	for _, field := range fields {
		field.AddTo(encoded)
	}

	line := c.render(entry, encoded.Fields)
	if _, err := c.writer.Write([]byte(line + "\n")); err != nil {
		return err
	}

	if entry.Level > zapcore.ErrorLevel {
		_ = c.Sync()
	}

	return nil
}

func (c *Core) Sync() error {
	return c.writer.Sync()
}

func (c *Core) render(entry zapcore.Entry, fields map[string]any) string {
	parts := []string{
		c.styles.timestamp.Render(entry.Time.Format("2006-01-02T15:04:05Z07:00")),
		c.renderLevel(entry.Level),
		c.styles.caller.Render(entry.Caller.TrimmedPath()),
		c.styles.message.Render(entry.Message),
	}

	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", c.styles.key.Render(key), formatValue(fields[key])))
	}

	return strings.Join(parts, " ")
}

func (c *Core) renderLevel(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return c.styles.debug.String()
	case zapcore.InfoLevel:
		return c.styles.info.String()
	case zapcore.WarnLevel:
		return c.styles.warn.String()
	case zapcore.ErrorLevel:
		return c.styles.err.String()
	case zapcore.FatalLevel:
		return c.styles.fatal.String()
	default:
		return strings.ToUpper(level.String())
	}
}

func formatValue(value any) string {
	switch typed := value.(type) {
	case string:
		if strings.ContainsAny(typed, " \t") {
			return fmt.Sprintf("%q", typed)
		}
		return typed
	default:
		payload, err := json.Marshal(typed)
		if err != nil {
			return fmt.Sprintf("%v", typed)
		}
		return string(payload)
	}
}
