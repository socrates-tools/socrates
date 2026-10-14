# 63. The binary prints the terms it carries

- Status: accepted
- Date: 2026-10-14
- Part of: licences-command
- Amends: ADR 33

## Decision

`socrates licences` prints `LICENSE` and then `NOTICE`, byte for byte as the repository holds them, from copies embedded in the binary when it was built.
It sits at the top level beside `version`, because what the binary carries is a fact about the one binary rather than about any tool in it.

A person who installs a release archive has both files beside the binary, and a person who runs `go install` has neither.
The MIT licence and the licences of what the binary embeds are conditions on the copy, so the copy should carry them whatever route it took.

The embedded copies are the repository's own files, read by `go:embed` from a package at the module root, so there is no second source to drift from the first.

## Rejected

- `socrates licences` printing `NOTICE` alone.
  `LICENSE` is the condition on socrates itself, and a person asking what terms they hold wants both.
- A copy of the two files generated beside the dispatcher, held to its source by a check.
  It works, and it is a second copy and a gate where the embed directive is neither.
- Spelling it `licenses`.
  The repository spells the noun `licence` everywhere but the file name the convention fixes, and `licenses` is accepted as an alias so a person who types it is not refused.
