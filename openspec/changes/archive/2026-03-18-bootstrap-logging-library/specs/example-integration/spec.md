## ADDED Requirements

### Requirement: Repository includes a compiled example program
The repository MUST include an example program that builds against the public logging API and demonstrates intended usage from a Go application entrypoint. The example MUST live in version-controlled source so it can be exercised in local development and CI.

#### Scenario: Build the example program
- **WHEN** a maintainer runs the repository's standard Go build or test workflow
- **THEN** the example program compiles successfully against the current public API

### Requirement: Example demonstrates structured logging workflow
The example MUST show construction of a logger, emission of at least two severity levels, and the addition of structured context to log entries. The example MUST use the same public APIs that downstream consumers are expected to call.

#### Scenario: Read the example source
- **WHEN** a developer inspects the example program
- **THEN** they can see how to create a logger, attach fields, and emit structured log messages

### Requirement: Example documents default formatting behavior
The example MUST explain or demonstrate that terminal output defaults to pretty formatting and non-terminal output defaults to JSON formatting. This documentation MAY live in source comments or nearby README content, but it MUST remain adjacent to the example implementation.

#### Scenario: Run the example in different environments
- **WHEN** a developer runs the example once in an interactive terminal and once with output redirected
- **THEN** the example makes the difference in default output behavior discoverable
