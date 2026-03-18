## Why

The library currently stops at `error`, which forces callers that need process-terminating logs to handle emission and exit behavior themselves. Adding a first-class `fatal` log level closes that gap now that the public logging API and example CLI are in place and ready to support a complete severity ladder.

## What Changes

- Extend the public logging API with a `Fatal` logging method and a `fatal` level constant so callers can emit terminal events through the same structured interface as other severities.
- Modify logger configuration and level parsing so `fatal` is accepted anywhere the library currently accepts `debug`, `info`, `warn`, or `error`.
- Define fatal emission behavior normatively: the logger must emit the entry, include the same structured fields and metadata as other levels, sync best-effort, and then terminate the process with a non-zero exit code.
- Define fatal idempotency and retry expectations: each `Fatal` call is terminal and MUST NOT retry failed writes because termination begins after the single emission attempt.
- Define fatal failure handling and recovery expectations: if sync or write cleanup fails, the implementation must not suppress the fatal termination path, and callers must treat fatal logging as non-recoverable.
- Define observability expectations so fatal entries remain auditable in both pretty and JSON output and are covered by unit tests and example updates.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `logging-api`: Extend the stable API requirements to include fatal-level configuration and fatal log emission semantics.

## Impact

- Affected code includes public level definitions, logger methods, configuration parsing, tests, README examples, and the sample CLI in `example/`.
- The public API changes in a backward-compatible way by adding new constants and methods without removing existing behavior.
- No new runtime dependencies are expected.
- Unit tests should remain the primary verification method, with targeted coverage for fatal output shape, termination behavior, level parsing, and example usage.
