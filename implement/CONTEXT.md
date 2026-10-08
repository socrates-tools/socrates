# Context: implement

`socrates implement` lands a change's tickets: the session that invokes it hands each ticket to a subagent and reviews what comes back.
These are implement's words, on top of the suite's shared ones in `/CONTEXT.md`.
A ticket, a blocker, a wave, a main ticket, abandoned and needs a person are design's words in `/design/CONTEXT.md`, and implement reads the tickets design wrote.
`docs/adr/` here holds implement's decisions and `AGENTS.md` its invariants.

## The implementation

- **overseer**: the session that invokes `socrates-implement`. It writes no code: it hands tickets to implementers, reviews what each one hands back, and is the one channel through which the person talks about the change while it is built.
- **implementer**: a subagent the overseer hands one ticket to. It writes the code for that ticket, and patches it when the overseer's review sends feedback back to it.
- **frontier**: every ticket whose blockers are done. It is what the overseer may hand out.
- **notes**: what one subagent writes about the code a change touches before the first hand-out, kept outside the repository under `~/.socrates/implement/notes/` and pointed at by every brief. Never committed.
- **default branch**: the branch a person starts a feature from in one repository, which socrates asks for once and remembers. Not git's or a forge's default branch, which may be a different one: in a repository whose forge shows `main`, a person who starts every feature from `dev` has `dev` as the default branch.
