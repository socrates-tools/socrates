# 87. An attachment travels with the answer at the press and is never served back

- Status: accepted
- Date: 2026-10-26
- Part of: attachments

## Decision

A person can only answer a card in words, and some answers are a picture: a screenshot of what a page looks like, a trace, a file the caller never saw.
Every decision's words and the composer's message now carry a tray of up to ten attachments of any type.

The bytes travel at the press and only then, in the same request as the answer.
The page holds each file in memory and previews it from a `blob:` URL, and either press sends one `multipart/form-data` POST holding the answer as JSON and every file every tray holds.
What the controls hold at the submit is still exactly what is sent, so nothing is committed per decision, an amend sends every file again, and no file reaches the surface that no answer names.

The surface stores an answer whole or refuses it whole.
It writes every file into a fresh folder before the exchange swaps its uncollected answer, and refuses a file over 20 MiB, an answer over 100 MiB, an eleventh file in a tray or a name the answer did not state, keeping nothing.
The ceilings are stated once in `agora/surface` beside `BodyCeiling`, which keeps guarding every other route at 1 MiB.
A stored file is named by agora, as its position and a sanitised base name inside the answer's folder, and the person's name for it travels as `name`.

The surface never serves an attachment.
No route reads back what a person sent, so a person's bytes are never rendered from agora's origin, and `socrates agora poll` hands the caller an absolute path for each instead, which the caller's agent reads as it reads any file.

## Rejected

These were set aside on `docs/changes/attachments/design.html`, the page this change was designed from.

- Uploading each file when it is attached.
  It needs an upload route, leaves files no answer names when a person removes one or never presses, and commits something per decision.
- Inlining each file as base64 in what `socrates agora poll` prints.
  The caller's agent would read megabytes of encoded bytes into its context rather than reading a file by its path.
- Serving previews from the surface.
  It would add a route that renders bytes a person or a caller supplied from agora's own origin, which the page's preview from memory never needs.
