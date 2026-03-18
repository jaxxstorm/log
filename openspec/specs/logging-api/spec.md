## Purpose

Define the stable public logging API for constructing loggers, emitting structured events, and handling invalid configuration.

## Requirements

### Requirement: Logger construction provides deterministic defaults
The library MUST provide a constructor for creating a logger instance without requiring callers to configure every internal Zap detail. When optional configuration is omitted, the constructor MUST apply deterministic defaults for log level, metadata inclusion, and output formatting selection. Repeated construction with the same explicit configuration MUST produce equivalent logger behavior.

#### Scenario: Create logger with defaults
- **WHEN** a caller constructs a logger without overriding optional settings
- **THEN** the library returns a usable logger configured with the library's default level and metadata behavior

#### Scenario: Repeat construction with the same configuration
- **WHEN** a caller constructs two loggers with the same explicit settings
- **THEN** both loggers emit the same shape of log records for the same input

### Requirement: Logger supports leveled structured log emission
The library MUST expose methods for emitting `debug`, `info`, `warn`, and `error` events with a message and optional structured fields. Each emitted entry MUST include the selected severity, the original message text, and any caller-supplied fields that are valid for emission.

#### Scenario: Emit an info log with fields
- **WHEN** a caller logs an `info` message with structured fields
- **THEN** the resulting log entry contains the `info` level, the provided message, and the supplied fields

#### Scenario: Filter messages below the configured level
- **WHEN** a caller emits a message below the logger's configured minimum level
- **THEN** the library suppresses that entry without returning an error

### Requirement: Logger supports contextual child loggers
The library MUST allow callers to derive a child logger that appends structured context to subsequent entries. Fields added to a child logger MUST be included on all entries emitted through that child unless explicitly overridden by per-call fields.

#### Scenario: Add persistent context to a child logger
- **WHEN** a caller derives a child logger with a `request_id` field and logs multiple messages
- **THEN** each emitted entry from that child includes the `request_id` field

#### Scenario: Override child context with per-call fields
- **WHEN** a caller logs with a per-call field that uses the same key as a child logger field
- **THEN** the emitted entry prefers the per-call value for that key

### Requirement: Logger includes baseline event metadata
The library MUST include a timestamp and caller information on every emitted log entry by default. The timestamp and caller metadata MUST be present in both pretty and JSON output modes so operators can audit where and when an event was emitted.

#### Scenario: Emit an error log
- **WHEN** a caller emits an `error` message
- **THEN** the resulting entry includes timestamp and caller metadata in addition to the level and message

### Requirement: Invalid configuration fails during construction
The library MUST validate explicit configuration at construction time and MUST return an error instead of a partially initialized logger when configuration cannot be applied. Log emission methods MUST NOT implement retries because failed construction is the point at which misconfiguration is reported and recovered.

#### Scenario: Reject unsupported configuration
- **WHEN** a caller requests a configuration combination that the library cannot apply
- **THEN** the constructor returns an error and no usable logger instance
