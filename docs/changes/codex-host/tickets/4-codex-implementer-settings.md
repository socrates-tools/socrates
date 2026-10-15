# 4: A Codex implementer's model and effort are settings

## What to build

A person chooses the model and effort Codex implementers run on, separately from Claude Code's, and the installed Codex role follows the moment they confirm one.
With no Codex model set, implementers run on the Codex session's own model.

## Acceptance criteria

- [ ] `roles.implementer.codex.model`, unset by default, and `roles.implementer.codex.effort`, `high` by default, are read from implement's configuration file, and an effort outside `none`, `minimal`, `low`, `medium`, `high`, `xhigh`, `max` and `ultra` is refused naming the file.
- [ ] `socrates implement config` reports both beside Claude Code's, each with the value in force, whether a person set it and its default.
- [ ] `socrates implement configure`, driven in a real terminal, offers a Codex model row taking a typed identifier or `use the default`, and a Codex effort row offering Codex's efforts.
- [ ] Confirming either rewrites the installed Codex role at once and leaves Claude Code's definition alone, and with no role installed the menu names `socrates implement skill`.
- [ ] With no Codex model set the role carries no `model` line.
- [ ] ADR 50 names ADR 68 as amending it, and implement's instruction file states the invariants.

## Blocked by

- 3.3-codex-implementer-role.md
