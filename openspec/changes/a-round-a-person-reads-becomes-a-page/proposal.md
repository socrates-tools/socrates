## Why

Some rounds end with something only a person can act on, and socrates hands those over as a wall of findings in a terminal: an identifier, a severity, a path, a line number and a sentence written by a reviewer for another machine to parse.
A person deciding whether an endpoint stays has to reconstruct what the finding means for the thing they are shipping, from prose that was never written to tell them.

The two places it happens are the two places the person is the one acting: a run that escalated, where socrates is finished and the findings are theirs, and a round parked on a question about intent, where the session cannot move until they answer.
Everywhere else the session is the reader, the terminal is the right surface, and nothing about this change touches it.

## What Changes

- Which rounds a person reads becomes a rule socrates owns, beside the rule that decides what a round leads to: `escalated` always, and `awaiting-response` holding at least one finding that needs a person.
  A `clean` run gets no page, because nothing is being asked of anybody.
- `scrat findings` states it, as one line and as one field in its JSON, so that nothing downstream re-derives it from the state and the counts.
  `scrat status` is untouched: it is polled, every line it prints lands in the session's context, and the fact is only needed by the one command that is read once.
- The skill says to render such a round as an HTML page through the `artifact-design` skill, hand the person the link, and summarise rather than paste.
- The page is drafted inline where the run escalated and by a delegated subagent where the round parked, and the skill says which and why.
  At `escalated` nothing follows, so the drafting is paid once. At a parked round every token added is paid again on each of the turns that patch, respond, continue and judge, which is the context the round cap exists to protect.
- The quality bar the page is held to is stated once in the skill: plain English about impact rather than restated technical prose, an example where a finding is purely technical, a diagram only where a workflow needs one, and the decisions foregrounded with the rest of the round listed briefly.
- The page is read only, and it is matchable with the terminal in both directions.
  There is nothing behind it, so it collects nothing and declares no capability: a control on it would take an answer that reaches no run.
  Every finding keeps the identifier, severity and audience `scrat findings` reported, and each decision says what answering it looks like in the terminal, so the person reads, decides, and says which finding they decided in one exchange.
- **socrates itself reaches nothing new.** It states a fact about a round. The session publishes the page, so the network reach and the code that a findings document quotes leaving the machine are the session's, in the harness the person is already driving, and no build tag or new dependency appears here.

## Capabilities

### New Capabilities
- `findings-page`: when a round becomes a page rather than terminal output, what that page is held to, that it is read only and matchable with the terminal, who drafts it and why that differs by branch, and what socrates does and does not do about it.

### Modified Capabilities
- `review-dialogue`: which rounds a person reads becomes one rule beside the routing rule, over the same input, in the same package that reaches nothing outward.
- `cli-foundation`: `scrat findings` reports whether the round is one a person reads, in its lines and in its JSON, and `scrat status` still does not.
- `agent-skill`: the `escalated` and `awaiting-response` branches say to render the page, the quality bar is stated once, and the branch that delegates says so.

## Impact

- `internal/dialogue`: one predicate over the classification and the round, tested against its whole input space like the rule it sits beside.
- `internal/report`: one line and one field on `Findings`, and nothing on `Status`.
- `internal/skill`: the quality bar, the two branches, and the regenerated `.claude/skills/scrat/SKILL.md`.
- `internal/e2e`: the journey that drives the loop gains the assertion that the signal appears on exactly the branches the rule names and on no other. Whether a page was published is outside what a test can drive, so what is driven is the fact socrates states.
- No new command, no new flag, no new dependency, and nothing written in the repository under review.
