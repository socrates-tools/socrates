# 38. A design session has no round cap

- Status: accepted
- Date: 2026-09-17
- Part of: socrates-design

## Decision

A design runs as many rounds as the frontier has questions.

Review caps rounds at two, nothing configures it, and a configuration file naming one is refused.
That cap exists because a round spawns a model process against a snapshot, and an uncapped review is an uncapped spend that a person is not watching.

A design round spawns nothing.
It costs the person's attention, and the person is the one spending it, in front of the page, deciding each round whether there is another.

So the frontier decides when to stop, and the person can end it at any round.
Where they end it early, the spec is written from what is settled and the unsettled questions go to Further Notes, because four rounds of work are worth more than nothing.

## Rejected

- Capping design rounds for symmetry with review.
  The cap's reason does not transfer, and a cap that stops an interview mid-tree produces a spec that states neither what was agreed nor what was still open.
- Letting a configuration file set a cap.
  Review refuses one for the round cap it has, and a setting nobody needs is a setting to keep working.
