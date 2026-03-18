package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type syncBuffer struct {
	bytes.Buffer
	syncErr error
}

func (b *syncBuffer) Sync() error {
	return b.syncErr
}

func TestNewDefaultsToJSONForNonTTYOutput(t *testing.T) {
	var output bytes.Buffer

	logger, err := New(Config{Output: &output})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	t.Cleanup(func() {
		_ = logger.Close()
	})

	logger.Info("hello", String("service", "api"))

	record := decodeRecord(t, output.Bytes())
	if got := record["level"]; got != "info" {
		t.Fatalf("level = %v, want info", got)
	}
	if got := record["msg"]; got != "hello" {
		t.Fatalf("msg = %v, want hello", got)
	}
	if got := record["service"]; got != "api" {
		t.Fatalf("service = %v, want api", got)
	}
	if record["ts"] == nil {
		t.Fatal("expected timestamp in record")
	}
	if record["caller"] == nil {
		t.Fatal("expected caller in record")
	}
}

func TestLevelFilteringSuppressesLowerPriorityMessages(t *testing.T) {
	var output bytes.Buffer

	logger, err := New(Config{
		Output: &output,
		Level:  WarnLevel,
	})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	t.Cleanup(func() {
		_ = logger.Close()
	})

	logger.Info("ignore me")
	if output.Len() != 0 {
		t.Fatalf("expected info log to be filtered, got %q", output.String())
	}

	logger.Warn("keep me")
	if output.Len() == 0 {
		t.Fatal("expected warn log to be emitted")
	}
}

func TestWithFieldsAllowsPerCallOverride(t *testing.T) {
	var output bytes.Buffer

	logger, err := New(Config{Output: &output})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	t.Cleanup(func() {
		_ = logger.Close()
	})

	child := logger.With(String("request_id", "initial"), String("component", "worker"))
	child.Info("override", String("request_id", "override"), Int("attempt", 2))

	record := decodeRecord(t, output.Bytes())
	if got := record["request_id"]; got != "override" {
		t.Fatalf("request_id = %v, want override", got)
	}
	if got := record["component"]; got != "worker" {
		t.Fatalf("component = %v, want worker", got)
	}
	if got := record["attempt"]; got != float64(2) {
		t.Fatalf("attempt = %v, want 2", got)
	}
}

func TestFormatOverrideSupportsPrettyAndJSON(t *testing.T) {
	var prettyOutput bytes.Buffer
	prettyLogger, err := New(Config{
		Output: &prettyOutput,
		Format: PrettyFormat,
	})
	if err != nil {
		t.Fatalf("new pretty logger: %v", err)
	}
	t.Cleanup(func() {
		_ = prettyLogger.Close()
	})

	prettyLogger.Info("started", String("component", "cli"))
	prettyLine := strings.TrimSpace(prettyOutput.String())
	for _, want := range []string{"INFO", "started", "component=cli"} {
		if !strings.Contains(prettyLine, want) {
			t.Fatalf("pretty output %q missing %q", prettyLine, want)
		}
	}

	var jsonOutput bytes.Buffer
	jsonLogger, err := New(Config{
		Output: &jsonOutput,
		Format: JSONFormat,
	})
	if err != nil {
		t.Fatalf("new json logger: %v", err)
	}
	t.Cleanup(func() {
		_ = jsonLogger.Close()
	})

	jsonLogger.Info("started", String("component", "cli"))
	record := decodeRecord(t, jsonOutput.Bytes())
	if got := record["component"]; got != "cli" {
		t.Fatalf("component = %v, want cli", got)
	}
	if got := record["msg"]; got != "started" {
		t.Fatalf("msg = %v, want started", got)
	}
}

func TestInvalidConfigReturnsError(t *testing.T) {
	if _, err := New(Config{Level: Level("trace")}); err == nil {
		t.Fatal("expected invalid level to fail")
	}

	if _, err := New(Config{Format: Format("console")}); err == nil {
		t.Fatal("expected invalid format to fail")
	}

	if _, err := New(Config{
		Output:     &bytes.Buffer{},
		OutputPath: "test.log",
	}); err == nil {
		t.Fatal("expected conflicting output settings to fail")
	}
}

func TestOutputInitializationFailureReturnsError(t *testing.T) {
	_, err := New(Config{OutputPath: "/tmp/does-not-exist/log/output.json"})
	if err == nil {
		t.Fatal("expected invalid output path to fail")
	}
}

func TestFatalLevelFiltersLowerSeverityMessages(t *testing.T) {
	var output syncBuffer

	logger, err := New(Config{
		Output: &output,
		Format: JSONFormat,
		Level:  FatalLevel,
	})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}

	logger.Error("ignore me")
	if output.Len() != 0 {
		t.Fatalf("expected error log to be filtered, got %q", output.String())
	}
}

func TestFatalWritesStructuredJSONAndExits(t *testing.T) {
	var output syncBuffer

	logger, err := New(Config{
		Output: &output,
		Format: JSONFormat,
		Level:  FatalLevel,
	})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}

	exitCode := 0
	logger.exitFn = func(code int) {
		exitCode = code
	}

	logger.Fatal("fatal failure", String("component", "cli"))

	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}

	record := decodeRecord(t, output.Bytes())
	if got := record["level"]; got != "fatal" {
		t.Fatalf("level = %v, want fatal", got)
	}
	if got := record["msg"]; got != "fatal failure" {
		t.Fatalf("msg = %v, want fatal failure", got)
	}
	if got := record["component"]; got != "cli" {
		t.Fatalf("component = %v, want cli", got)
	}
	if record["ts"] == nil {
		t.Fatal("expected timestamp in record")
	}
	if record["caller"] == nil {
		t.Fatal("expected caller in record")
	}
}

func TestFatalPrettyOutputUsesExplicitLabel(t *testing.T) {
	var output syncBuffer

	logger, err := New(Config{
		Output: &output,
		Format: PrettyFormat,
	})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}

	exitCode := 0
	logger.exitFn = func(code int) {
		exitCode = code
	}

	logger.Fatal("fatal failure", String("component", "cli"))

	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}

	line := strings.TrimSpace(output.String())
	for _, want := range []string{"FATAL", "fatal failure", "component=cli"} {
		if !strings.Contains(line, want) {
			t.Fatalf("pretty output %q missing %q", line, want)
		}
	}
}

func TestFatalStillExitsWhenSyncOrCloseFails(t *testing.T) {
	output := &syncBuffer{syncErr: errors.New("sync failed")}

	logger, err := New(Config{
		Output: output,
		Format: JSONFormat,
	})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}

	exitCode := 0
	closeCalled := false
	logger.exitFn = func(code int) {
		exitCode = code
	}
	logger.closeFn = func() error {
		closeCalled = true
		return errors.New("close failed")
	}

	logger.Fatal("fatal failure")

	if !closeCalled {
		t.Fatal("expected fatal path to call close function")
	}
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}

	record := decodeRecord(t, output.Bytes())
	if got := record["level"]; got != "fatal" {
		t.Fatalf("level = %v, want fatal", got)
	}
}

func decodeRecord(t *testing.T, payload []byte) map[string]any {
	t.Helper()

	lines := bytes.Split(bytes.TrimSpace(payload), []byte("\n"))
	if len(lines) == 0 || len(lines[0]) == 0 {
		t.Fatal("expected at least one log line")
	}

	var record map[string]any
	if err := json.Unmarshal(lines[0], &record); err != nil {
		t.Fatalf("decode log line %q: %v", lines[0], err)
	}

	return record
}
