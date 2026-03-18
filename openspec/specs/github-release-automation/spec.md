## Purpose

Define the GitHub Actions behavior for publishing tagged releases and preventing duplicate release records.

## Requirements

### Requirement: Version tags publish GitHub releases
The repository MUST provide a GitHub Actions workflow that publishes a GitHub release when a version tag matching the repository's release pattern is pushed. The workflow MUST check out the repository history required for release generation and MUST publish the release against the triggering tag.

#### Scenario: Push a version tag
- **WHEN** a maintainer pushes a version tag that matches the configured release pattern
- **THEN** the release workflow runs and publishes a GitHub release for that tag

### Requirement: Release publication is idempotent for a tag
The release workflow MUST behave deterministically for a given tag and MUST NOT create duplicate GitHub releases when the workflow is rerun. Rerunning the workflow for the same tag MUST either update the existing release consistently or exit without creating a second release record.

#### Scenario: Rerun a release job
- **WHEN** a maintainer reruns the release workflow for a tag that already has a GitHub release
- **THEN** the workflow does not create a duplicate release for that tag

### Requirement: Release workflow fails visibly on missing prerequisites
The release workflow MUST fail when required prerequisites are unavailable, including release credentials, required git history, or invalid tag context. The workflow MUST surface that failure in GitHub Actions logs so the maintainer can diagnose and recover before retrying.

#### Scenario: Missing release credential
- **WHEN** the release workflow runs without the permissions or token required to create a GitHub release
- **THEN** the workflow fails and GitHub Actions shows the release failure in the job output

### Requirement: Release workflow verifies repository state before publishing
Before publishing a release, the workflow MUST run the repository's automated verification needed to ensure the tag points to a buildable state. If verification fails, the workflow MUST stop before publishing a GitHub release.

#### Scenario: Release verification fails
- **WHEN** the release workflow encounters a failing repository verification step
- **THEN** no GitHub release is published for that run
