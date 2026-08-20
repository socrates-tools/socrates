## Context

The core establishes the asymmetry that makes socrates worth having: the reviewer has independence and never authored anything, the author has context and does the patching.

What the core does not have is the return path.
Findings go one way, and whether the author agreed with them is not recorded anywhere, not checked by anyone, and not visible afterwards.

This layer adds the return path and nothing else.
It changes no verdict, no exit code and no part of the agent seam.

## Goals / Non-Goals

**Goals:**

- An author who disagrees has somewhere to say so, and saying so is judged rather than accepted.
- A second reviewer that is genuinely fresh and genuinely informed: it inherits no session, and it inherits the claims explicitly.
- A bound on the loop that is chosen for the right reason and stated as the reason, so nobody raises it to three for cost.
- Everything that challenges intent leaves the machine and reaches the person, at the first round it appears.

**Non-Goals:**

- socrates patching anything. The author patches. socrates never writes in the repository, and that is unchanged by this layer.
- A negotiation protocol, a severity the author can lower, or a way to mark a finding permanently ignored. A rejection is a claim to be judged, not a suppression.
- Rounds beyond two.
- Any inference about whether a patch actually addressed a finding. That is the next reviewer's job, from the diff, and socrates asserts nothing about it.

## Decisions

### The routing rule, stated once

After a completed review of round `n`, with findings classified by the core's single helper:

- Any finding whose effective action is `ask-user`: `escalated`, whatever else is in the set. It challenges the author's deliberate intent, and a machine arguing with intent is the thing this design refuses to build.
- Nothing blocking: `clean`. Patchability does not enter it, because a finding that holds nothing is not something the author is being asked to answer.
- Nothing blocking that a response could address: `escalated`. A finding that holds the change while its own action asks for nothing is a contradiction, and a person resolves it rather than being asked to answer it.
- `n` is the cap: `escalated`, naming the cap.
- Otherwise: `awaiting-response`.

The rule is a pure function of the classified findings and the round number, so it is tested against its whole input space rather than sampled.

The `ask-user` case dominating everything else is deliberate and is the conservative reading.
A set holding one intent question and five mechanical defects goes to the person with all six, rather than patching five and holding one back, because the person's answer to the sixth may change what the other five should have been.

### A rejection is a written claim, not a veto

`scrat respond` records, per finding identifier, one of:

- `fixed`, with an optional note about what was done.
- `rejected`, with a reason that is required and is not accepted empty.

Both go to the next reviewer.
A `fixed` response is a claim too: the reviewer verifies it against the current code rather than believing it, which is the provenance clause's whole job.

`scrat continue` refuses while any blocking finding has no response.
Silence is not a rejection and it is not a fix, and a round that advanced on silence would present an unpatched, unargued finding to the next reviewer as though it had been dealt with.

### What socrates tells the author, and why socrates says it at all

socrates does not run the author's session, so it cannot constrain it.
What it can do is put the rules in front of it, at the moment it has findings in hand, in the output of the command it reads.

Four sentences, and the first one is the inversion:

> You did not review this code, an independent reviewer did.
> For each finding, either apply a fix or record a rejection with a concrete reason.
> A rejection is a claim you are making to the next reviewer, not a decision.
> Do not resolve a finding by weakening or deleting behaviour you added on purpose; if the finding is right about a symptom and wrong about the cause, say so and fix the cause.

The last sentence is the anti-lobotomy rule, which is what stops a finding being resolved by quietly deleting the work.

### The next reviewer receives the claims, labelled

This is the single explicit exception to "the reviewer never receives the author's rationale", and stating it as an exception rather than softening the rule is the point.

The prior findings and the responses arrive as a block labelled as claims to be judged, carrying:

- what was found, by identifier
- what the author said they did about it, or why they say it is wrong
- that both are claims and neither is evidence
- that a test written in the same round as the code it exercises is part of the claim, judged on whether its asserted outcome is the right outcome and whether it could still pass with the code wrong

It does **not** carry any instruction not to re-report a finding.
`no-mistakes` carries one because a finding a person chose to ignore should stay ignored there.
Carrying it here would let a rejection suppress the finding it disputes, which inverts the entire mechanism: the rejection exists to be judged, and a reviewer told not to re-report it cannot judge it.

### A round is snapshotted afresh, against the original base

Round two must see the patches as commits it can diff, not as a description of what changed, so the top of every round takes a new snapshot.

The base is the run's original base commit, not the previous round's snapshot.
Reviewing round two against round one's snapshot shows only the patches, which reads them out of the context they were applied to, and a patch that is wrong because of what it sits next to is exactly the kind this layer exists to catch.

### Neutralising is per field, on both channels

Findings were model output re-embedded in a prompt in the core.
Responses are model output re-embedded in a prompt too, and they are written by a session that is arguing.

For every field before it is re-embedded: neutralise conflict markers, normalise line endings, collapse whitespace runs, trim.
Per field rather than across the assembled block, so that a field's boundaries survive and a long description cannot swallow the label of the field after it.

## Risks / Trade-offs

- **A determined author can defeat this.** The rebuttal is free text and a session that wants to talk its way past a reviewer has a channel to try. What the design buys is that the attempt is written down, is judged by something that did not write the code, and ends at the person rather than at a pass. That is the honest claim; anything stronger would be a lie about what a text channel can enforce.
- **The cap of two escalates changes that a third round would have settled.** Accepted, and the reason is stated on the requirement rather than left to be re-derived, because the pressure to raise it will come from somebody looking at a run that nearly finished.
- **`ask-user` dominating the set makes escalation common.** A reviewer that defaults to `ask-user` when in doubt, which it is instructed to do, will escalate often at first. The alternative is a machine patching things that challenge intent, and the prompt's action definitions are where this is tuned, not the routing rule.
- **`fixed` responses can be lies.** By construction. Every one of them is verified against the code by the next reviewer, which is why a `fixed` response is carried as a claim and not as a fact, and why the round it appears in is not the round that certifies it.
