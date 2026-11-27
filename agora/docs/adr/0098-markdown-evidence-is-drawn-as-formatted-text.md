# 98. Markdown evidence is drawn as formatted text

- Status: accepted
- Date: 2026-11-27
- Amends: ADR 78

## Decision

A block of evidence whose language is Markdown, by any name the page highlights it under, is drawn as formatted text rather than as numbered source.
Its headings read as headings, its task list as boxes ticked where the file ticks them, its tables as tables and its fences as highlighted code.
Its heading names the caption and the language as any block's does, and its copy control puts the Markdown on the clipboard exactly as the document stated it.

`ui-markdown-block` draws it.
It reads the text with marked's lexer and draws each token with Lit, so nothing the caller wrote is ever parsed as markup.
HTML written in the Markdown is shown as the words it is, an image as its alt text, and a link as a link only where it names a web page or a mail address.
A relative link, such as a ticket naming `../spec.md`, would lead nowhere from the page, so it is drawn as its words.

socrates-tickets shows every proposed ticket as one Markdown block of the whole file it would be written as.
Drawn as source, a ticket read as a wall of wrapped monospace lines with `##` and `- [ ]` in them, which is not how a person reads a ticket anywhere else.

## Rejected

- Keeping every block as code and changing the tickets skill to show a ticket as an illustration.
  The skill would write HTML for every ticket, the page would frame it in a sandbox, and the copy control would copy nothing a person could commit.
- Rendering marked's HTML into the page.
  It needs a sanitiser beside marked, and the page's policy embeds nothing a caller wrote, which drawing tokens with Lit keeps true by construction.
- A control flipping a block between the formatted text and its source.
  The copy control already hands a person the source, and the toggle is one more state to draw for what is read once.
- A new field on a block for prose.
  ADR 78 settled that evidence is one list of blocks, and a language already says what the code is.
