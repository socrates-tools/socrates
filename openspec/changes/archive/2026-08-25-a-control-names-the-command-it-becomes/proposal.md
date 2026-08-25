## Why

The page a person reads is rendered through `artifact-design`, which puts a document quoting the code it was written about on a host outside the machine.
`internal/corpus` refuses to record such a document for that exact reason, and the skill tells a session to publish one.
That asymmetry is recorded as issue 1 and has been open since the page existed.

The page also collects nothing, and the reason it collects nothing is stated in the spec: there is nothing behind it, so a control would take a decision that reaches no run.
That reasoning was correct and is now false.
`a-page-that-answers-back` builds the thing behind it, and a control on the page can name the command it becomes and hand that command back to the session.

So the rule inverts, and it inverts for the reason it was written rather than in spite of it.
What replaces "the page carries no control" is not "the page may do anything".
It is that every control names the command it becomes, and nothing is recorded until the session runs it.

## What Changes

- **The skill drives `agora` instead of `artifact-design`.**
  The session writes a document, serves it, blocks on `agora poll`, and records what came back.
  Nothing about the page leaves the machine, and the asymmetry issue 1 records is closed rather than restated.
- **The page stops being read only, and the reason it was read only is kept.**
  A control is permitted exactly where it names the command it becomes.
  Nothing on the page is recorded against a run, and the session running `socrates respond` is still what records it.
  A page that stored an answer, or that carried a control whose command it could not name, is refused by the same reasoning that refused all controls before.
- **The person's own words reach the record.**
  Today a decision is recorded in the session's paraphrase of what the person said in the terminal. It becomes the reason string the person's chosen option carries, or the words they typed.
  `socrates respond --decided` is unchanged, and it is still the session that runs it.
- **The identifiers stop being a matter of the session's care.**
  agora copies each finding's identifier, severity and audience out of `socrates findings --json`, so the requirement that they match is met mechanically rather than by instruction.
- **The surface check learns there are two binaries.**
  Every command and flag the skill names is checked against the binary that offers it. Today the extraction only sees words following `socrates`, so a skill naming an `agora` flag that does not exist would pass unnoticed, which is the failure this check exists to prevent.
- **Where the page is drafted is revisited.**
  The inline-against-delegated split was decided when drafting meant writing a whole HTML page. Writing a document is a fraction of that, so the reason for delegating at a parked round is re-derived rather than carried forward unexamined.
- **socrates still renders nothing, publishes nothing, stores nothing and verifies nothing.**
  No command and no flag is added to `socrates`, no run record carries anything about a page, and socrates gains no network reach. The surface is the session's reach, exactly as it was when the surface was somebody else's host.

## Capabilities

### Modified Capabilities
- `findings-page`: the handover names the local surface rather than `artifact-design`; the read-only requirement inverts to a control that names the command it becomes, with the record still made in the terminal; the identifier match becomes mechanical; and where the page is drafted is re-derived against what a document now costs.
- `agent-skill`: the branches that hand a round over drive `agora`, and every command and flag the skill names is checked against the binary that offers it rather than against one binary.

## Impact

- `internal/skill`: the page section, `PageBar`, the `escalated` and `awaiting-response` branches, the surface extraction, and the regenerated `.claude/skills/socrates/SKILL.md`.
- `internal/e2e`: the journeys that drive the loop drive the handover too, against the real `agora` binary, so the sequence the skill describes is executed rather than read.
- **No change to** `internal/dialogue`, `internal/findings`, `internal/reviewer`, `internal/runstore`, `internal/report`, `internal/agent` or `internal/prompt`. Which rounds a person reads, what a round leads to, and what `socrates findings` reports are all already right.
- No new command, no new flag, no new dependency for `socrates`, and nothing written in the repository under review.

## Dependencies

This change depends on `a-page-that-answers-back`, which builds the binary the skill here names. The skill cannot name a surface that does not exist, and the check that every named command is real would fail first.

Both changes further depend on `a-round-a-person-reads-becomes-a-page` being archived, because `findings-page` lives in that change rather than in `openspec/specs` and the deltas here have nothing to compare against until it is promoted.
That archive is not this change's to perform.
