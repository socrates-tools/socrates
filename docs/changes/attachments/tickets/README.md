# Attachments beside the words a person writes

A person answering a page agora serves can attach files beside each decision's words and beside the message, and the caller's agent reads each one by the path `socrates agora poll` prints.
socrates review carries a person's files to the next round's reviewer, and socrates design keeps what a decision rests on beside the spec.

Spec: [../spec.md](../spec.md). Decisions: [ADR 87](../../../../agora/docs/adr/0087-an-attachment-travels-with-the-answer-at-the-press-and-is-never-served-back.md), [ADR 88](../../../../agora/docs/adr/0088-every-collected-answer-is-kept-whole-in-its-session.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1.1: Every collected answer is kept whole in its session](1.1-every-collected-answer-kept-whole.md) | None | done |
| [1.2: The file tray primitive](1.2-the-file-tray-primitive.md) | None | done |
| [1.3: socrates review respond records attachments](1.3-review-respond-attach.md) | None | done |
| [2: Attach by the button, send at the press, poll the paths](2-attach-and-poll-the-paths.md) | 1.1, 1.2 | done |
| [3.1: Drop and paste attach](3.1-drop-and-paste.md) | 2 | done |
| [3.2: An answer is stored whole or refused whole](3.2-answer-ceilings.md) | 2 | done |
| [3.3: Files alone are an answer, and what goes names them](3.3-files-alone-are-an-answer.md) | 2 | in progress |
| [3.4: A settled answer names its files](3.4-settled-answers-name-their-files.md) | 2 | done |
| [3.5: The review and design skills read attachments](3.5-skills-read-attachments.md) | 2, 1.3 | done |
