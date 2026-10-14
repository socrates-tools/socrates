# 65. A release names its signature and provenance the way verifiers look for them

- Status: accepted
- Date: 2026-10-14
- Part of: recognisable-signatures

## Decision

The Sigstore bundle over the checksums is published as `checksums.txt.sigstore.json`, the extension Sigstore names for a bundle, rather than `checksums.txt.bundle`.
The build provenance the release workflow attests is also attached to the release, as `socrates.intoto.jsonl`, beside being stored in GitHub's attestation store.

Both were already produced, and neither was seen.
Scorecard reads a release's asset names, and a signature or provenance under a name it does not recognise scores as none, so v0.1.0 to v0.3.0 read as unsigned and without provenance though each was signed and attested.
A person verifying by hand is helped by the same names, because `.sigstore.json` says what tool reads the file and an attached `.intoto.jsonl` can be checked without asking GitHub's API for it.

The installer verifies against `checksums.txt.sigstore.json` and falls back to `checksums.txt.bundle`, so a release named with `SOCRATES_VERSION` from before this still verifies where `cosign` is installed.

## Rejected

- Uploading renamed copies to the releases already published.
  A published release is what somebody already installed, and its assets are not edited by hand, so the score rises as new releases replace old ones in Scorecard's window.
- Leaving the names alone because the files are verifiable either way.
  True, and a release whose signature nobody can find is worth less to the person deciding whether to trust it than one whose signature is where they look.
- A SLSA generator workflow in place of the attestation action, which publishes `.intoto.jsonl` on its own.
  It is a second build path with its own pinned toolchain, to produce a file the existing attestation already holds.
