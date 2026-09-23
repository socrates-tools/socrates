# 42. Each tool reads a configuration file of its own

- Status: accepted
- Date: 2026-09-23
- Part of: monorepo

## Decision

Each tool reads one optional configuration file, in its own directory under the brand's root.
socrates review reads `~/.socrates/review/config.yaml`, and socrates agora reads `~/.socrates/agora/config.yaml`.
`SOCRATES_HOME` moves the root, and both files move with it.

The `surface` block leaves review's file with both its fields, `bind` and `reached_at`, and is agora's file's only block.
The block keeps its name and its fields, so moving it is a cut and a paste, and every document naming `surface.bind` or `surface.reached_at` stays true.
`socrates review config` and `socrates review configure` report and change review's settings only.

A review file still carrying a `surface` block is refused by name, naming agora's file as where the block now belongs.
Every review command refuses it before a snapshot is taken and before a reviewer is spawned, except `configure`, which opens it and offers to reset the block as a whole.
Nothing moves the block into agora's file.

This supersedes ADR 24.
Its window moved out of the server only so that review's configuration could report agora's default bind without linking the network into review.
Review no longer reports that setting, so the window is stated in `agora/surface` beside the guards that judge an address, and `internal/loopback` is gone.

## Rejected

- One file for the suite, with a block per tool.
  It keeps ADR 24's problem, one tool's settings read through another tool's package, and multiplies it by the number of tools.
- Tolerating and ignoring the block for one release.
  socrates review does not quietly run on defaults while a configuration file says otherwise, and an ignored block is exactly that for the person who set it.
- Moving the block into agora's file automatically.
  It writes a file the person owns, and the migration stays in the tool for as long as any machine might still hold the old file, which nothing can tell.
- agora refusing a `surface` block it finds in review's file.
  agora would then read review's file, which is the coupling this decision removes.
  The residual is stated: a person who runs only `socrates design` or `socrates agora` never runs a review command, so nothing refuses a block left in review's file, and the page serves on defaults.
  Only an unreleased build could have written the block there, because v0.1.0 read `~/.socrates/config.yaml`, which nothing reads now.
- Flattening the block to top-level `bind` and `reached_at` in agora's file.
  A person would have to rewrite the block rather than move it, and the key every document and message names would change.
