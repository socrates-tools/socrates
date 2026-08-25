package skill

const (
	frontMatter = `---
name: ` + Name + `
description: Have an independent reviewer read the change you just made, and drive that review to a verdict. Use after implementing anything nobody else has read.
---`

	purpose = `# Have somebody who did not write it review it

You have just implemented something and nobody has read it.
|socrates review| freezes your working tree into a commit nothing can move and hands it to a reviewer that wrote none of it.
You drive that review to a verdict, and every patch along the way is yours to make.
What the reviewer says while it works is not for you: if the person wants to see it happen, suggest that they run |socrates watch| in another terminal.`

	rulesHeading = "## How to answer a finding"
	ruleLine     = "- %s"

	pageHeading = "## The page a person reads"

	PageBar = `- Get to the point.
  Length is not thoroughness, and a person who has to search the page for their own decision has been handed the terminal output again.
- Say in plain English what the finding means for the change.
  Do not restate the reviewer's technical prose in a nicer typeface.
- Where a finding is purely technical, show a small example of it going wrong.
  A reader should not have to work out for themselves what the term costs them.
- Draw a diagram only where a workflow is what makes a finding hard to follow, and not otherwise.
- Spend tokens as though they were the person's, because they are.`

	thePage = pageHeading + `

Some rounds hold something only a person can act on.
|socrates findings| says which: it prints a line beginning "read by a person" on exactly those rounds and on no others.
That line is what decides it. Do not work it out for yourself from the state or from any count, because a round you decide wrongly about is a round nobody renders and nobody notices was not rendered.

Where it says so, do not paste the findings into the terminal.
Write a document for the round and serve it: |agora serve --document ./round.json| prints a loopback address.
Hand the person that address, then say briefly what it holds.
The page is served on this machine and reaches nowhere beyond it. A findings document quotes the code it was written about, so it goes to no other machine and is put in front of nobody else.

The person answers on the page.
Every control it carries names the command it becomes, and nothing is recorded until you run that command here.
An option whose command you cannot name is one you do not offer, because a person choosing it would be deciding something whose destination they cannot see.
Wait for them with |agora poll|, which prints what they chose and anything they wrote, and record it with |socrates respond| in their own words rather than in yours.
Do not ask them here for a decision they have already made there.

The page and this terminal have to be matchable in both directions.
Give the surface the round as |socrates findings --json| reported it, so every finding keeps the identifier, the severity and the audience socrates gave it and the person can say f3 and you know which finding they mean.
Each decision on the page names its finding and says what answering it looks like here.
A page that renamed a finding to something more readable hands the person a search problem in exchange for a nicer sentence.

Hold the page to this:

` + PageBar + `

The findings a person answers go on the page in full.
The ones you are patching yourself go on it as a short list, so the person can see the whole round without reading it.`

	loopHeading = "## The loop"

	theLoop = loopHeading + `

1. Implement what the person asked for.
2. Start the review, passing what the person asked for in their own words:

       socrates review --intent "count the characters of a document"

   It returns as soon as the run exists and prints the run identifier, which every later command takes as its argument.
3. Wait for the review to end: poll |socrates status| every five seconds or so until it reports one of the states below.
   Poll rather than follow what the reviewer is saying, because everything you read lands in your own context and you are about to need that context for the patching.
   |socrates status| prints the same handful of lines whether the reviewer was silent or said a great deal.
   Do not set a timeout of your own. socrates gives the reviewer a deadline, a reviewer that answers nothing before it passes ends the run as failed, and so every run reaches one of the states below without you bounding the wait.
4. Act on the state the run ended in. Each of them is below, and a run rests in no other.
5. Report what happened. Say how many rounds ran and what the last one found, not just the verdict.
   Two rounds ending with nothing blocking and one round ending the same way are different results, and the verdict alone does not say which of them the person got.

### clean

An independent reviewer read the change and found nothing blocking.
Say so, say what was reviewed, and stop.
It means nobody found a problem, and not that the tests pass: socrates runs none.

### failed

Say that no review completed, and why. |socrates status| prints the reason.
Nothing has been checked, so this is not a clean result and not a pass.
Report the change as unreviewed, and offer to start another review.

### escalated

Read |socrates findings| and hand every finding to the person, as a page.
Draft it here, in this session: the run is over, so nothing follows that would re-read the drafting, and a fresh subagent would cost its own start-up for no benefit.
Serve it with |agora serve|, hand the person the address it prints, and wait on |agora poll|.
These are theirs to decide, and socrates has nothing further it can do about them.
A run escalates for one of two reasons, and |socrates status| prints which: a finding holds the change while its own action asks for nothing, so no answer of yours addresses it, or the last permitted round still found something blocking.
No further round is possible. A run has two rounds at most, and no flag gives it a third.
Looking again means a fresh |socrates review| of the code as it now stands, and the reviewer of a new run is told nothing about this run, nothing about what it found and nothing about the answers recorded against it. It reads the change cold.
Offer that, and start one only where the person asks for it.

### awaiting-response

Read |socrates findings|. It says of every finding whether it can be patched or needs a person, and that line is what splits the work in two.

A finding it reports as one that can be patched is yours. Patch it, or say why it is wrong, one finding per command:

       socrates respond --finding f1 --fixed --note "an empty document now answers minus one"
       socrates respond --finding f2 --rejected --reason "there is no retry loop in this function at all"

A finding it reports as needing a person is the person's. Put it in front of them as a page, wait for them, and record what they decided in the words they used:

       socrates respond --finding f3 --decided --reason "the endpoint stays, because the migration reads it"

Delegate the drafting of that document to a subagent rather than writing it here.
You are about to patch, record one answer per finding, continue, and read another round, and everything you add to this session now is read again on every one of those turns.
Give the subagent the run identifier and let it read |socrates findings --json| itself, so what comes back to you is the document and a sentence.
Serve that document yourself with |agora serve|, hand the person the address, and wait on |agora poll|.
A subagent that served it would exit owning the surface you are the one waiting on.

You do not answer such a finding on your own behalf. It asks what the change is for, you are the one who wrote the change, and a recorded decision that no person made is a claim a later reader will find and a reviewer will judge.

Once every finding that takes an answer has one, run |socrates continue|, and go back to waiting on |socrates status|.
Do not ask the person whether the round is worth spending. A patch is unreviewed new code you wrote just after being told you were wrong, which is the code most worth reading, and the round is what the findings earned.
The next round is a fresh reviewer that reads your answers and judges them, and there are two rounds at most.

### cancelled

Somebody stopped the review, so nothing was decided about the change.
Say that, and do not report the change as reviewed.
Whoever cancelled it did so on purpose, so ask before starting another one.`

	prohibitions = `## What socrates will not do for you

socrates patches nothing.
Every fix is yours, made in your own working tree, and recorded afterwards.

There is no way to tell the reviewer about the change.
It gets the diff, the repository as it was frozen, and the one line of intent you passed to |socrates review|.
That is the design and not a limitation to work around: an account written by the session that wrote the code is the thing an independent reviewer exists to do without.

A rejection is a claim, and it goes to the next reviewer, which judges it.
It is not a way to make a finding go away, and a rejection without a reason is refused.`

	noSuchBinary  = "the skill names %s %s, and this repository builds no binary called %s"
	noSuchCommand = "the skill names %s %s, and %s offers no such command"
	noSuchFlag    = "the skill names %s %s %s, and that command offers no such flag"

	theTargetCannotBeWritten = "socrates could not write the skill to %s: %s. Name a directory you can write with socrates skill --to"
	noHomeDirectory          = "socrates installs the skill under your home directory, and it could not find one: %s. Name where to write it with socrates skill --to"
)
