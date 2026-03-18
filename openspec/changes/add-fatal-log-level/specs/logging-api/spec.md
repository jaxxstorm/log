## ADDED Requirements

### Requirement: Fatal logging terminates after emission
The library MUST treat `fatal` logging as a terminal operation. A `Fatal` call MUST attempt to emit exactly one structured fatal entry with the same message, caller metadata, timestamp, and valid structured fields used by other levels. After the emission attempt, the library MUST perform best-effort final synchronization for the configured writer and MUST terminate the process with a non-zero exit status. The library MUST NOT retry a failed fatal write or continue normal execution after the fatal path begins.

#### Scenario: Emit a fatal log
- **WHEN** a caller emits a `fatal` message with structured fields
- **THEN** the library writes a fatal entry containing the `fatal` level, the provided message, and the supplied fields before terminating with a non-zero exit status

#### Scenario: Fatal finalization fails
- **WHEN** a fatal log write or final sync encounters an error
- **THEN** the library does not attempt a second emission and still terminates the process with a non-zero exit status

## MODIFIED Requirements

### Requirement: Logger supports leveled structured log emission
The library MUST expose methods for emitting `debug`, `info`, `warn`, `error`, and `fatal` events with a message and optional structured fields. Each emitted entry MUST include the selected severity, the original message text, and any caller-supplied fields that are valid for emission. The library MUST accept `fatal` anywhere it accepts a configured minimum level so callers can filter lower-severity messages and reserve output for terminal failures.

#### Scenario: Emit an info log with fields
- **WHEN** a caller logs an `info` message with structured fields
- **THEN** the resulting log entry contains the `info` level, the provided message, and the supplied fields

#### Scenario: Filter messages below the configured level
- **WHEN** a caller emits a message below the logger's configured minimum level
- **THEN** the library suppresses that entry without returning an error

#### Scenario: Configure fatal as the minimum level
- **WHEN** a caller constructs a logger with `fatal` as the configured minimum level
- **THEN** lower-severity messages are suppressed without returning an error and fatal messages remain eligible for emission
