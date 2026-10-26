# 90. The last read puts socrates review in front of the whole change

- Status: accepted
- Date: 2026-10-26
- Part of: last-read-review
- Amends: ADR 49

## Decision

Once the last ticket has landed and passed, the overseer runs `socrates review` over the whole change, from the default branch's fork point, beside her own reading of the spec's user stories.
Every finding that asks for a fix goes to one fresh implementer with the stories the change falls short of, and the second round judges what landed.
Nobody is asked anything: the person reads the outcome at hand-back.

This reverses one entry of ADR 49's Rejected section, "one socrates review over the whole change once every ticket is done, its findings handed to one implementer", which the person reopened in this design.
The rest of ADR 49 stands.
Every skill stays a person's to start, and the overseer invokes none: she drives review's commands, `begin`, `findings`, `respond` and `continue`, and never the `socrates-review` skill.

What forces it is what the last read could not see.
It read the whole diff against the user stories only, so what shows only across tickets, the same helper written twice by two tickets or one concept spread over five packages, had no reader, because each ticket's review saw half of it.
By the end of a long run the overseer's own context has compacted, and a reviewer that wrote nothing and read nothing else is the reader that gap needs.

The surprise for the next engineer is that implement now names review's commands while still knowing nothing of review's skills, and that the test which held "the overseer names no `socrates review` invocation" is replaced rather than broken.

## Rejected

- Keep ADR 49, and have the last read review the whole diff on the repository's standards as well as the spec, by the overseer's own reading or by a fresh subagent, without socrates review.
  Declined without a word.
- Keep ADR 49 and the last read as they are, and reject the gap.
  Declined without a word.
