# 77. An illustration may load from three CDNs and sends nothing out

- Status: accepted
- Date: 2026-10-16
- Part of: html-illustrations

## Decision

agora promised that the page reaches nowhere beyond this machine, and the page's policy named nothing remote.
A drawing that may run script, per ADR 76, mostly wants a library: Mermaid, D3, a chart library, a CSS framework.
A drawing may also hold private code or screenshots of an unreleased screen.

The frame's policy admits scripts, styles and fonts from `https://cdnjs.cloudflare.com`, `https://cdn.jsdelivr.net` and `https://unpkg.com`, beside the files the illustration brought with it.
It names nothing remote for `connect-src`, `img-src`, `media-src`, `frame-src` and `form-action`, and the sandbox grants no popups and no top-level navigation.
Libraries can come in, and what the drawing holds has no way out.

The promise changes from "reaches nowhere beyond this machine" to "sends nothing beyond this machine": opening a page that carries an illustration may fetch a library from one of those three hosts, and nothing the drawing holds is sent to them or anywhere else.
Every document that states the promise says the new one.

## Rejected

- Nothing remote at all, so everything a drawing uses travels with it.
  The promise would stay word for word, but an agent would inline a library such as Mermaid or D3 into every drawing that uses one.
- No network rule, as lavish-axi has.
  A drawing's script could send what it holds anywhere.
