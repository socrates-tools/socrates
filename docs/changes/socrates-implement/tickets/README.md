# socrates: implement, an overseer landing a change's tickets through implementers

`socrates implement` is a fourth namespace of the one binary.
A person invokes the `socrates-implement` skill, and the session that does becomes the overseer: she writes no code, hands each ticket on the frontier to an implementer subagent in a worktree of its own, reviews what comes back, and lands it on a branch named for the change.
The namespace also holds the settings that shape her: how many implementers run at once, their model and effort, and the default branch of each repository.

Spec: [../spec.md](../spec.md). Decisions: [ADR 48](../../../adr/0048-socrates-implement-is-a-fourth-namespace-with-a-file-of-its-own.md), [ADR 49](../../../adr/0049-every-socrates-skill-is-invoked-by-a-person-only.md), [ADR 50](../../../adr/0050-an-implementers-model-and-effort-are-settings-written-into-an-agent-definition.md), [ADR 51](../../../adr/0051-the-default-branch-is-remembered-per-repository-in-socrates-implements-own-file.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1.1: the configure menu is shared machinery](1.1-the-configure-menu-is-shared.md) | None | in progress |
| [1.2: every socrates skill is invoked by a person only](1.2-every-socrates-skill-is-a-persons-only.md) | None | in progress |
| [1.3: socrates implement exists and reports its settings](1.3-socrates-implement-reports-its-settings.md) | None | done |
| [2.1: the default branch is remembered per repository](2.1-the-default-branch-is-remembered.md) | 1.3 | in progress |
| [2.2: the implementer's agent definition is installed](2.2-the-implementers-agent-definition-is-installed.md) | 1.3 | in progress |
| [3.1: the socrates-implement skill teaches the overseer](3.1-the-socrates-implement-skill.md) | 1.2, 2.1, 2.2 | pending |
| [3.2: socrates implement configure changes one setting from a menu](3.2-socrates-implement-configure.md) | 1.1, 2.1, 2.2 | pending |
