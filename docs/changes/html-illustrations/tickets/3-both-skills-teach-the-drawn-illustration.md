# 3: Both skills teach the drawn illustration

## What to build

`socrates-design` and `socrates-review` teach the illustration as an HTML file named by path, drawn in a sandboxed frame, with the files it references, the 16 MiB and type bounds, the three CDN hosts, the tokens stylesheet's address, and what is refused.
Both keep the rule to draw only where the choice or the finding is a shape.
Design adds that a prototype's result can be the prototype itself.
Review adds that a finding decided by looking carries a before/after drawing, with screenshots the agent captures by whatever browser it has, and words where it has none.
Every document that says the page reaches nowhere beyond this machine says instead that it sends nothing beyond it, and `SECURITY.md` names the three CDN hosts as what a drawing may load from.
The committed renderings are regenerated for both hosts.

## Acceptance criteria

- [ ] An agent following either skill writes a document naming an illustration by path that agora serves without refusal.
- [ ] The review skill describes a before/after pair for a finding decided by looking, and the fallback to words.
- [ ] No document says the page reaches nowhere beyond this machine, and each that stated it says it sends nothing beyond it.
- [ ] `SECURITY.md` names cdnjs, jsDelivr and unpkg as the hosts a drawing may load from and says a drawing sends nothing.
- [ ] Both skills stay within the ceiling a skill is held to.

## User stories

- 15, 16

## Blocked by

- 2.1-an-illustration-brings-the-files-it-references.md
- 2.4-a-drawing-may-link-the-page-tokens.md
