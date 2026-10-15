# 68. A Codex implementer has settings of its own and a root it can write

- Status: accepted
- Date: 2026-10-15
- Part of: codex-host
- Amends: ADR 50

## Decision

On Codex, the implementer is a role file, `$CODEX_HOME/agents/socrates-implementer.toml`, carrying `name`, `description`, `developer_instructions` and, where set, `model` and `model_reasoning_effort`.
The overseer spawns it by naming it as `agent_type`, because Codex applies a role only to a spawn that names it.

Its model and effort are settings of their own, `roles.implementer.codex.model` and `roles.implementer.codex.effort`, beside Claude Code's `roles.implementer.model` and `roles.implementer.effort`.
With no Codex model set the role file names none, so the implementer runs on the Codex session's own model, and the Codex effort is `high` by default.
`socrates implement configure` offers both rows and rewrites the Codex role file the moment one is confirmed, as ADR 50 has it for Claude Code's definition, and the record of where each was installed holds one directory per host.

A Codex subagent inherits its parent's sandbox and working directory and cannot be given one of its own, and Codex's default sandbox cannot write `~/.socrates/implement`, where the worktrees and the notes live.
So the Codex overseer checks at the start of every run that it can write there, and where it cannot it starts nothing and tells the person the one line that grants it: `writable_roots` under `[sandbox_workspace_write]` in Codex's configuration, or `--add-dir` when starting Codex.
The README says the same.

The surprise is twofold: the same setting has two keys, one per host, and a run on Codex can refuse to start over a sandbox it does not own.
What is bought is that each host runs on values it accepts, and that a run never fails halfway through at its first `git worktree add`.

## Rejected

- The Codex role file carries instructions only, and a Codex implementer inherits the session's model and effort.
  Declined without a word.
- One model and one effort read by both hosts, written as typed into both files.
  Declined without a word.
- The README states the writable root and the overseer does not check, so a run without it fails at its first worktree.
  Declined without a word.
- Codex worktrees inside the repository's own workspace, which the sandbox already lets Codex write.
  Declined without a word.
