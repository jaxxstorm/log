package log_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	loglib "github.com/jaxxstorm/log"
)

func TestNewDefaultsToJSONForNonTTYOutput(t *testing.T) {
	var output bytes.Buffer

	logger, err := loglib.New(loglib.Config{Output: &output})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	t.Cleanup(func() {
		_ = logger.Close()
	})

	logger.Info("hello", loglib.String("service", "api"))

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

	logger, err := loglib.New(loglib.Config{
		Output: &output,
		Level:  loglib.WarnLevel,
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

	logger, err := loglib.New(loglib.Config{Output: &output})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	t.Cleanup(func() {
		_ = logger.Close()
	})

	child := logger.With(loglib.String("request_id", "initial"), loglib.String("component", "worker"))
	child.Info("override", loglib.String("request_id", "override"), loglib.Int("attempt", 2))

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
	prettyLogger, err := loglib.New(loglib.Config{
		Output: &prettyOutput,
		Format: loglib.PrettyFormat,
	})
	if err != nil {
		t.Fatalf("new pretty logger: %v", err)
	}
	t.Cleanup(func() {
		_ = prettyLogger.Close()
	})

	prettyLogger.Info("started", loglib.String("component", "cli"))
	prettyLine := strings.TrimSpace(prettyOutput.String())
	for _, want := range []string{"INFO", "started", "component=cli"} {
		if !strings.Contains(prettyLine, want) {
			t.Fatalf("pretty output %q missing %q", prettyLine, want)
		}
	}

	var jsonOutput bytes.Buffer
	jsonLogger, err := loglib.New(loglib.Config{
		Output: &jsonOutput,
		Format: loglib.JSONFormat,
	})
	if err != nil {
		t.Fatalf("new json logger: %v", err)
	}
	t.Cleanup(func() {
		_ = jsonLogger.Close()
	})

	jsonLogger.Info("started", loglib.String("component", "cli"))
	record := decodeRecord(t, jsonOutput.Bytes())
	if got := record["component"]; got != "cli" {
		t.Fatalf("component = %v, want cli", got)
	}
	if got := record["msg"]; got != "started" {
		t.Fatalf("msg = %v, want started", got)
	}
}

func TestInvalidConfigReturnsError(t *testing.T) {
	if _, err := loglib.New(loglib.Config{Level: loglib.Level("trace")}); err == nil {
		t.Fatal("expected invalid level to fail")
	}

	if _, err := loglib.New(loglib.Config{Format: loglib.Format("console")}); err == nil {
		t.Fatal("expected invalid format to fail")
	}

	if _, err := loglib.New(loglib.Config{
		Output:     &bytes.Buffer{},
		OutputPath: "test.log",
	}); err == nil {
		t.Fatal("expected conflicting output settings to fail")
	}
}

func TestOutputInitializationFailureReturnsError(t *testing.T) {
	_, err := loglib.New(loglib.Config{OutputPath: "/tmp/does-not-exist/log/output.json"})
	if err == nil {
		t.Fatal("expected invalid output path to fail")
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
