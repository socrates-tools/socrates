# 2: Attach by the button, send at the press, poll the paths

## What to build

Every decision's words and the composer's message carry a tray, and either press sends one multipart POST holding the answer and every file the trays hold, per ADR 87.
The surface writes every file into a fresh folder before it swaps its uncollected answer, names each by its position and a sanitised base name, removes the folder of the answer it replaces, and removes every uncollected folder on a revision.
`socrates agora poll` and `poll --last` print an `attachments` list of name, type, size and absolute path beside a decision's words and beside the message, omitted where empty.
The surface serves no attachment back, the page's policy admits `blob:` images for the previews, a revision empties every tray, and a closed exchange draws no tray.
`agora/AGENTS.md` gains an invariant for each.

## Acceptance criteria

- [ ] A file attached by the button to a decision and one to the message, sent with a press, appear in what `socrates agora poll` prints with absolute paths whose bytes match what was attached.
- [ ] A poll of an answer with no attachments prints exactly what it prints today.
- [ ] Amending before the poll sends every file again, and the replaced answer's folder is gone from disk.
- [ ] A revision empties every tray, and no uncollected answer's folder is left on disk.
- [ ] A closed exchange draws no tray.
- [ ] No route of the surface returns the bytes of an attachment.
- [ ] A stated name holding a path separator or `..` is stored inside the answer's folder.

## User stories

- 1, 15, 17, 18, 20, 23

## Blocked by

- 1.1-every-collected-answer-kept-whole.md
- 1.2-the-file-tray-primitive.md
