# 1: Skills are rendered for a host

## What to build

The skill machinery and every namespace's skills learn that a skill is rendered for a host, with Claude Code the only host so far.
A skill's document, its front matter, its person-only switch and the directory it installs to are each decided by the host it is rendered for.
Nothing a person sees changes: every skill installed and every committed copy is what it is today, byte for byte.

## Acceptance criteria

- [ ] The supported hosts are stated once, in the machinery no tool owns, each with its skills directory, how a person starts a skill and its person-only switch, and Claude Code is the only one stated.
- [ ] Every namespace's skills are rendered from one source per skill, given the host.
- [ ] `make skill` leaves every committed skill and the committed agent definition unchanged.
- [ ] `socrates skill` and each namespace's `skill` install exactly the files they install today, with the same contents.
- [ ] Every existing skill check (ceiling, named commands, prohibitions, committed copy, person-only) runs per host rather than over one rendering.

## Blocked by

None
