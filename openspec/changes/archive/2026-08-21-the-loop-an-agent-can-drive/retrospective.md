# The run

`the-loop-a-session-drives`, `opus` at `high` effort, round 1, against base `d4d527e`, recorded on 2026-08-21.
It escalated with 8 findings, 2 of them for the person, at `medium` risk, and it is recorded whatever it found.
The corpus now holds six answers and `internal/corpus/outcomes/distribution.txt` was rewritten to say so.

That is the whole of what this task may require of a live run.
What follows is what the run was for: the places the sequence needed a sentence the skill did not have.

## What the driving session was, and was not

The session that drove the loop was the one that wrote the change, and it had been told the sequence at length before the skill existed.
So this run is evidence that the sequence the skill names works end to end against the real CLI.
It is not evidence about the wording, because the reader already knew the answer.
A naive session driving it is still worth doing, and it is not what happened here.

The person authorised the run, and then authorised fixing what it found, which is why the escalated branch was not obeyed to the letter: the skill says hand the findings over and stop, and stopping is what happened until the person said otherwise.

## Where the sequence needed a sentence

**No polling interval and no bound.**
Step 3 says to poll `scrat status` and does not say how often or for how long.
The driving session invented five seconds and an unbounded loop.
A session that picks a fifth of a second spends the review hammering a command that reads `run.json` and a findings document each time, and a session that picks a bound too short reports a live review as stuck.
The skill still does not say, because what to say depends on `agent_deadline`, and that is a sentence to write once somebody has driven this against a slow review rather than a guess made here.

**`cancelled` had no branch, and the loop had no exit for it.**
The reviewer found this as f1 and it is the sharpest thing the run produced.
The skill said "one of the four states below", the four being `clean`, `failed`, `escalated` and `awaiting-response`, and a run also rests in `cancelled`.
The concrete path is one the skill itself opens: it tells the session to suggest that the person run `scrat watch` in another terminal, and README documents `watch` stopping when the person cancels from there.
So the person cancels, `scrat status` prints `cancelled` for as long as anybody lets it, `respond` and `continue` both refuse, and the session polls or improvises.

The gap was in the specification before it was in the text.
The requirement named four states, the test named the same four, and nothing in the suite could notice a fifth.
`cancelled` is now a requirement, two scenarios and a section, and `TestEveryStateARunRestsInIsNamed` compares what the skill names against `runstate.States()` rather than against a list somebody typed twice.

**The escalated branch gave one cause for escalation.**
It explained escalation as a finding that questions what the change is for, which is one of two ways a run escalates.
The other is the round cap: a second round that still holds something blocking escalates with ordinary patchable findings.
A session reading the old text would have handed those over with the wrong explanation, and `scrat findings` would not have corrected it, because the recorded reason prints only where nothing was reviewed.
The section now names both and says that `scrat status` prints which one it was.

**Nothing said the run identifier is optional.**
Every command takes the most recent run for the repository you are standing in when no identifier is passed.
The driving session passed it every time, which is right and is what the skill implies, and a session that omits it is equally right and the skill does not say so.
Left alone: it costs a session nothing to pass what `scrat review` printed, and one fewer sentence is one fewer thing to get wrong.

## Where the checks were weaker than they read

**A flag check that could pass while the skill named a missing flag.**
`surfaceOf` built each command's flag set with `flagsIn`, which collected every `--` token on every line of a `Flags:` block, including tokens inside another flag's usage text.
`--rejected`'s usage names `--reason`, so `--reason` was reported as offered by `respond` on the strength of a description.
Deleting the `--reason` registration from `cmd/scrat/respond.go` left `skill.Unavailable` silent while the skill went on telling sessions to pass it.
`flagsIn` now takes the first `--` token per line, and that mutation was run: it now fails, naming the flag.

`TestARemovedFlagFailsTheSuite` also removed only the first flag-bearing invocation, which was `respond --fixed`, and `--fixed` appears in no other flag's description, so the hole was exactly where the guard was not looking.
It now removes each named flag in turn, as a subtest per flag.

**A guard that could not fail.**
The check that `scrat watch` appears in no loop the session runs was written as a conjunction whose second half was already an error, so it could never report anything the line above it had not reported.
The suggestion also sat inside numbered step 3 of the loop, which the specification's own scenario says it must not.
The suggestion moved out of the loop entirely, into the paragraph that says what the skill is for, and the guard is now an index comparison against the loop's heading.

**A section helper that read past its section.**
`sectionAbout` bounded a section at the next `###` only, so the last one, `awaiting-response`, carried the whole prohibitions section with it.
Two tests were asserting over a region twice the size of the section they named, and moving `scrat continue` out of the parked-run instructions and into the prohibitions would have left `TestTheLoopCloses` green with the loop no longer closing.
It now stops at `##` as well.

**A check that asserted the wrong thing.**
The half of `TestNoChannelIsImpliedForTheAuthorsAccount` that scanned for filenames would have failed on `AGENTS.md`, on `~/.claude/skills/scrat/SKILL.md`, and on the words "and/or", with an error message describing none of those.
A filename is neither necessary nor sufficient for the thing the scenario is about.
It is replaced by an allowlist of the commands the skill may name, beside the one it already had for flags, each entry saying why that command is not a channel for the author's own account.

## Two notes, decided rather than fixed

**The reviewer can see the skill.**
`scrat skill` installs at `~/.claude/skills`, which is exactly the tree the reviewer loads, because socrates invokes the CLI with user-level settings only and that is the invariant keeping the reviewed project's instructions away from it.
So every review after an install spawns a reviewer whose skill list holds an entry telling it how to drive a review.
It is bounded: a reviewer's tools are `Read`, `Grep` and `Glob`, so it cannot run `scrat`, and `scrat review` refuses inside a review.
Recorded rather than fixed, because the alternative is installing where the vendor does not look, and README now says it.

**`.serena/` was in the snapshot.**
The reviewer reported a 133-line Serena config as something this change committed.
It committed nothing: the directory is untracked, and snapshots include untracked files on purpose, because a freshly implemented feature is mostly untracked files.
The reviewer was right that it had nothing to do with the change, and reading it was attention spent on nothing.
It is in `.gitignore` now, so it stays out of the snapshot rather than out of the commit.
