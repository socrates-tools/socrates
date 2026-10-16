# 76. An illustration is an HTML document drawn in a sandboxed frame

- Status: accepted
- Date: 2026-10-16
- Part of: html-illustrations

## Decision

An illustration was a fragment of markup, reduced by `agora/sanitise` to a fixed allowlist and painted by one document stylesheet that `agora/stylesheet` parsed, reduced and scoped, inline on the page.
That kept a caller's drawing safe on the page itself, and it kept out everything a drawing most often needs: script, images, style attributes and any reference but `url(#name)`.
Issue 77 could not show two screenshots, and nothing like a chart, a Mermaid diagram or a toggle between a before and an after could be drawn at all.

An illustration is now a whole HTML document, drawn in an iframe sandboxed with `allow-scripts` and without `allow-same-origin`, served from a route of the session.
The frame runs in an opaque origin, so its script may run but cannot read the page, its token or its storage, and cannot navigate the page.
The same sandbox is stated again as a `Content-Security-Policy` header on every response of that route, so a drawing's file opened on its own is sandboxed as well.
A drawing only shows: it has no channel into the answer, and the page reads one thing from it, its height, as a finite number posted by a script agora adds.

`agora/sanitise` and `agora/stylesheet` are removed, the document's `stylesheet` and the decision's `fragment` go, and both skills teach one way to draw.

## Rejected

- Keeping both kinds, the allowlisted fragment for small inline drawings and a second field for a full document.
  Two ways to draw is two pipelines to keep, and the frame draws everything the fragment could.
- Serving the drawing from a second loopback port, as an origin of its own.
  It is real origin separation, but agora would bind two ports, against ADR 57's one port.
- Inlining the HTML under a shadow root.
  A shadow root keeps styles apart but not scripts, so scripts would stay forbidden and an allowlist would stay.
- Forbidding script inside the frame, so a drawing is HTML, CSS and SVG only.
  It is safer to reason about, but most of what an agent draws would stay out of reach, which is what the change is for.
