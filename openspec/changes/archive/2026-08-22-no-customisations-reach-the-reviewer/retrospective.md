# What this is, and what it is not

This was written after the change was archived, from the corpus, the commits and the change's own documents.
Every claim in it is one of those three, and where the record is silent about why something was done it says so rather than filling it in.
That is a weaker retrospective than one written by the session that made the decisions, and it is here because the alternative is the only archived change with no account of itself.

# The runs

Two runs, four rounds, all `opus` at `high` effort against base `26ada92`, recorded on 2026-08-22.

| run | round | findings | risk |
| --- | --- | --- | --- |
| `20260822T035656Z` | 1 | 4 | `medium` |
| `20260822T035656Z` | 2 | 6 | `medium` |
| `20260822T051041Z` | 1 | 5 | `medium` |
| `20260822T051041Z` | 2 | 6 | `medium` |

Both runs ended `escalated` on the round cap, and both second rounds are in the corpus, which now holds thirteen answers.
Two decisions were recorded against each first round, which is what the second round's own subject line says it was reading.

That is the whole of what the verification task may require of a live run: the runs happened, their answers are in the corpus, and socrates routed what came back.
What task 8.5 did require, and what all four corpus entries record in their own subject lines, is that each round was answered by a session holding the tools socrates asked for and no others.
The tool list that session held is in the design's table rather than here, because that was a probe and these were reviews.

# What the change did

It replaced a vendor feature named through five packages with one flag, and gave the reviewer the web.

The property was never about one source of tools: it is that the reviewer holds the tools socrates asked for and nothing else.
`--safe-mode` disables the customisations the machine brings on its own account, it is not refused where a managed configuration is deployed, and the change measured it rather than trusting the documentation: four external sources named without it, none with it, and a plugin's own tool gone from the session's list at the same time.
With it, the settings allowlist, the refused flag, the retry, an event kind, a rendered sentence and four tests whose premise no longer existed all left the program.

`WebSearch` and `WebFetch` joined the three reading tools, because a reviewer that cannot check what a library does writes a finding about its own recollection.

# What the reviewer found, and what was done

Twenty-one findings across four rounds, and none of them was about the mechanism.
Every one was about what the documents claimed the mechanism did.

**The claim was broader than the flag.**
The rewritten requirement said the one mechanism disables every customisation the machine brings, and the vendor's own description keeps the managed settings policy in force.
Round 1 asked a person; the answer was to state the residual rather than to widen the flag.
Round 2 then found the carve-out written in the other direction, as though every policy-configured customisation survived, which would have made the deleted allowlist a regression on exactly the machines it was measured on.
It now says what the vendor says: a policy-configured hook, status line and file-suggestion commands survive, and a policy-configured tool source does not.

**The requirement contradicted itself twice, fifteen lines apart each time.**
It justified the tool-list check with "a hook contributes tools" and justified the residual with "a hook contributes no tool".
It also still required the reviewer to load user-level memory, two paragraphs above the sentence requiring the mechanism that removes memory.
Both were the stale half of a rewrite, and both were corrected to the accurate half.

**The README claimed less protection than the code gives.**
It said the policy hook's command runs outside the worktree the seal bounds.
The hook's command is a child of the reviewer's session, and the session's working directory is the sealed snapshot worktree, so two of the three mechanisms do bound a write it makes there.
What is genuinely uncovered is where that command reaches outside the worktree, which is what the requirement and the README now say.

**"Measured" was doing work the measurement did not do.**
The record holds one measurement: four sources went to none, and one plugin-contributed tool left the tool list.
Nothing recorded distinguishes a managed plugin, a managed skill or managed memory from a person's own, because the session's tool list does not say what contributed a tool.
The prose was split into what socrates measured and what the vendor documents, and the remaining half of that split was carried as informational into the second run.

**The enumeration of what the flag disables lost items twice.**
Once it dropped external tool sources, which is the entire reason the flag replaced the allowlist, leaving the file arguing with itself.
Once it dropped the status line, the file-suggestion commands, LSP servers and auto memory, which made the residual stated below it read as broader than it is.
Both were fixed in the file the reviewer named.

**A machine-specific observation was written as a general one.**
`docs/gotchas.md` said that with the flag no hook runs and the opening event is first.
On a machine whose policy configures a `SessionStart` hook it would still fire, and the entry now says so.

**The change's own documents kept the claim the specifications had dropped.**
The proposal and the design still opened with the unqualified version, which is where a later reader goes for the reasoning.
They were qualified, one of them only after a second round found the sentence the first round's fix had missed.

# Where the checks were weaker than they read

Two findings were recorded four times each and answered by neither run.
They are the reason `absence-of-evidence-is-not-a-pass` exists, and both are the same defect: a check whose empty answer is printed in the shape of its clean answer.

**A session that never said what it held read as a clean session.**
socrates subtracts the tools it asked for from the session's tool list and records the remainder.
An opening event carrying no tool list leaves that remainder empty, so the run recorded nothing, and the suite already knew the distinction mattered: two tests fail with `the session reported no tools at all, so this test would pass on an empty reading` before they assert anything.

**The drop gate was disarmed by the very change that taught it renames.**
The promoted specifications were written at the same time as the delta, so the gate's before-state and after-state were byte-identical, the five declared scenario renames matched no promoted name, and `make scenario-drop` was run, ticked and reported nothing about this change.
The reviewer named it in round 1, again in round 2, and again in the first round of the second run.
This change was archived in that state, which is what the gate that now fails on it was written from.

# What has since been reversed

One decision in this change's design is no longer socrates' position.

It rejected reading the sources the opening event names, on the reasoning that doing so "would swap one incomplete signal for another, and would keep the vocabulary this change removes".
The first half is true of a swap, and reading both is not one: the tool list misses a source whose tools have not arrived, the source list misses a built-in the vendor added, and this change's own design recorded the measurement that proves the first gap is real, an opening event naming four sources whose tools were not in the list yet.
The second half was a real constraint and it is met: the vendor's field name appears once, in a JSON tag beside every other vendor field name, and no requirement, event kind, recorded field or rendered line says anything but source.

What made the difference decidable is a measurement this change took and did not use for it: under the flag the sources named went from four to none while the `plugins` inventory beside them did not move.
That is what separates a field which answers for the session from one which inventories the machine, and it is why one is read and the others are not.

# What is left, and why it is left

The residual is stated rather than covered: a hook an administrator's policy configures runs, it contributes no tool for the tool list to see, and where its command reaches outside the sealed worktree nothing socrates holds bounds it.
That is accepted for the reason the review runs at all on such a machine, which is that a policy set by somebody entitled to set it should still get reviews.

The reviewer loses the person's own memory and skills, which is deliberate and points the same way as refusing the reviewed repository's instructions.
The benefit is the one the README now claims: the `scrat` skill a person installs cannot seat the reviewer that reads their code.

Both runs ended on the round cap rather than on a person deciding they were finished, so no fresh reviewer ever read the fixes the second rounds prescribed.
The second run is what read the first run's fixes, and two of the five findings it opened with were the two the first run had already recorded twice.
Nothing has read the second run's, and the two findings it repeated were answered by the change after this one rather than by a third round, which the cap makes a fresh review somebody asks for.
