## Why

A page is served only where a round cannot move until a person answers, so the two most common outcomes of a review, a clean run and a round the author patches alone, reach the person as terminal output or as nothing at all.
The page exists because a person reads a review result better on a page than in a scrollback, and that is as true of a round nobody is being asked to answer as of one they are.

The condition is doing two jobs at once and only one of them is load bearing.
It decides whether a page is served, and it decides whether the session stops and waits on it.
The second job is real: a session that waited on a clean run would hang forever.
The first is what makes a person's only view of a clean review a sentence in a terminal.

## What Changes

- Every round that completed a review is handed over as a page: `clean`, `awaiting-response` whether or not it holds a finding a person answers, and `escalated`.
  A `failed` or `cancelled` run gets no page: a failed run completed no review, and a cancelled one was stopped on purpose, so what its first round got as far as is not a result anybody is still waiting for.
- Whether the session waits on the page becomes a separate fact from whether it serves one.
  It waits where a decision on the page is the person's, and does not wait otherwise: it hands over the address, says briefly what the page holds, and gets on with patching or reporting.
- **BREAKING** The one rule stays one rule and is renamed for what it now decides.
  `dialogue.APersonReads` becomes `dialogue.AwaitsAPerson`, the line `socrates findings` prints changes from `read by a person` to one naming the wait, and the `--json` field `read_by_a_person` becomes `awaits_a_person`.
  Nothing consumes that field but the skill, which this change rewrites.
- **BREAKING** A document may hold no decision.
  `agora serve` refuses one today, on the reasoning that there would be nothing for a person to answer, which is exactly the page a clean round needs.
  A document must still say something: a title, and at least one of a decision, a source entry or a preamble.
- The source a document cites is rendered.
  Today it is copied onto the decisions that cite it and is otherwise invisible, so the short list of the rest of the round is prose the session retypes.
  The page shows every entry of the source as the round at a glance, taken from the source rather than retyped, and shows in full the ones a decision covers.
- A page holding no decision offers no way to answer and no way to send.
  There is nothing to carry, and a control that delivers an empty answer is a control a person can wonder about.
- The skill teaches one shape rather than a branch: draft a document, serve it, hand over the address, and wait only where socrates says the round waits on a person.
  Where to draft it follows what a token costs at that branch, which is the rule already recorded: inline where nothing follows, and delegated where another round does.

## Capabilities

### New Capabilities

None.
Every part of this is a change to what an existing capability requires.

### Modified Capabilities

- `findings-page`: which rounds are handed over as a page, and the separation of serving one from waiting on it.
- `review-dialogue`: the rule renamed for what it decides, and stated as the wait rather than the reading.
- `cli-foundation`: what `socrates findings` reports, as a line and as a field.
- `agent-skill`: the four branches of the loop, each serving a page, and which of them waits.
- `surface-document`: a document that holds no decision, and what a document must carry.
- `surface-page`: the source rendered as the round at a glance, and a page with nothing to answer offering neither control.

## Impact

- `internal/dialogue`: the predicate and its name, and the exhaustive test over the whole input space.
- `internal/report`: the reported field, its line, and the message that states it.
- `internal/skill`: the loop's four branches, the page section, and the tests that read the rendered skill.
- `internal/document`: validation of a document holding no decision.
- `internal/render`, `web/`: the source list on the page, and the send control's absence where nothing is answered.
- `internal/e2e`: the handover tests, the loop tests, and the test that asserts which rounds socrates says a person reads.
- `AGENTS.md`: the invariants naming `dialogue.APersonReads` and what socrates states about a page.
- No new command, flag or run-record field on `socrates`, and no network reach.
