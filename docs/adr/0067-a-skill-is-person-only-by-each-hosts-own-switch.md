# 67. A skill is person-only by each host's own switch

- Status: accepted
- Date: 2026-10-15
- Part of: codex-host
- Amends: ADR 49

## Decision

Every skill stays a person's to start, per ADR 49, and each host's rendering says so with that host's own switch.
Claude Code's rendering carries `disable-model-invocation: true` in its front matter.
Codex's rendering carries no such field, because Codex reads only `name` and `description` from `SKILL.md` and ignores every other key, and is installed with `agents/openai.yaml` beside it holding `policy.allow_implicit_invocation: false`.

The surprise is that the two switches are not equally strong.
Claude Code hides the skill from the model and blocks a session that tries to invoke it.
Codex only keeps the skill out of the model's context, and a session that types `$socrates-tickets` still reaches it.
Every skill already tells the session to invoke no skill, and that sentence is what holds the rest of ADR 49 on Codex.

The suite's check over every committed skill and everything a `skill` operation installs reads each host's switch rather than one front matter field.

## Rejected

- Shipping nothing to Codex until it can refuse a session's explicit invocation, because the guarantee ADR 49 bought is the point.
  Declined without a word.
