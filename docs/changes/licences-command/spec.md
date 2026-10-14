# The binary prints the terms it carries

Decisions: ADR 63.

## Problem Statement

A person who installs socrates from a release archive has `LICENSE` and `NOTICE` beside the binary.
A person who installs it with `go install`, or who copies the binary somewhere, has neither, so they hold a binary without the terms it carries.

## Solution

`socrates licences` prints `LICENSE`, then `NOTICE`, exactly as the repository holds them, from copies embedded when the binary was built.
`socrates licenses` is the same command.

## User Stories

1. As a person who installed socrates with `go install`, I want to read the licence socrates is under, so that I know what I may do with it.
2. As a person redistributing the binary, I want to read every borrowed work it carries and the terms of each, so that I can meet them.
3. As a person who types the American spelling, I want `socrates licenses` to work, so that I am not refused over a spelling.
4. As a contributor, I want the suite to fail where what the binary prints differs from the repository's files, so that the two cannot drift.

## Implementation Decisions

- A package at the module root embeds `LICENSE` and `NOTICE` with `go:embed`, because an embed directive can only reach files in its own package's directory.
- `internal/dispatcher` mounts a `licences` command at the top level, with `licenses` as an alias, that prints the two in that order and refuses an argument.
- The convention list of the suite's own top-level commands gains `licences`.
- `README.md`'s Licence section names the command.

## Testing Decisions

- `internal/e2e` runs the built binary and asserts its output is `LICENSE` followed by `NOTICE`, byte for byte against the repository's files, under both spellings, and that an argument is refused.
- The guard is mutated once by printing `NOTICE` alone.

## Out of Scope

- Paging the output. It is printed to standard output for a pager or a file.
- Listing the Go modules the binary links in separately from `NOTICE`, which already names them.
