# Retrospective: a-review-nobody-authored

## The live run

Task 13.3 asks for one review of a real change in a real repository with the real `claude` CLI, and for what it found to be written down.
The change reviewed was socrates' own group 12: the scenario coverage gate, `AGENTS.md`, `README.md`, `docs/gotchas.md`, and the tool server restriction that the first attempt at this live run made necessary.
The diff was 38 KB across eleven files, under the inlining ceiling, so the reviewer was handed the whole thing.
It ran for about twenty minutes at `opus` and `high`, read the repository around the diff, and returned nine findings and a `high` risk rating.

The suite was green before the run: `make fmt`, `make lint`, `make test`, `make scenario-coverage` and the three fidelity targets all passed.
That is the ratio worth recording.
**A green suite against a simulator, plus a coverage gate with no gaps, plus a drift check against recordings of the real CLI, and one live review still found two defects that hold the change and five more worth fixing.**

## What socrates itself got wrong, discovered by trying to run it

These are not findings the reviewer reported.
They are what happened to the run.

### The reviewer's tool servers wrote in the snapshot worktree

The first live run failed with `the reviewer wrote to the snapshot worktree, which it may not do: .serena/.gitignore, .serena/project.yml`.
The reviewer never asked to write anything.
The machine's configured tool servers are loaded by the CLI regardless of `--tools Read,Grep,Glob`, and one of them writes a project directory into its working directory the moment it starts.

The second mechanism did its job: the comparison of the worktree against the snapshot commit caught it, named the paths and kept the worktree.
The first mechanism was simply incomplete, because a tool restriction says nothing about a second source of tools.
The specification gained a requirement, socrates now asks for `--mcp-config {"mcpServers":{}} --strict-mcp-config`, and `TestNoToolServerIsLoaded` holds it.

### The CLI refuses to ignore an enterprise tool server configuration

The second live run failed differently: `claude exited: exit status 1: You cannot use --strict-mcp-config when an enterprise MCP config is present`.
This devcontainer has `/etc/claude-code/managed-mcp.json`, which is outside the person's control, and the CLI will not honour strictness while it exists.

socrates now asks for strictness, and on that specific refusal runs the review again without it, recording `agent-tool-servers-loaded` in the run's event log.
That is a real weakening of the first mechanism, stated in the specification rather than hidden: where the CLI will not let the servers be ignored, the worktree comparison is the only thing holding the reviewer read-only.
On this machine the consequence is that every review still fails, because the server named in that enterprise configuration writes into the worktree.
The remedy is to take that server out of the enterprise configuration, and not to weaken the comparison, which is the one mechanism that holds whatever the vendor's flags are called next year.

### The simulator had invented a field

The drift check found this before the live run, on its first execution: the simulator emitted `structured_output` in its result event even when no schema was supplied, which the real CLI does not do.
Every test that used the simulator supplied a schema, so no test in the suite could have found it.
This is the whole argument for the recording and the drift check, demonstrated within an hour of them existing.

## What the reviewer reported

Nine findings, one `high` risk rating, and the rationale it gave: two defects hold the change, and the security boundary the change exists to establish is dropped on a substring match against arbitrary CLI stderr.

| id | severity | what it found | what happened |
| --- | --- | --- | --- |
| f1 | error | The coverage gate does not pass on this repository, while tasks 12.2 and 13.2 are ticked claiming it does: the new scenario `A machine that will not let them be ignored` is named by no test | Already fixed between the snapshot and the reviewer's answer, by the `// scenario:` marker the gate asked for. The reviewer was right about the snapshot it was given |
| f2 | error | The retry breaks cancellation: `noteAgentProcess` latched on the first `agent-started` event, so `run.Agent` held the dead first process and `scrat cancel` killed nothing, marked the run `cancelled`, and left a live model process with its working directory deleted underneath it | Fixed. The run now records whichever process is answering, and `TestCancellingStopsEverythingAfterASecondInvocation` fails on the old code |
| f3 | warning | Dropping the tool server restriction was decided by `strings.Contains(err.Error(), "strict-mcp-config")` against up to 8 KiB of the CLI's stderr, so an unrelated failure mentioning the flag would silently load every server on the machine | Fixed. The condition is now the CLI's refusal phrase and an attempt that produced nothing at all |
| f4 | warning | `Gaps` reported a scenario as covered when its text normalised to empty, because `strings.Contains(claim, "")` is true, so a malformed `#### Scenario:` heading was silently counted as tested | Fixed. An empty normalised scenario is an unconditional gap, with `TestAScenarioWithNoTextIsAGap` |
| f5 | warning | The documentation the change was written to produce omits the invariant the same change introduced: neither `AGENTS.md` nor `README.md` mentioned tool servers | Fixed in both |
| f6 | warning | The fallback test built the agent with no recorder, so the event the specification requires was emitted into nothing and asserted by nobody | Fixed. The test now records and asserts the event, and the end-to-end test asserts it in the run's own log |
| f7 | info | `attempt` built a fresh deadline per call, so a review that took the retry path was bounded by twice the configured deadline | Fixed. One deadline now spans both attempts |
| f8 | info | The round's transcript is not reset between attempts, so a first attempt that wrote to stdout before failing would leave two concatenated streams in `stream.jsonl` | Accepted. The transcript is an `io.Writer` that socrates cannot truncate, the parse that matters is the second attempt's, and the cost is legibility of a file rather than a wrong answer |
| f9 | info | `AGENTS.md` claimed every invariant is held by a test that fails, while scenario coverage was only a separate `make` target that `make test` never invoked, which is how f1 got past two ticked verification tasks | Fixed. The gate is now `TestEveryScenarioUnderWayIsNamedByATest` inside the ordinary suite |

Seven fixed, one already fixed, one accepted with the reasoning written down.
No finding was wrong.

## What this says about the design

- **The two-mechanism rule earned its keep immediately.** The vendor's tool restriction was incomplete in a way nobody predicted, and the worktree comparison caught it on the first real invocation. Had socrates trusted the flags, the reviewer would have been writing in the snapshot and nothing would have said so.
- **`failed` not being `clean` earned its keep too.** Both failing live runs produced a reviewer that answered perfectly well; the round still failed, and no verdict claimed the change was reviewed. The findings the reviewer returned were kept on disk and are what this retrospective is written from.
- **The reviewer used its independence.** f2 and f9 are not local defects: they are consequences of the change that only appear when you read the code the diff does not show, which is exactly what the "trace a counterexample" and "require a concrete path" rules ask for. f1 and f9 are criticisms of the author's own bookkeeping, which is what an author reviewing their own work does not write down.
- **A simulator can only be as honest as its drift check.** One field of fiction existed for as long as there was no check, and the check found it the first time it ran.
