# 66. Every skill is rendered once per host and installed for every host

- Status: accepted
- Date: 2026-10-15
- Part of: codex-host
- Amends: ADR 61

## Decision

socrates supports two hosts, Claude Code and Codex, and every skill the binary ships is rendered once per host from one source.
Each rendering names its own host's directories, its own way of starting a skill, its own way of running a subagent and its own way of asking the person before a command runs.
Each is held to its own ceiling and to the commands the binary offers, as a single rendering is today.

Every install writes for every supported host, whatever is on `PATH`: `socrates skill` and each namespace's `skill` write Claude Code's rendering under `~/.claude` and Codex's under `~/.agents/skills` and `$CODEX_HOME/agents`.
What lands therefore never depends on the machine it was run on, and a person who installs Codex after socrates needs no second install.

`--to` and `--agents-to` now need `--host claude` or `--host codex`, naming whose rendering goes to the directory named, and either flag without `--host` is refused naming `--host`.
Without `--to` every host's default is written.
`socrates skill` still takes no flags, per ADR 61.

The surprise is that a person who never runs Codex finds `~/.agents/skills` and `~/.codex/agents` written by socrates.
What is bought is one install that is the same everywhere, and a text that tells each host exactly what to do rather than describing every host's way at once.

## Rejected

- One host-neutral document, installed for both hosts byte for byte, worded around what differs.
  Declined without a word.
- Writing for each host whose CLI is on `PATH`, so a host installed later needs another install.
  Declined without a word.
- Writing for the hosts a person names, by a setting or a `--host` flag, Claude Code by default.
  Declined without a word.
- A directory flag per host, `--codex-to` and `--codex-agents-to` beside `--to` and `--agents-to`.
  Declined without a word.
- Dropping `--to` and `--agents-to`, so every skill goes to each host's default directory and nowhere else.
  Declined without a word.
