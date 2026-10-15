# 69. The reviewer's vendor follows the host where nobody set one

- Status: accepted
- Date: 2026-10-15
- Part of: codex-host
- Amends: ADR 60

## Decision

`socrates review begin` takes `--host claude` or `--host codex`, and each rendering of `socrates-review` passes its own host.
Where `reviewer_cli` is unset, the host's vendor is the reviewer CLI: `claude` for Claude Code and `codex` for Codex.
A `reviewer_cli` a person set wins whatever the host, and a `begin` without `--host` keeps `claude` as the default it has today.
The vendor is resolved once at `begin` and recorded on the run, so every later round of it is reviewed by the same CLI.

ADR 60 keeps the reviewer independent of the session that called it, and rejected a flag the skill tells the session to pass because it puts the reviewer's routing under the control of the thing being reviewed.
That still holds for the model and the effort: they come from the vendor's defaults and `roles.reviewer`, and nothing about the session's own model is read, asked for or passed.
What changes is only the vendor, and only where nobody chose one.

The surprise is that a person on Codex gets a Codex reviewer by default, the same vendor as the session that wrote the change, which is the kind of likeness ADR 60 argued against.
What is bought is that socrates review works on a machine with only Codex on it without a setting first, and a person who wants a reviewer of another vendor says so once in `reviewer_cli`.

## Rejected

- The reviewer CLI stays a setting of its own, `claude` by default, and a person on Codex sets `reviewer_cli: codex` once, so the reviewer stays independent of the host.
  Declined without a word, in round two and again in round three.
- `socrates review begin --reviewer-cli`, which only the Codex rendering passes, and only where `socrates review config` reports `reviewer_cli` unset.
  Declined without a word.
