# 39. A change is a folder and a decision is a file

- Status: accepted
- Date: 2026-09-17
- Part of: socrates-design

## Decision

`docs/changes/<slug>/` holds a change's spec.
`docs/adr/NNNN-<sentence>.md` holds the decisions, numbered across the repository, each carrying `Part of: <slug>`.

The folder holds what is only ever about one change, and the tracker holds the tickets, because a ticket has a state, an assignee and a blocking edge that a committed file cannot carry.
A tenth of what a spec produces is not: an ADR outlives the change that produced it, is read by people who never heard of it, and amends ADRs that came from other changes.

Three properties make its home repository-wide rather than per change.
Its number is a repository-wide identifier already, because ADRs cite each other by bare number.
Supersession crosses changes, so a decision overturned a year later is amended by an ADR from a different change, and the two belong in one chronological list.
And the reader of an ADR has an area rather than a slug, which `docs/agents/domain.md` already assumes when it tells every skill to read the ADRs touching the area it is about to work in.

`docs/one-command/` and `docs/monorepo/` move under `docs/changes/`.

The spec lists the ADRs it produced and each ADR names its change, so both links exist and the skill writes both.

## Rejected

- One folder per change holding the spec and the ADRs, so that nothing links across directories.
  The link it saves is two lines a skill writes.
  What it costs is paid by people: `docs/adr/` is named in fourteen places outside itself, including the guide the reviewer reads at a run's base commit and the refusal a skill over its ceiling is given; the bare-number citations between ADRs either become ambiguous or become cross-folder paths, which is the same link multiplied by every ADR rather than by every change; and a change that overturns an earlier decision must write into an earlier change's finished folder.
- `docs/features/<slug>/`, which is what this change was first asked for.
  Two of the three specs this repository holds are a rename and a restructure, and neither is a feature.
- `docs/design/<slug>/`, naming the workflow that produced the folder.
  The folder is named for the change, which is the repository's own word for a unit of work, and it is opened by a design and read by whoever implements it.
