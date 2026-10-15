# Codex is a host for every socrates skill

This change lets a person run every socrates skill from Codex as well as from Claude Code.
Each skill is rendered once per host, every install writes both, the implementer becomes a Codex role with settings of its own, and a review started from Codex reviews with `codex` where nobody chose a reviewer.
It also makes Claude Code's install follow `CLAUDE_CONFIG_DIR`.

Spec: [../spec.md](../spec.md). Decisions: [ADR 66](../../../adr/0066-every-skill-is-rendered-once-per-host-and-installed-for-every-host.md), [ADR 67](../../../adr/0067-a-skill-is-person-only-by-each-hosts-own-switch.md), [ADR 68](../../../../implement/docs/adr/0068-a-codex-implementer-has-settings-of-its-own-and-a-root-it-can-write.md), [ADR 69](../../../../review/docs/adr/0069-the-reviewers-vendor-follows-the-host-where-nobody-set-one.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1-skills-rendered-for-a-host](1-skills-rendered-for-a-host.md) | None | done |
| [2.1-codex-renderings-installed](2.1-codex-renderings-installed.md) | 1 | done |
| [2.2-claude-config-dir-honoured](2.2-claude-config-dir-honoured.md) | 1 | done |
| [3.1-codex-escalates-the-page](3.1-codex-escalates-the-page.md) | 2.1 | done |
| [3.2-reviewer-follows-the-host](3.2-reviewer-follows-the-host.md) | 2.1 | done |
| [3.3-codex-implementer-role](3.3-codex-implementer-role.md) | 2.1 | done |
| [3.4-codex-overseer](3.4-codex-overseer.md) | 2.1 | done |
| [4-codex-implementer-settings](4-codex-implementer-settings.md) | 3.3 | done |
| [5.1-readme-names-codex](5.1-readme-names-codex.md) | 2.2, 3.1, 3.2, 3.4, 4 | in progress |
| [5.2-live-runs-under-codex](5.2-live-runs-under-codex.md) | 3.1, 3.2, 3.4, 4 | pending |
