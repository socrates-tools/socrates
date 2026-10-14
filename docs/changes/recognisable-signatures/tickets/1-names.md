# 1: The signature and the provenance are published under names verifiers recognise

## What to build

What the spec describes, across the release configuration, the release workflow, the installer and every document naming the signature.

## Acceptance criteria

- [x] The signature over the checksums is published as `checksums.txt.sigstore.json`, and the suite fails where it is named otherwise.
- [x] The release workflow attaches the attested provenance as `socrates.intoto.jsonl`, and the suite fails where nothing attaches one.
- [x] The installer verifies the new name, and falls back to `checksums.txt.bundle` for an earlier release.
- [x] The README, the release footer and the release skill name the new file.
- [x] The next release carries both, and Scorecard counts it signed.

## User stories

- 1, 2, 3, 4, 5

## Blocked by

None
