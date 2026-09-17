---
name: socrates-design
description: Interview a person over a page rather than a scrollback, and leave a spec and its decisions behind. Use when a change needs a spec before it has code.
---

# socrates design

A design is an interview.
It settles what a change is for and leaves behind a record of the decisions, rather than a change to the code.

What it leaves behind:

- `docs/changes/<slug>/spec.md`, the spec.
- `docs/adr/NNNN-<sentence>.md`, one per decision that is hard to reverse, surprising without context, and the result of a real trade-off.
- `CONTEXT.md`, maintained as terms settle.

The interview is put to a person on a page, the same way socrates review hands a finding to one, rather than as prose in a terminal:

    socrates design skill          install this skill
    socrates agora serve           put a round of questions in front of a person
    socrates agora poll            read what they answered

## Installing elsewhere

This document installs under `~/.claude/skills` by default.
Write it somewhere else instead with `socrates design skill --to <directory>`.
