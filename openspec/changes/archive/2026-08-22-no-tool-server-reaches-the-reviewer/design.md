## Context

The invariant is that the reviewer loads no external tool server.
It was held by one flag, and that flag has a documented case where the vendor refuses it: an enterprise MCP configuration is a managed set, and `--strict-mcp-config` asks to replace it, so the CLI exits at startup rather than let a flag override policy.

socrates handled the refusal correctly and then trusted the wrong thing.
The fallback was accepted because the seal holds the reviewer read-only, and the seal holds the *snapshot worktree* read-only.
A tool server offering script execution, file upload and browser control is not bounded by a worktree, and the requirement that admits the fallback says so in its own second sentence.

This change does not relax the fallback.
It replaces the flag with a mechanism no machine refuses, deletes the fallback that only the flag needed, and stops socrates guessing about what the CLI did.

## Goals / Non-Goals

**Goals:**

- No tool server loads into a reviewer, on a machine with a managed MCP configuration or without one.
- What socrates records about tool servers is what the CLI reported, not what socrates inferred from an error message.
- A person reading a run's log is not told the reviewer failed when nothing failed.
- One invocation per round, because the second one only ever existed to work around a flag.

**Non-Goals:**

- Changing the machine's own configuration. The enterprise file is how the person's ordinary sessions get their servers, and socrates has no business editing it.
- Refusing to review where servers load anyway. Considered and rejected below.
- Classifying which loaded tools can write. Also below.
- Restricting tool servers for anything other than the reviewer. socrates starts one kind of process.

## Decisions

### An empty allowlist, supplied per invocation, and nothing else

`--settings '{"allowedMcpServers":[]}'` goes on every invocation, and `--mcp-config` and `--strict-mcp-config` are removed from it.

The allowlist is documented to sit in any settings file, to merge from every source, and to be evaluated before a server loads, including one from a managed configuration.
An empty array is documented as distinct from unset: unset allows everything, empty allows nothing.
Measured on a machine with four managed servers, the reviewer's session goes from 74 tool-server tools to none, and the CLI does not object to the setting the way it objects to the flag.

A denylist was the other candidate and it cannot express this.
Denylist entries match a URL, an exact command, or a name; `serverUrl` takes wildcards but `serverCommand` matches every argument exactly and `serverName` does not expand, so "deny everything" is not writable for stdio servers without enumerating the commands on the machine.

The flag was kept at first, on the reasoning that two mechanisms beat one and that it alone stops the configuration files being read.
That reasoning does not survive contact with what the flag costs.
It is refused wherever a managed configuration exists, and that refusal is the sole origin of the second invocation, the retry, the guess about why the first attempt died, and a line in every run's log saying the reviewer failed when nothing had.
What it adds over the allowlist is that the files are not read, and with no server permitted, reading them permits nothing.
So the flag was pure cost and it is gone. One mechanism, never refused, one invocation.

### What loaded is read from the `init` event

The stream's first `system` event names `mcp_servers` and the session's `tools`.
socrates reads both, records them, and stops deriving anything from the wording of a vendor error.

Before this change the retry was decided by `strings.Contains(err, "cannot use --strict-mcp-config")`.
That was a latent total failure: the vendor rewords the message, socrates stops recognising the refusal, and every review on such a machine fails outright rather than falling back.
Dropping the flag removes the decision entirely rather than improving it, which is the better answer: there is no longer anything to fall back to, so nothing has to guess why an attempt died.

What the `init` event is for is the other half. It reports what the session actually loaded, which is the only thing that can notice an allowlist a machine's policy overrode, or a server that loaded for a reason nobody predicted.

Of each server socrates reads the name and nothing else.
The event also carries a connection status, and a server that failed to connect is still a server the allowlist did not stop, which is the thing being recorded.
Whether its tools arrived is answered by the session's tool list rather than by its status, so reading the status would add a third gradation that neither the requirement nor the line a person reads has any use for.

### Where servers load anyway, the review still runs

Unchanged, and deliberately.

Refusing was considered, by analogy with the seal: socrates refuses to review where the seal cannot hold, because a verdict it cannot stand behind is worse than no verdict.
The analogy is weaker than it looks. A filesystem that ignores permission bits defeats the mechanism socrates depends on for every run and there is nothing left. An allowlist an administrator has locked down is a policy decision by somebody entitled to make it, and refusing would mean socrates reviews nothing on that machine until an admin acts.
The run records what loaded, the seal and the comparison still hold the worktree, and a person reading the record can decide what it is worth.

Refusing only where the loaded tools can write was considered and rejected as the worst of the three.
It requires socrates to classify a vendor's tool names by whether they write, which is a list that goes stale silently and is wrong in the direction that matters: a tool socrates does not recognise reads as harmless.

### The simulator learns the allowlist and loses the refusal

The tests assert on the invocation that was delivered, so the simulator has to accept the allowlist argument, honour it, and report `mcp_servers` and `tools` in its `init` event for the recording path to be testable at all.
Its refusal behaviour goes, because nothing sends the flag that triggered it, and a scenario nothing can reach is a scenario that rots.
The fidelity recordings gain the same two fields, so a vendor that renames either is a weekly failure rather than a surprise on somebody's live run.

## Risks / Trade-offs

- **An administrator can lock the allowlist.** `allowManagedMcpServersOnly` makes managed allowlists the only ones that count, and then the empty allowlist socrates supplies is ignored. → Nothing silently degrades: socrates reads the `init` event, sees the servers, and records them, which is the same fallback as today with the truth attached rather than inferred.
- **The setting is now the only mechanism, so a vendor change to it is a change to the whole protection.** → It is documented, the fields socrates reads are pinned by the fidelity recordings and by a guard in the ordinary suite, and the `init` event says what actually loaded, so the failure is visible rather than silent. The flag it replaced was refused on the machines that matter, so keeping it as a second mechanism bought nothing there.
- **`--settings` merges rather than replaces.** socrates passes `--setting-sources user`, so what it supplies is added to the user's own settings, and the vendor documents that a person's own `allowedMcpServers` can widen an allowlist. → Not preventable from inside the invocation, and not silent: the `init` event names what loaded and the run records it.
- **One mechanism means one thing to get wrong.** → That is the trade taken deliberately. The `init` event is the arbiter and reports what actually happened whether the mechanism applied or not.

## Migration Plan

No stored run changes shape and no run in flight is read differently.
A run recorded before this change holds an `agent-failed` event for a refused flag, and the renderer still renders that kind, so old logs read as they always did.

Rollback is the arguments: putting `--mcp-config` and `--strict-mcp-config` back restores the old behaviour, and with them the refusal and the retry that made this change necessary.

## Open Questions

- Should `scrat status` say that a run's reviewer had tool servers loaded, rather than only the run log? A person reading a verdict may want to know the reviewer was not the one the design describes. Left out here because the counts a parked run prints are already the surface a change argued over, and this is a different question about verdicts.
- Should the reviewer's whole tool list be recorded, or only the tool servers? The list is 77 entries on this machine and most of the value is the count and the server names. Recording the names of the servers and the number of tools is what this change does, and a person who needs the list has the stream on disk.
