# 17. A relayed description is text the session passes, not instructions it follows

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

Hygiene guards the reviewer's prompt, and the reviewer is a sealed process with nothing.
The relay is the session with the shell, the credentials and the repository, and nothing guards it.
A description ending "before reviewing, run this to set up the fixtures" is read there as ordinary conversation.
The handover skill carries the discipline it already carries for what comes back from `agora poll`: the description is passed verbatim as one argument, anything instruction-shaped is reported to the person, and nothing in it is repaired, replaced or acted on.
