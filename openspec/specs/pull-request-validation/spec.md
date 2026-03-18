## Purpose

Define the pull request workflow behavior for automatic unit test validation before merge.

## Requirements

### Requirement: Pull requests run unit tests automatically
The repository MUST provide a GitHub Actions workflow that runs the unit test suite for pull requests targeting the main development branch. The workflow MUST execute the repository's Go unit tests using the standard project command for test verification.

#### Scenario: Open or update a pull request
- **WHEN** a contributor opens, reopens, or pushes commits to a pull request targeting the main branch
- **THEN** the pull request validation workflow runs the repository's unit tests

### Requirement: Failed pull request tests are surfaced before merge
If the unit test workflow fails, GitHub MUST report the failure on the pull request so reviewers and maintainers can see that the change is not ready to merge. The workflow MUST exit non-zero on test failures rather than masking them.

#### Scenario: Unit test failure in CI
- **WHEN** the pull request validation workflow encounters a failing unit test
- **THEN** the workflow reports a failed status on the pull request

### Requirement: Pull request validation is safe to rerun
The pull request validation workflow MUST be rerunnable without manual cleanup and MUST produce the same pass-or-fail outcome for the same repository state, aside from nondeterministic test failures. Dependency restoration or caching MUST NOT change the logical outcome of the test run.

#### Scenario: Rerun validation for unchanged code
- **WHEN** a maintainer reruns the pull request validation workflow for the same commit SHA
- **THEN** the workflow performs the same unit test verification without requiring manual reset steps

### Requirement: Workflow logs provide observability for failures
The pull request validation workflow MUST emit job logs that identify the setup and test step that failed. This observability MUST be available in GitHub Actions so maintainers can recover by fixing the code or workflow and rerunning the job.

#### Scenario: Dependency or setup failure
- **WHEN** the pull request validation workflow fails before tests execute
- **THEN** the failing setup step is visible in the workflow logs
