# 54. The reviewer runs in safe mode, and what it held is read rather than trusted

- Status: accepted
- Date: 2026-10-08
- Part of: docs-cleanup

## Decision

review starts the claude CLI with `--safe-mode`, which disables every customisation the machine brings on its own account: memory files and auto memory, skills, plugins, hooks, external tool sources, custom commands and agents, output styles, workflows, themes, keybindings, the status line, the file-suggestion commands and LSP servers.
Auth, model selection, built-in tools and permissions keep working, so the reviewer loses the machine's additions and keeps the CLI.

`--tools` alone does not do this.
It is an allowlist over the vendor's built-in set and says nothing about what the machine loads beside it.
review's first live run failed because a configured tool source wrote a `.serena` directory into the snapshot worktree, and a later run found the session holding script execution, file upload and browser control that `--tools` had no opinion about.

What survives the flag is what an administrator's policy configures: a policy hook, a policy status line and policy file-suggestion commands.
A policy hook contributes no tool, but its command is a child of the session whose working directory is the sealed worktree, so a write it makes there is refused and the comparison after the round catches one the seal did not.
Where that command reaches outside the worktree is bounded by nothing review holds, so it is a residual review states rather than closes.

`--safe-mode` arrived in claude 2.1.169, and review keeps no fallback for an older CLI.
An older CLI dies at startup on an unknown option, and its own message is what a person reads, because the floor is stated in `README.md` rather than recognised from the CLI's text.

The flag is a request, so what the session actually held is read from the stream's opening event rather than assumed.
review matches that event by type and subtype wherever it arrives, because a hook can fire before it.
It reads two fields and records the rest as an inventory of the machine:

- `tools`, minus what review asked for, so anything any source contributed lands in one remainder.
- `mcp_servers`, because external sources connect asynchronously and can be named at the opening event before any of their tools arrive.

The two are blind in different places, neither fails a round, and a session that names no inventory records that it named none.
`skills`, `plugins`, `slash_commands` and `agents` name what is installed on the machine whatever the session loaded, so a check that read them would report a customisation on every run of every machine that has one.

## Rejected

- `--strict-mcp-config` with an empty configuration.
  A managed configuration refuses it at startup, deliberately, because it asks to replace a managed set.
  review once sent it and fell back on the refusal, which put a second invocation, a retry and a `the reviewer failed` log line into every round of every run.
  The simulator refuses it, so putting it back fails every test that runs a reviewer.
- An empty `allowedMcpServers` through `--settings`.
  It works, and took a session from 74 externally contributed tools to none, and it was what review used before `--safe-mode` existed.
  It covers tool sources only, leaving hooks and plugins, and it names a vendor feature in review's own vocabulary.
- Detecting an older CLI from its error text and saying something better.
  Matching a vendor's prose is the pattern this repository deleted for being fragile.
- Trusting the flag and reading nothing back.
  A machine that loaded a source anyway would record nothing, which is indistinguishable from a clean session.
