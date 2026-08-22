## Context

The change before this one made the tool-source protection correct.
This one makes it small, by noticing two things.

The property being protected was never about one vendor feature: it is that the reviewer holds the tools socrates asked for and nothing else.
And the vendor already ships a flag that disables every source of the ones socrates did not ask for that the machine brings on its own account.

## Goals / Non-Goals

**Goals:**

- One statement of the reviewer's tool surface, checked against what the session reported.
- No vendor feature named in a requirement, an event, or a rendered line.
- A reviewer equipped to check a claim rather than guess at one.

**Non-Goals:**

- Giving the reviewer the tools that write. Patching is the author's half of the asymmetry the whole design rests on.
- Deciding whether an unexpected tool can write. Rejected in the change before this one for the reason that still holds: a list of tool names that write goes stale silently, and it fails permissively.
- Refusing to review where an unexpected tool appears. Unchanged: the review runs and the run records it.

## Decisions

### `--safe-mode` replaces the settings allowlist

The vendor documents it as disabling the customisations the machine brings on its own account, and lists them: memory files and auto memory, skills, plugins, hooks, external tool sources, custom commands and agents, output styles, workflows, themes, keybindings, the status line, the file-suggestion commands and LSP servers.
What it keeps in force is the managed settings policy, which is the carve-out recorded below.
It also documents what keeps working: auth, model selection, built-in tools, permissions.

Measured on a machine that deploys a managed configuration, which is the case that defeated everything before it:

| invocation | external sources named | tools |
| --- | --- | --- |
| neither flag | four | 27 |
| `--safe-mode` | none | 26, the plugin-contributed one gone |
| `--safe-mode` with socrates' own flags and a schema | none | `Glob`, `Grep`, `Read`, `StructuredOutput`, `WebFetch`, `WebSearch` |

It is not refused by a managed configuration, so the refusal that produced a second invocation, a retry and a `the reviewer failed` line on every round has nothing left to attach to.

It covers more than the allowlist it replaces. A hook and a plugin are tool sources the allowlist said nothing about, and the requirement it replaces admitted as much: the tool restriction "does not cover a tool server, a hook, or any other process the vendor's session starts". The plugin tool leaving the session's list is that admission being answered rather than restated.

It also takes the user's own memory and skills out of the reviewer, which is the same direction socrates already went by refusing the reviewed repository's instructions. The `scrat` skill a person installs at user level can no longer seat the reviewer that reads their code.

What it does not cover, measured after the fact and stated rather than closed: a hook an administrator's policy configures still runs, along with the status line and the file-suggestion commands, because the vendor keeps the managed settings policy in force by design. A policy hook contributes no tool, so the tool-list check cannot see it, and outside the snapshot worktree nothing socrates holds bounds it. A policy-configured tool source is not in that residual, which is what the table above shows: the four external sources are gone with the flag and so is the plugin's tool. Managed plugins, managed skills and managed memory go the same way by the vendor's documentation rather than by anything measured here, because the tool list does not say what contributed a tool.

### The tool list is the check

socrates reads `tools` from the stream's opening event, subtracts the tools it asked for, and records the remainder.

The expected set is derived from the invocation rather than written twice: the `--tools` value socrates passes, split, plus the one tool the vendor adds when a schema is supplied.

One honest limit, measured rather than assumed: external sources connect asynchronously, so a probe on this machine named four of them in the opening event while their tools were not in the list yet.
The tool list is therefore what the reviewer can use at the point socrates reads it, not a proof that nothing arrives later.
What covers later is `--safe-mode` not loading the sources at all, and the comparison after the round catching a write from anything that got in regardless.
Reading the source list instead would swap one incomplete signal for another, and would keep the vocabulary this change removes.

The failure mode taken deliberately: a vendor release that adds a built-in outside the `--tools` allowlist makes every run record an unexpected tool until somebody looks.
That is one line, it is true when it appears, and it is the direction to be wrong in.

### The reviewer gets the web and not the pen

`WebSearch` and `WebFetch` join the three reading tools.
A reviewer that can check what a library actually does writes a finding about the code; one guessing from memory writes a finding about its own recollection.

`Write`, `Edit`, `Bash` and `NotebookEdit` stay out.
The first two are patching, which is the author's half.
`Bash` is the one that reaches outside the snapshot worktree, where neither the seal nor the comparison bounds anything, so admitting it would turn "socrates writes nothing in your repository" from a proven property into a request.
`Task` and `SlashCommand` stay out because they start further model processes, and one package owns that.

`WebFetch` widens what a reviewer can send outward. The code already reaches the model provider, so this widens exposure rather than creating it, and it is recorded here rather than discovered later.

### The fence is the simulator refusing the flags

The simulator's argument table no longer lists `--mcp-config` or `--strict-mcp-config`, so an invocation carrying either fails with `error: unknown option`, exit 1.

That is a harder fence than the test it replaces, which asserted their absence: a test can be deleted along with the code it guards, and a rejected argument fails every test that runs a reviewer.

### A scenario rename is expressible

The drop gate compared scenario names and could not tell a deliberate rename from a loss, so removing the vocabulary from five scenario names read to it as dropping five scenarios.

`RENAMED Requirements` already carried `FROM`/`TO` pairs for requirements. The same section now takes them for scenarios, and a promoted scenario counts as carried when the delta declares its new name.
A scenario renamed without declaring it is still reported, which is the case the gate exists for, and both directions are held by a test.

### What a run records, and what a person reads

One event and one sentence, naming the tools.
The sentence says what is true and stops: these tools were in the reviewer's session, socrates did not ask for them, and nothing socrates holds bounds what a tool reaches outside the snapshot worktree.

The event kind is renamed, so an event recorded before this change renders through the fallback that names a kind and its fields.
Retention is 336h, so the window is two weeks and the rendering is legible rather than wrong.
Carrying a second rendering for a kind nothing writes any more would keep the vocabulary this change exists to remove.

## Risks / Trade-offs

- **One flag now carries the whole protection.** → It is documented, it is asserted on the delivered invocation, and the session's tool list is read rather than trusted, so a release that changed its meaning shows up as recorded unexpected tools rather than as silence.
- **`--safe-mode` is described in the vendor's help as useful for troubleshooting.** → What it does is specified and what it preserves is specified, and both were measured here. If a release narrows it, the tool-list check is what notices.
- **The reviewer loses the person's own memory and skills.** → Deliberate, and it points the same way as refusing the reviewed repository's instructions. The cost is the same cost already accepted there: it does not know local conventions, and the prompt's stay-in-lane rule covers most of what that produces.
- **A vendor release that adds a built-in makes every run record it.** → Visible rather than silent, one line, and the weekly fidelity check is where a change to the tool list surfaces first.

## Migration Plan

No stored run changes shape.
A run in the retention window that recorded the old event kind renders through the unrecognised-event fallback, which prints the kind and its fields.

Rollback is the previous change's code, and with it the vocabulary and the refused flag.

## Open Questions

None. The subject is one flag, and this is the last change about it.
