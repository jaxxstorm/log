## Purpose

Define how tagged GitHub releases generate and publish changelog content.

## Requirements

### Requirement: Each published release includes generated changelog content
The release process MUST generate changelog content for each published GitHub release and include that content in the release notes visible to users. The generated notes MUST summarize the changes contained in the release rather than leaving the release body empty.

#### Scenario: Publish a release
- **WHEN** the release workflow publishes a GitHub release for a version tag
- **THEN** the resulting GitHub release includes generated changelog content

### Requirement: Changelog generation is deterministic for a given tag
For a given tag and git history, changelog generation MUST produce stable content across reruns of the same release workflow. The workflow MUST use a defined source of release history so repeated runs do not produce conflicting changelog bodies for the same release.

#### Scenario: Rerun release notes generation
- **WHEN** a maintainer reruns the release workflow for an unchanged tag and commit history
- **THEN** the generated changelog content is consistent with the prior run

### Requirement: Changelog scope is limited to release history since the prior release
Generated release notes MUST be derived from the changes included between the current release tag and the previous relevant release boundary, or from repository history when no prior release exists. The workflow MUST avoid including unrelated future changes in the published notes.

#### Scenario: Publish a second release
- **WHEN** the repository publishes a release after at least one prior release exists
- **THEN** the changelog includes changes since the previous release rather than the full repository history

### Requirement: Changelog generation failures block release publication
If changelog generation cannot complete, the release workflow MUST fail instead of publishing a release with missing or partial notes. That failure MUST be visible in GitHub Actions logs so the maintainer can correct the issue and retry.

#### Scenario: Release notes generation fails
- **WHEN** the release workflow cannot generate changelog content from the configured repository history
- **THEN** the workflow fails before publishing the GitHub release
