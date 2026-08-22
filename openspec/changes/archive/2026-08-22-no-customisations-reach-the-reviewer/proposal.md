## Why

socrates named a vendor feature in its invariants, and then reasoned about that feature instead of about the property it wanted.

The requirement said the reviewer loads no external tool server.
Holding it took an argument, then a second argument, then a refusal, then a retry, then an event, then a rendered sentence, then a fidelity guard, then a paragraph of specification and four of documentation.
None of it was about reviewing.

An external tool source was never the thing socrates cares about.
What socrates cares about is that the reviewer holds the tools socrates asked for and nothing else.
An external tool source is one way that stops being true.
A hook is another, a plugin is another, a memory file that seats the reviewer is a third, and a built-in the vendor adds in a release nobody read is a fourth.
Naming one of them left the rest uncovered and put the vendor's vocabulary through five packages.

The vendor already has a flag for all of it.
`--safe-mode` disables the customisations the machine brings on its own account and leaves auth, model selection and the built-in tools working, which is exactly the reviewer socrates wants.
What it leaves in force is a hook an administrator's policy configures, which contributes no tool and is a residual socrates states rather than covers.
Measured on a machine that deploys a managed configuration: the opening event named four external sources without it and none with it, and a plugin-contributed tool left the session's tool list at the same time.
Unlike the flags that ask to replace a managed set, it is not refused where one exists.

The other half is that the reviewer was starved.
`--tools Read,Grep,Glob` gave it no way to check a library's actual behaviour, and a reviewer guessing an API from memory writes a finding about nothing.

## What Changes

- **One flag replaces the tool-source machinery.** Every invocation carries `--safe-mode`, and the settings allowlist, the vendor's feature name, and every mention of tool servers leave the program with it.
- **The tool list is the check.** socrates reads the session's tool list from the stream's opening event, subtracts the tools it asked for, and records the remainder. Anything a customisation contributes lands in that list, so one check covers every source rather than the one it was written for.
- **The reviewer gets the web.** `WebSearch` and `WebFetch` join `Read`, `Grep` and `Glob`. The tools that write stay out, because patching is the author's half of the asymmetry and the reviewer's independence is the other half.
- **The fence is now the simulator.** It no longer accepts `--mcp-config` or `--strict-mcp-config` at all, so putting either back fails every test that runs a reviewer instead of quietly reintroducing the fallback they needed.
- **The fidelity guard pins the tool list**, rather than a field whose presence socrates' own correct behaviour prevents.
- **A deliberate scenario rename is expressible.** The drop gate could not tell a renamed scenario from a lost one, which made removing the vocabulary from a scenario name impossible without the gate reporting a loss.

## Capabilities

### Modified Capabilities

- `agent-seam`: what the reviewer's tool surface is required to be, how the machine's customisations are kept out of it, and that what the session held is read rather than assumed.
- `reviewer-run`: the three mechanisms, with the tool surface named as what covers a tool that reaches past the worktree.

## Impact

- `internal/agent`: one flag replaces the settings allowlist, the tool list is read and checked, and the tool-server parsing and its event are gone.
- `internal/presentation`: one sentence replaces two, naming tools rather than servers.
- `internal/claudesim`: extra tools replace loaded servers, the allowlist modelling goes, and the two refused flags are no longer accepted arguments.
- `internal/fidelity`: the guard pins the opening event's tool list.
- `internal/promotion`: the drop gate follows a scenario rename declared in the delta, and still reports one that is not.
- `internal/convention`: the event kind the renderer owns is renamed.
- No stored run changes shape. A run inside the 336h retention window that recorded the old event kind renders through the fallback that names a kind and its fields.
