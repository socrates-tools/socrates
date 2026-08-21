---
name: scrat
description: Have an independent reviewer read the change you just made, and drive that review to a verdict. Use after implementing anything nobody else has read.
---

# Have somebody who did not write it review it

You have just implemented something and nobody has read it.
scrat review freezes your working tree into a commit nothing can move and hands it to a reviewer that wrote none of it.
You drive that review to a verdict, and every patch along the way is yours to make.
What the reviewer says while it works is not for you: if the person wants to see it happen, suggest that they run scrat watch in another terminal.

## How to answer a finding

- You did not review this code. An independent reviewer did, and it did not write any of it.
- For each finding, either apply a fix and record it with scrat respond --fixed, or record a rejection with scrat respond --rejected and a concrete reason.
- A rejection is a claim you are making to the next reviewer, which judges it. It is not a decision, and it does not make the finding go away.
- Do not resolve a finding by weakening or deleting behaviour you added on purpose. Where the finding is right about a symptom and wrong about the cause, say so and fix the cause.

## The loop

1. Implement what the person asked for.
2. Start the review, passing what the person asked for in their own words:

       scrat review --intent "count the characters of a document"

   It returns as soon as the run exists and prints the run identifier, which every later command takes as its argument.
3. Wait for the review to end: poll scrat status every five seconds or so until it reports one of the states below.
   Poll rather than follow what the reviewer is saying, because everything you read lands in your own context and you are about to need that context for the patching.
   scrat status prints the same handful of lines whether the reviewer was silent or said a great deal.
   Do not set a timeout of your own. socrates gives the reviewer a deadline, a reviewer that answers nothing before it passes ends the run as failed, and so every run reaches one of the states below without you bounding the wait.
4. Act on the state the run ended in. Each of them is below, and a run rests in no other.

### clean

An independent reviewer read the change and found nothing blocking.
Say so, say what was reviewed, and stop.
It means nobody found a problem, and not that the tests pass: socrates runs none.

### failed

Say that no review completed, and why. scrat status prints the reason.
Nothing has been checked, so this is not a clean result and not a pass.
Report the change as unreviewed, and offer to start another review.

### escalated

Read scrat findings, hand every finding to the person, and stop.
These are theirs to decide.
A finding that questions what the change is for is not one to patch your way out of, which is why the reviewer escalated it instead of asking you to answer it.
A run also escalates where a second round still found something blocking, and then the findings are ordinary ones that ran out of rounds rather than questions of intent.
scrat status prints which of the two it was.

### awaiting-response

Read scrat findings.
Every finding it reports as one that can be patched needs an answer from you, one finding per command:

       scrat respond --finding f1 --fixed --note "an empty document now answers minus one"
       scrat respond --finding f2 --rejected --reason "there is no retry loop in this function at all"

Then run scrat continue, and go back to waiting on scrat status.
The next round is a fresh reviewer that reads your answers and judges them, and there are two rounds at most.

### cancelled

Somebody stopped the review, so nothing was decided about the change.
Say that, and do not report the change as reviewed.
Whoever cancelled it did so on purpose, so ask before starting another one.

## What socrates will not do for you

socrates patches nothing.
Every fix is yours, made in your own working tree, and recorded afterwards.

There is no way to tell the reviewer about the change.
It gets the diff, the repository as it was frozen, and the one line of intent you passed to scrat review.
That is the design and not a limitation to work around: an account written by the session that wrote the code is the thing an independent reviewer exists to do without.

A rejection is a claim, and it goes to the next reviewer, which judges it.
It is not a way to make a finding go away, and a rejection without a reason is refused.
