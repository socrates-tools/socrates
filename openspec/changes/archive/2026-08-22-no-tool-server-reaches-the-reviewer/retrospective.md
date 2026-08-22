# The run

`no-tool-server-reaches-the-reviewer`, `opus` at `high` effort, round 1, against base `7a3c3b0`, recorded on 2026-08-22.
It parked with 5 findings at `medium` risk, and it is recorded whatever it found.
The corpus now holds nine answers and `internal/corpus/outcomes/distribution.txt` was rewritten to say so.

That is the whole of what the verification task may require of a live run: the run happened, its answer is in the corpus, and socrates routed what came back.

What the task did require, and what the run reported, is that the invocation carried the empty allowlist and that the session loaded nothing.
The opening event read `mcp_servers: []` with four tools: `Glob`, `Grep`, `Read`, `StructuredOutput`.
Two live runs earlier the same session held 52 and then 74 tool-server tools on the same machine, including `mcp__playwright__browser_run_code_unsafe` and `mcp__chrome-devtools__evaluate_script`.

## How the change was found

Not by a review, and not by a test.
A person said the `--strict-mcp-config` refusal in every run's log was annoying, and asked for socrates to stop tripping on it rather than for the machine to be reconfigured.

Reading the refusal properly is what turned a log annoyance into the change.
The fallback that runs the review without the flag was specified, implemented and tested, and it was accepted on the belief that the seal answers a loaded tool server.
The seal bounds the snapshot worktree.
Nothing bounds a browser-control tool, and the requirement admitting the fallback said so in its own second sentence, which nobody had read against what the fallback actually left running.

The measurement is what settled it: a probe on this machine reported 77 tools, 74 of them from servers.
The invariant in `AGENTS.md` claimed the reviewer loads no external tool server, and on every machine with a managed configuration it had been false since the day it was written.

## What the first attempt got wrong

The first version of this change kept `--strict-mcp-config` and added the allowlist beside it, on the reasoning that two mechanisms beat one and that only the flag stops the configuration files being read at all.

That reasoning survived until somebody asked why socrates could not simply want what an ordinary session wants.
The flag is refused wherever a managed configuration exists, and that refusal was the sole origin of the second invocation, the retry, the guess at why the first attempt died, and a line in every run's log reading `the reviewer failed` when nothing had.
What it bought over the allowlist is that the files are not read, and with no server permitted, reading them permits nothing.

So the flag was pure cost.
Dropping it deleted the retry, the second invocation, the `agent-startup-rejected` event and its rendering, the simulator's refusal behaviour, and four tests whose premise no longer existed.
The change got smaller and the protection got stronger, and the four fronts the author had been defending were three fronts the flag had created.

## Where the checks were weaker than they read

**A one-directional drift check would not have noticed the field going away.**
`fidelity.Compare` reports what the simulator produces that the recordings do not, and says nothing about a recording field the vendor stopped sending.
socrates now depends on the opening event's `mcp_servers` and `tools`, so the vendor removing either would be a silent total failure of the reading path.
Fixed by a guard in the ordinary suite over the checked-in recordings, and the guard was mutated to prove it fails.

**Two of the author's own tests could not fail.**
`Contains(line, "2")` was satisfied by the `09:14:02` timestamp of the test's fixed clock, and `Contains(printed, "person")` was satisfied by the unrelated `awaiting responses` line.
Both were found by the live reviewer of the change before this one, and both were rewritten to assert the noun beside the number.

**A test whose name claimed what it did not exercise.**
`TestARewordedRefusalIsStillRetried` was vacuous, because the simulator always wrote the canned refusal.
It was fixed, and then deleted with the retry it was about.

## What the reviewer found, and what was done

Five findings: three asking a person, two informational.
None of them was about the mechanism, and four of them were defects.

**The fidelity guard was set to fail the first time the fixtures are re-recorded.**
It required `mcp_servers[].name` from every recording, and recordings are taken through socrates' own invocation, which now permits no server.
The next `make fixtures` records `mcp_servers` as an empty array, no `[].name` is derived, and the guard fails reading as though the vendor had removed a field socrates depends on.
Fixed: `mcp_servers` and `tools` are required from every recording, and `[].name` is required only from a recording that lists a server.
Both directions are held by mutation, and `docs/gotchas.md` records that doing this right takes the evidence for that field away, so the pin lasts only while the two pre-allowlist recordings stand.

**A subtest's name claimed a case the simulator cannot express.**
"A machine that names them in its own settings" was the previous subtest's scenario with one more server, and a user allowlist widening socrates' own is the case the design says socrates cannot prevent, so the assertion was the opposite of the documented behaviour.
Fixed by naming the two cases what they are, one server configured and several.
The case socrates cannot prevent stays covered by the test where the machine ignores the allowlist outright.

**An old event rendered as a sentence with its content missing.**
The loaded-servers line was rewritten to read `servers` and `tools`, and every such event written before this change carried `reason` instead.
A run still inside the 336h retention window rendered as "the machine's tool servers were loaded and socrates could not stop them: , adding 0 tools...", losing the only account it held.
The design's migration plan had considered `agent-failed` and not this.
Fixed: the renderer falls back to the recorded reason where no servers are named, and a test supplies the old shape on its own, because the test that covers every kind at once supplies both shapes together and could not see it.

**An exported field nothing consumed.**
`Response.Opened` was written on both paths and read nowhere.
It was the distinction between "the stream reported no servers" and "the stream never opened", which nothing asks: a stream that never opened belongs to a run that failed, and socrates records no servers because it saw none.
Dropped, along with the parser field behind it, which the retry had been the last reader of.

**A server the vendor lists but never started reads as loaded.**
The reviewer noted that the per-server `status` was parsed into `wireToolServer.Status` and never consulted, so a server named but not started would print "socrates could not stop them" beside a count of zero tools.
Asked, the person declined to arbitrate a tool-server detail and told the author to own it.
The call kept is that presence in `mcp_servers` is the signal, because a server the vendor lists is a server the allowlist did not stop, whatever its connection status, and whether its tools arrived is already answered by the session tool list recorded beside it.
The unread field is deleted rather than left parsed and ignored, and the reasoning is stated in the design.

## What is left, and why it is left

Round 2 was available and was not spent: the person chose to archive rather than have a fresh reviewer read the five fixes.
Four of the five are defects with one sensible resolution each and the fifth is a decision, so what a second round would have judged is mostly the evidence, which was mutated instead.

The mechanism rests on one vendor setting, which the reviewer named as the risk and which the design takes deliberately.
Nothing in the suite can verify that `allowedMcpServers` still means what the vendor documents; what covers it is that socrates reads the opening event rather than trusting the request, that the fields it reads are pinned by the recordings, and that the weekly fidelity workflow is where a vendor change surfaces.
