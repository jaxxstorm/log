## ADDED Requirements

### Requirement: Output format defaults to terminal-aware selection
The library MUST choose pretty human-readable output when writing to a TTY and MUST choose JSON output when writing to a non-TTY destination, unless the caller explicitly overrides the format. Format selection MUST be resolved during logger construction so the decision is stable for the lifetime of the logger.

#### Scenario: Write to an interactive terminal
- **WHEN** a logger is constructed against a TTY-backed output without an explicit format override
- **THEN** the logger emits pretty formatted log entries

#### Scenario: Write to redirected output
- **WHEN** a logger is constructed against a non-TTY output without an explicit format override
- **THEN** the logger emits JSON log entries

### Requirement: Callers can override automatic format selection
The library MUST allow callers to explicitly request pretty or JSON formatting, and that override MUST take precedence over automatic TTY detection. The same override MUST be honored across repeated constructions using the same configuration.

#### Scenario: Force JSON on a terminal
- **WHEN** a caller explicitly requests JSON formatting for a TTY-backed output
- **THEN** the logger emits JSON entries instead of pretty terminal output

#### Scenario: Force pretty formatting in tests
- **WHEN** a caller explicitly requests pretty formatting for a non-TTY output
- **THEN** the logger emits pretty formatted entries instead of JSON

### Requirement: Output modes preserve canonical event data
Pretty and JSON rendering MUST preserve the same canonical event information, including level, message, timestamp, caller, and structured fields, even if presentation differs. Pretty formatting MAY reorder or style fields for readability, but it MUST NOT drop required event data.

#### Scenario: Compare pretty and JSON output
- **WHEN** two loggers with the same event input emit one entry in pretty mode and one in JSON mode
- **THEN** both outputs contain the same event data even though the encoded representation differs

### Requirement: Output initialization failures are surfaced to callers
If the logger cannot initialize its configured writer or encoder, construction MUST fail with an error. The library MUST NOT silently fall back to a different destination or format because that would hide operational misconfiguration.

#### Scenario: Output destination cannot be initialized
- **WHEN** a caller provides an output target that cannot be opened or configured
- **THEN** the constructor returns an error describing the initialization failure
