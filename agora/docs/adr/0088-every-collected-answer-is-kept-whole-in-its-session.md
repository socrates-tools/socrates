# 88. Every collected answer is kept whole in its session

- Status: accepted
- Date: 2026-10-26
- Part of: attachments

## Decision

ADR 73 kept the last collected answer in `answer.json`, and a later collection replaced it.
Words could be replaced safely, because a caller that needs an earlier round's words carries them forward itself, as `socrates design` does in its settled list.
Attachments change that: a design runs many rounds in one session, and the screenshot a person attached in round 1 is what the spec has to cite once the last round is answered.

Every collected answer is now kept whole in its session's directory, as `answers/<n>/` holding its `answer.json` beside its files, for the life of the session.
A path `socrates agora poll` printed stays readable until the surface collects the session, and each round's words stay beside that round's attachments.
`socrates agora poll --last` still prints the last answer collected, without waiting and without the surface.

Nothing new removes them.
They go with the rest of the session's directory when the surface collects the session, once no page has been open on it for the idle period, so a caller writing its spec after `socrates agora end` still finds them.

In the person's words: "we do keep the records of the questions of answers of each round, right? so naturally we should keep the attachments? Once the spec and ADRs are written, we can remove those attachments."

## Rejected

- Keeping only the last collected answer's files, with each caller copying what it keeps right after the poll.
  It matched ADR 73, but a caller that copies late would find round 1's files gone, and the record of each round would live nowhere agora keeps it.
- Keeping every collected answer's files while `answer.json` stays the last answer only.
  An earlier round's attachments would stay while the words written beside them were replaced, so a file and its words would part.
- A new `socrates agora` command the design skill runs after writing the spec, to remove the attachments at once.
  The session already goes when the surface collects it, and agora's command tree would gain a fifth verb for nothing the idle period does not do.
