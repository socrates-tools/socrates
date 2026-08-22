## Why

socrates asks the CLI to ignore the machine's tool servers with `--strict-mcp-config`.
Where an enterprise MCP configuration is deployed, the CLI refuses that flag at startup by design, because the flag asks to replace a managed set.
socrates then runs the review without it and records that the servers were loaded, which is what the specification says to do.

What nobody had measured is what the reviewer holds in that case.
The `init` event of socrates' own live review names 56 tools, 52 of them from tool servers, and a throwaway probe on the same machine reports 77 and 74.
Among them are `mcp__playwright__browser_run_code_unsafe`, `mcp__playwright__browser_file_upload` and `mcp__chrome-devtools__evaluate_script`.
`--tools Read,Grep,Glob` does not exclude them, because that flag restricts the built-in set and nothing else, and `--dangerously-skip-permissions` means nothing would have asked.

So the reviewer of every run on such a machine can execute arbitrary script and drive a browser.
The seal covers the snapshot worktree and the comparison catches what defeats it, and neither covers anything outside the worktree, which is where those tools reach.
The requirement already says a loaded tool server "is a reviewer that can write, both to the snapshot worktree and outside it", and the fallback was accepted on the belief that the seal answered it. It does not.

There is a documented mechanism the CLI does not refuse.
`allowedMcpServers` may be set in any settings file, entries merge from every source, an empty array means no server is allowed, and the check runs before a server loads, including one from a managed configuration.
Supplying it with `--settings` on the same machine takes the reviewer's session from 74 tool-server tools to none.

## What Changes

- **Every invocation carries an empty tool-server allowlist.** socrates supplies `allowedMcpServers: []` through `--settings`, which no machine refuses.
- **`--strict-mcp-config` and its empty `--mcp-config` go away.** That request is what a managed configuration rejects at startup, and it was the only reason socrates ever invoked the CLI twice. With no server permitted, reading those configuration files achieves nothing, so one mechanism that is never refused replaces two of which one always was. The retry, the second invocation and the event that reported the refusal all go with it.
- **What loaded is read rather than inferred.** socrates takes the tool servers and the session's tool list from the stream's `init` event and records them, instead of deciding from the text of an error message what the CLI did.
- **The line that read `the reviewer failed` on every round of every run is gone,** because there is no rejected attempt left to report.
- **The policy does not change.** Where a tool server loads anyway, the review runs and the run records that it did, because a machine that locks its allowlist down should still get reviews.

## Capabilities

### Modified Capabilities

- `agent-seam`: how the implementation stops tool servers loading, that it no longer asks to replace a managed set, and that it reads what loaded from the `init` event and records it.
- `reviewer-run`: that the three mechanisms hold the snapshot worktree and nothing else, so a tool server that loads is not covered by any of them.

## Impact

- `internal/agent`: the arguments every invocation carries, the stream parser, what the agent records, and the deletion of the retry along with the flag that needed it.
- `internal/claudesim`: the simulator learns the allowlist argument, honours it, and reports tool servers and tools in its `init` event, because the tests assert on the invocation that was delivered. Its refusal behaviour goes, because nothing sends the flag that triggered it.
- `internal/fidelity`: the recorded shapes gain the `init` event's `mcp_servers` and `tools`, which is what makes a vendor change to either one visible weekly rather than on the next live run.
- No stored run changes shape, and nothing about what socrates writes in a reviewed repository changes.
- **BREAKING for nobody's interface**: the event a parked round writes is unchanged, and the `agent-failed` line this removes from a refused startup was never a state anything branched on.
