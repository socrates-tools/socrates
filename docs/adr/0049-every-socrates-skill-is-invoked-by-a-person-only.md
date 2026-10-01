# 49. Every socrates skill is invoked by a person only

- Status: accepted
- Date: 2026-10-01
- Part of: socrates-implement

## Decision

Every skill the binary ships carries `disable-model-invocation: true`: `socrates-review`, `socrates-review-handover`, `socrates-design`, `socrates-tickets` and `socrates-implement`.
Claude Code then lists none of them to a model, blocks a session that tries to invoke one, and tells it not to reproduce the skill another way.
No skill reaches for another.

The person stated it while answering how the overseer reviews a ticket, in their own words: "socrates-review skill must not be auto invoked by any other skill - this has to be invoked by human only. In fact, all socrates-* skills can only be invoked manually by humans. We do this to avoid inter-skill-invocation that will cause chaos."

Two things follow, and both are the surprise.
The overseer reviews every ticket herself and socrates review is not named, though the suite ships an independent reviewer: "Whenever a subagent is done the work, overseer will review the diff, knowing that the diff will be accumulated as more subagents finish their work. No socrates-review is needed in this case."
And a design's last card stops offering to cut the spec's tickets next, which `socrates-tickets` added: the session would have to read a skill that is now a person's only, so the person types `/socrates-tickets` themselves.

## Rejected

- Both: the implementer drives socrates review on its own change first, and the overseer then reviews against the ticket.
  Rejected in the person's words above.
- One socrates review over the whole change once every ticket is done, its findings handed to one implementer.
  Rejected in the person's words above.
- All five person-only, with `socrates-tickets` left invocable by a session as the one exception, so the last card's hand-on still works.
  Declined without a word.
- Only `socrates-implement` marked in this change, leaving the other four and the hand-on for a change of their own.
  Declined without a word.
