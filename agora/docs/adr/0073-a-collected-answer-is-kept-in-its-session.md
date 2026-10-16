# 73. A collected answer is kept in its session

- Status: accepted
- Date: 2026-10-16
- Part of: issue 76

## Decision

Collecting an answer handed it to the caller and kept no copy anywhere a caller could reach.
A caller that lost what `agora poll` printed, by piping it through `tail` or by its own process ending, had lost it for good: a second poll waits for a new send, and the session's directory held only the document and how to reach the page.
The person had to answer the whole page again, and what a person answered is the one thing in an exchange nothing can regenerate.

Collecting an answer now writes it to `answer.json` in its session's directory before the exchange lets go of it.
A write that fails leaves the answer uncollected and the poll failing, so the next poll collects it, and no answer is handed out that the session did not keep.
The file holds the last answer collected for that session, so a later collection replaces it, and it is removed with the rest of the session's directory when the surface collects the session.

`agora poll --last` prints the answer the session last kept, in the form `agora poll` printed it, and never waits.
It reads the file rather than asking the surface, so it works once the surface has stopped, and it collects nothing, so an answer sent since stays held for the next poll.
A session whose directory holds no answer is refused, naming the poll that waits for one, and a name that is not one path element is refused before anything is read.

This leaves ADR 30 standing: the exchange still holds at most one uncollected answer and a waiting poll still delivers it once.
What is kept is what was collected, and ADR 59 binds the same answer.

## Rejected

- Writing every send as it arrives, rather than what is collected.
  An answer nobody collected was never handed to the caller, so there is nothing for it to have lost, and printing one again would collect it without the page learning that it had been.
- Asking the surface for the last answer it handed out.
  The answer would live only as long as the surface process, and a surface replaced or stopped before the caller noticed the loss would take it along.
- A second poll printing the last answer again whenever nothing new has arrived.
  `agora poll` waiting on a person is the ordinary case, and one that sometimes returns at once with an old answer would hand a caller an earlier round as though it were the reply to this one.
