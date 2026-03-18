## Context

The current logger wrapper exposes `Debug`, `Info`, `Warn`, and `Error` methods and maps configured levels onto Zap's corresponding severities. The public spec, README, example CLI, and unit tests all describe the supported range as stopping at `error`. Adding `fatal` touches multiple surfaces at once: public constants and methods, level parsing, formatter output, tests, and example documentation.

This repo already centralizes log emission inside `Logger.log`, keeps configuration parsing in `config.go`, and owns the pretty renderer in `internal/pretty`. That structure is sufficient for a fatal-level change, but the implementation needs an explicit decision about process termination so the behavior stays testable and consistent across JSON and pretty output.

## Goals / Non-Goals

**Goals:**
- Extend the public API to support `fatal` as both a configured minimum level and an emission method.
- Ensure fatal entries carry the same structured fields, timestamp, and caller metadata as other levels in both JSON and pretty output.
- Guarantee that fatal logging terminates the process after the entry is emitted and final synchronization is attempted.
- Keep the implementation unit-testable without requiring the main test process to exit.

**Non-Goals:**
- Introducing package-global logging state or a separate high-level crash-reporting API.
- Changing the semantics of existing `debug`, `info`, `warn`, or `error` methods.
- Adding stack traces, panic recovery, or crash dump handling beyond Zap's current behavior.
- Reworking the library around Zap's full option surface.

## Decisions

### Add fatal support to the library-owned level and method surface
`config.go` will gain a `FatalLevel` constant and level resolution will map it to `zapcore.FatalLevel`. `logger.go` will gain a `Fatal(msg string, fields ...Field)` method so callers use the same structured interface as every other severity.

Alternative considered: leaving configuration unchanged and only adding a convenience helper around `os.Exit(1)`. This was rejected because it would create an inconsistent API where fatal behavior exists at call sites but cannot participate in normal level filtering.

### Keep fatal control flow in the wrapper instead of delegating termination entirely to Zap options
The wrapper will own fatal termination by checking or writing a fatal entry through the underlying Zap logger, attempting final sync/cleanup, and then invoking an internal exit hook with a non-zero status. The exit hook will default to `os.Exit` in production and be overridable in tests through unexported construction seams.

Alternative considered: configuring `zap.WithFatalHook(zapcore.WriteThenFatal)` and delegating all fatal behavior to Zap. This would reduce wrapper code, but it makes unit testing harder because the default path exits the process directly and gives the package less control over sync and cleanup expectations.

### Preserve formatter parity by updating both JSON and pretty paths
JSON output will inherit fatal-level serialization from Zap once level parsing and emission are wired correctly. The pretty renderer will add an explicit fatal label/style so interactive output remains readable and severity ordering stays obvious.

Alternative considered: allowing pretty output to fall back to the default uppercase level rendering. This was rejected because fatal events are operationally distinct and should be visibly differentiated from `ERROR`.

### Verify fatal behavior primarily with unit tests that isolate termination
Normal logger tests will expand to cover fatal level parsing and rendered output. Termination behavior will be verified either through an injected internal exit hook or a subprocess-style unit test, keeping `go test` as the primary verification path without introducing external integration dependencies.

Alternative considered: relying only on manual example execution for fatal behavior. This was rejected because the spec calls for auditable, repeatable verification and fatal semantics are too easy to regress silently.

## Risks / Trade-offs

- [Fatal termination path becomes difficult to test safely] -> Keep termination behind an internal hook or other test seam and verify behavior in unit tests.
- [Sync and close behavior differs between formatter paths] -> Route fatal through the same logger core path and explicitly test JSON and pretty output for fatal entries.
- [Callers accidentally configure `fatal` and suppress all lower-severity logs] -> Document the semantics in README and example help text so the filtering behavior is unsurprising.
- [Implementation duplicates parts of Zap fatal handling] -> Keep the wrapper logic narrow and reuse Zap's core level handling rather than introducing a parallel severity system.

## Migration Plan

This is an additive API change, so no data or config migration is required for existing callers. Implementation can land in four steps: add the level constant and parsing support, add the logger method and internal termination hook, update pretty/example/documentation surfaces, and then expand unit coverage. Rollback is straightforward: revert the added constant, method, renderer label, and docs if fatal behavior proves undesirable before release.

## Open Questions

- Whether the internal test seam for process termination should live on `Logger`, `resolvedConfig`, or as a package-level variable hidden from the public API.
- Whether the example CLI should expose `fatal` only in help text or also demonstrate the call in the example flow despite its terminating behavior.
