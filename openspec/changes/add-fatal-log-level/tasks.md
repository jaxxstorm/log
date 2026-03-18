## 1. Public API And Level Parsing

- [x] 1.1 Add `FatalLevel` to the public level constants and update configuration parsing so `fatal` maps to Zap's fatal severity.
- [x] 1.2 Add a `Fatal(msg string, fields ...Field)` method on `Logger` that emits structured fatal entries through the existing wrapper API.
- [x] 1.3 Update any public-facing strings and comments that enumerate supported levels so they include `fatal`.

## 2. Fatal Execution Path

- [x] 2.1 Implement the internal fatal termination flow so fatal logs emit once, attempt final sync/cleanup, and terminate with a non-zero exit status.
- [x] 2.2 Add any internal test seam needed to verify fatal termination behavior without exiting the main test process.
- [x] 2.3 Update the pretty formatter so fatal entries render with an explicit fatal label consistent with the new level support.

## 3. Verification And Examples

- [x] 3.1 Extend unit tests to cover fatal level parsing, fatal-level filtering, and fatal output shape in structured logs.
- [x] 3.2 Add unit coverage for fatal termination behavior, including the case where final sync or cleanup fails.
- [x] 3.3 Update the example program and README to document fatal support and verify the example/help text reflects the new level list.
