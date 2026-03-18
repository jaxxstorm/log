## Context

The repository now contains a working Go logging library, but it does not yet define any repository automation for pull request validation or tagged releases. That leaves merge quality dependent on local developer discipline and makes releases manual, inconsistent, and difficult to audit. This change introduces automation across GitHub Actions so that validation, release publication, and changelog generation follow a repeatable path.

The main stakeholders are the repository maintainer shipping tagged module releases and contributors opening pull requests. The core constraints are low operational overhead, compatibility with GitHub-hosted workflows, and deterministic release outputs that can be rerun safely without creating duplicate releases.

## Goals / Non-Goals

**Goals:**
- Run the repository's unit test suite automatically for pull requests that target the main branch.
- Publish GitHub releases from version tags using a repeatable workflow with explicit permissions and failure visibility.
- Generate changelog content automatically during release publication so releases include human-readable notes.
- Keep release automation idempotent so rerunning a job for the same tag does not create conflicting releases or divergent notes.

**Non-Goals:**
- Introducing multi-environment deployment automation or package publication beyond GitHub releases.
- Enforcing a full conventional-commit workflow across all contributors.
- Building a complex matrix test strategy beyond the baseline unit test coverage required for merge validation.
- Automating semantic version selection or release branching policy in this change.

## Decisions

### Use GitHub Actions for both pull request validation and release orchestration
The repository will define GitHub Actions workflows because the requested behaviors are all GitHub-native: merge-request checks, tagged releases, and surfaced workflow logs. This keeps repository automation co-located with the code and avoids external CI dependencies.

Alternative considered: relying on a separate CI service. This was rejected because it adds more credentials, duplicated configuration, and weaker integration with GitHub release state.

### Use GitHub CLI release creation with generated notes
The release workflow will use `gh release create` with generated notes rather than a separate release tool. For a Go library repository, this keeps the implementation small while still allowing tag verification, generated changelog content, and deterministic skip behavior when a release already exists for a tag.

Alternative considered: introducing GoReleaser. This was rejected because the repository does not need packaged binaries or a dedicated release assembly layer for library tags.

### Keep pull request validation narrow and fast
The pull request workflow will focus on `go test ./...` with a standard Go toolchain setup and dependency caching. That gives the repository a reliable merge gate without over-specifying a matrix or adding noisy non-blocking checks in the initial automation pass.

Alternative considered: adding linting and multi-version Go matrices immediately. This was rejected for the first change because the user only asked for unit tests on merge requests, and broader CI can be layered in later.

### Publish changelog content from git history scoped to each release
Release notes will be generated from the commits included since the prior release tag so each GitHub release includes a concise changelog. The release workflow will use the same git history as the source for the published notes, making reruns stable for the same tag.

Alternative considered: maintaining a hand-edited `CHANGELOG.md`. This was rejected because it creates release toil and is easy to forget in a small repository.

### Fail visibly and early on missing release prerequisites
The release workflow will require a matching semantic version tag, repository checkout with full history, and an authentication token that can create or update GitHub releases. If those prerequisites are missing, the workflow will fail rather than attempting a partial release.

Alternative considered: silently downgrading to a dry run. This was rejected because it obscures release failures and weakens operator confidence.

## Risks / Trade-offs

- [Generated changelog quality depends on commit history quality] -> Keep the generated notes deterministic now and refine commit hygiene or filtering rules later if needed.
- [Release reruns may behave differently if git history is shallow] -> Configure the release workflow to fetch full history before changelog generation.
- [Generated release notes follow GitHub's note generation heuristics] -> Keep commit history readable and refine the workflow later if maintainers need more control over changelog formatting.
- [PR validation may miss issues outside unit tests] -> Treat this as the baseline merge gate and add linting or additional checks in future changes if needed.
- [Tag-driven releases can still be triggered accidentally] -> Scope release execution to version-like tags and document the expected release trigger clearly.

## Migration Plan

Add the GitHub Actions workflows and release documentation, then validate them locally where possible using repository test execution and workflow review. After merging, enable the workflows in GitHub and create the first version tag to exercise the release path. Rollback is straightforward: disable or remove the workflow files if the automation behaves unexpectedly.

## Open Questions

- Whether the repository should eventually adopt a stricter commit message convention to improve changelog readability.
- Whether future releases should attach additional artifacts beyond the default GitHub release metadata.
- Whether branch protection rules in GitHub should be updated separately to require the pull request validation workflow.
