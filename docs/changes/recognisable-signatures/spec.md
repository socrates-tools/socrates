# A release's signature and provenance carry the names verifiers look for

Decisions: ADR 65.

## Problem Statement

Every release is signed over its checksums with Sigstore and attested with build provenance.
The signature is published as `checksums.txt.bundle`, and the provenance lives only in GitHub's attestation store.
OpenSSF Scorecard reads a release's asset names, recognises neither, and scores every release so far as unsigned and without provenance.

## Solution

The signature is published as `checksums.txt.sigstore.json`.
The release workflow attaches the provenance it attests to the release as `socrates.intoto.jsonl`.
The installer, the README and the release footer verify against the new name, and the installer falls back to the old one for an earlier release.

## User Stories

1. As a person deciding whether to trust socrates, I want Scorecard to see that releases are signed and carry provenance, so that the badge says what is true.
2. As a person verifying a download by hand, I want the signature named for the tool that reads it, so that I know what to run it through.
3. As a person verifying provenance without GitHub's API, I want the provenance attached to the release, so that I can download it beside the archive.
4. As a person installing an earlier release with `SOCRATES_VERSION` and `cosign` installed, I want the installer to still verify it, so that pinning a version does not break the install.
5. As a contributor, I want the suite to fail where the signature or the provenance is published under a name verifiers do not look for, so that the score cannot quietly drop back.

## Implementation Decisions

- `.goreleaser.yaml` names the signature `${artifact}.sigstore.json`.
- `.github/workflows/release.yml` gives the attestation step an id and uploads its `bundle-path` to the tag's release as `socrates.intoto.jsonl`, in the publishing job, which already holds `contents: write`.
- `install.sh` downloads `checksums.txt.sigstore.json` and, where the release has none, `checksums.txt.bundle`.
- `README.md`, the release footer and the release skill name the new file and the provenance.
- `internal/convention` reads the signature's published name and the release workflow's upload, and states the extensions verifiers recognise.

## Testing Decisions

- A convention test fails where the signature over the checksums is not named with `.sigstore.json`, or where no step of the release workflow attaches an `.intoto.jsonl` to the release.
  Each is mutated once.
- The installer's end-to-end tests serve a release whose signature is `checksums.txt.sigstore.json`, and one serving only `checksums.txt.bundle` asserts the fallback verifies with `cosign` on `PATH`.
- `make release-snapshot` shows the new signature name is accepted by goreleaser.
- The next release is checked live: its assets carry both files, and Scorecard's next reading counts it signed.

## Out of Scope

- Renaming or adding assets on releases already published.
