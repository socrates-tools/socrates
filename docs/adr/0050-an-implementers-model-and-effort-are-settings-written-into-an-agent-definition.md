# 50. An implementer's model and effort are settings written into an agent definition

- Status: accepted
- Date: 2026-10-01
- Part of: socrates-implement

## Decision

`roles.implementer.model` and `roles.implementer.effort` are settings of socrates implement, `sonnet` and `high` by default.
They reach an implementer through an agent definition, `socrates-implementer.md` in Claude Code's agents directory, which `socrates implement skill` writes from the settings in force and the overseer spawns.
`socrates implement configure` writes that definition again, wherever it was installed, as soon as a person confirms a new model or effort.

The definition is forced by Claude Code: a spawn can name a model, but only a definition carries an effort.
The surprise is that `configure` writes outside the brand's root, into a file of the harness, and that the binary therefore records where it last installed the definition, in `~/.socrates/implement/installed.yaml`, because the configuration file holds only what a person set, per ADR 28.
What is bought is that a setting `config` reports is the setting the next implementer runs on, with no second step to forget.

## Rejected

- The skill command installs a definition fixed at `sonnet` and `high`, and nothing configures it.
  Declined without a word.
- The overseer names `sonnet` per spawn, and the effort is inherited from her session.
  Declined without a word.
- A changed model reaches the implementers on the next `socrates implement skill`, and `configure` says so.
  Declined without a word.
- The overseer compares `socrates implement config` with the installed definition at the start of every run, and says in the terminal where they disagree.
  Declined without a word.
