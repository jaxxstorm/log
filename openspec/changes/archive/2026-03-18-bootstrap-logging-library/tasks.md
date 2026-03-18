## 1. Package Skeleton

- [x] 1.1 Create the initial package layout for the logging library, including the public package, internal formatting support, and an example entrypoint.
- [x] 1.2 Add and pin the initial dependencies for Zap, Lip Gloss, and any TTY-detection support required by the design.
- [x] 1.3 Define the exported logger type, field helpers, configuration types, and constructor signatures that make up the initial public API.

## 2. Core Logger Implementation

- [x] 2.1 Implement logger construction with deterministic defaults for level, metadata inclusion, and format selection.
- [x] 2.2 Implement `debug`, `info`, `warn`, and `error` logging methods that emit structured fields through the underlying Zap logger.
- [x] 2.3 Implement child logger creation so persistent contextual fields are attached to subsequent log entries.
- [x] 2.4 Validate explicit configuration during construction and return errors for unsupported or unbuildable logger setups.

## 3. Output Formatting

- [x] 3.1 Implement terminal detection during logger construction and choose pretty output for TTY destinations and JSON for non-TTY destinations by default.
- [x] 3.2 Implement explicit format override handling so callers can force pretty or JSON output regardless of TTY state.
- [x] 3.3 Ensure both output modes preserve the same canonical event data, including level, message, timestamp, caller, and structured fields.

## 4. Example And Verification

- [x] 4.1 Add a version-controlled example program that builds against the public API and demonstrates logger construction, multiple levels, and contextual fields.
- [x] 4.2 Document the default terminal-versus-redirected formatting behavior adjacent to the example implementation.
- [x] 4.3 Add unit tests covering logger construction defaults, level filtering, child logger field inheritance, format override behavior, and configuration failure cases.
- [x] 4.4 Add verification that the example compiles as part of the standard Go test or build workflow.
