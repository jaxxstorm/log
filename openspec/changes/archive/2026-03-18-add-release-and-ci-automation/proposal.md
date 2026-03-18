## Why

This repository can now build and test locally, but it still depends on manual release steps and has no merge-time guardrails in CI. Adding release and validation automation now reduces regressions, makes published versions reproducible, and ensures every release carries a usable changelog.

## What Changes

- Define a GitHub-based release workflow that creates tagged releases for the module and publishes generated release notes directly from GitHub Actions.
- Define a pull request validation workflow that runs unit tests automatically for merge requests before changes are merged.
- Define changelog generation behavior so each release includes a generated summary of user-facing changes.
- Define workflow retry and idempotency expectations so rerunning CI does not publish duplicate releases or inconsistent changelog content.
- Define failure handling and observability expectations for release and validation workflows, including surfaced workflow failures in GitHub Actions.

## Capabilities

### New Capabilities
- `github-release-automation`: GitHub workflow behavior for tagged releases, release artifact generation, and duplicate-release protection.
- `pull-request-validation`: GitHub workflow behavior for running unit tests on pull requests and surfacing failures before merge.
- `release-changelog-generation`: Automated changelog generation and attachment to GitHub releases.

### Modified Capabilities
- None.

## Impact

Affected areas include GitHub Actions workflow definitions, repository documentation, and tag-driven release publication behavior. This change adds CI automation dependencies and establishes the operational contract for release cadence, merge validation, and release note observability.
