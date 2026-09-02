---
name: socrates
description: Have an independent reviewer read the change you just made, and drive that review to a verdict. Use after implementing anything nobody else has read.
---

# Have somebody who did not write it review it

You have just implemented something and nobody has read it.
`socrates review` freezes your working tree into a commit nothing can move and hands it to a reviewer that wrote none of it.
You drive that review to a verdict, and every patch along the way is yours to make.
What the reviewer says while it works is not for you: if the person wants to see it happen, suggest that they run `socrates watch` in another terminal.

## How to answer a finding

- You did not review this code. An independent reviewer did, and it did not write any of it.
- For each finding, either apply a fix and record it with socrates respond --fixed, or record a rejection with socrates respond --rejected and a concrete reason.
- A rejection is a claim you are making to the next reviewer, which judges it. It is not a decision, and it does not make the finding go away.
- A finding that asks a person is the person's to answer and not yours: take it to them, and record what they decided with socrates respond --decided and a reason.
- Do not resolve a finding by weakening or deleting behaviour you added on purpose. Where the finding is right about a symptom and wrong about the cause, say so and fix the cause.

## The page a person reads

Every round that reached a result is handed over as a page, so do not paste the findings into the terminal.
Write a document for the round and serve it: `agora serve --document ./round.json` prints the loopback address of its page.
Hand the person that address, then say briefly what it holds.
The page is served on this machine and reaches nowhere beyond it. A findings document quotes the code it was written about, so it goes to no other machine and is put in front of nobody else.

Some rounds hold something only a person can act on, and those are the ones you stop at.
`socrates findings` says which: it prints a line beginning "awaits a person" on exactly those rounds and on no others.
That line is what decides it. Do not work it out for yourself from the state or from any count, because a round you decide wrongly about is one you stop at for an answer nobody was asked for, or one you patch around the person.
Where it says so, wait for them. Where it says nothing, hand over the address and get on with the work, because a page nobody is being asked to answer is a page nobody has to open before you carry on.

The document is one decision per finding a person answers, each citing the entry it came from:

    {"title": "...", "preamble": "...",
     "source":    [{"id": "f3", "name": "Finding 3", "severity": "error", "audience": "needs a person", "location": "src/api/windows.py:41"}],
     "decisions": [{"id": "f3", "question": "...", "meaning": "...", "mode": "single-choice",
                    "evidence": {"language": "python", "code": "...", "patch": "..."},
                    "options": [{"label": "...", "meaning": "...", "command": "socrates respond --finding f3 --decided --reason \"...\"",
                                 "recommended": true}]}]}

Where the round holds nothing a person answers, the document holds no decision at all: a title, a sentence or two, and the source.
Serve that the same way and hand over the address. The page says for itself that there is nothing on it to answer and that nobody is waiting on it.
Do not invent a decision to fill such a page. An option that records nothing is one a person can spend their attention wondering about.

One option of a decision may carry `"recommended": true`, and at most one: a document marking two is refused before anything is served.
The page reads that option first and says on it, in words, that it is the one you would take.
It selects nothing, so a decision nobody touched comes back with nothing chosen, and what to do about the finding is still theirs.

The person answers on the page.
Every control carries the command it becomes, `agora poll` hands that command back to you, and nothing is recorded until you run it here.
The page shows that command to nobody: a person deciding is not at a prompt, a command under an option is the same choice again in a language they did not come here to read, and one they may believe they have to type.
An option whose command you cannot name is one you do not offer, because you would be offering a decision you cannot honour.
Where you are waiting for them, do it with `agora poll --session`, giving it the session at the end of the address `agora serve` printed, because the surface holds a page per round and a poll naming none is told about none of them.
Keep that session for the whole review. A later round is written to the page they already have, by serving it with `agora serve --session`, and you hand over no second address: one review is one change and one person reading it, so a second address is a second tab whose first has gone stale with nothing on it saying so.
It prints what they chose and anything they wrote, and you record that with `socrates respond` in their own words rather than in yours.
Do not ask them here for a decision they have already made there.

The page and this terminal have to be matchable in both directions.
Give the surface the round as `socrates findings --json` reported it, so every finding keeps the severity and the audience socrates gave it.
Name each decision from its identifier and never title it yourself: f3 is "Finding 3", so a person saying finding 3 and you running `socrates respond --finding f3` are naming one thing.
The surface takes the name, the severity and the audience from the source you cite rather than from the decision you wrote, so a decision you titled yourself loses to it.
A page that titled a finding something more readable hands the person a search problem in exchange for a nicer sentence.
Each decision says what choosing each option means for the change, and nothing about what you will type.

Hold the page to this:

- Get to the point.
  Length is not thoroughness, and a person who has to search the page for their own decision has been handed the terminal output again.
- Say in plain English what the finding means for the change.
  Do not restate the reviewer's technical prose in a nicer typeface.
- Say what each option gets the person, in the words of the thing being built rather than the words of the code.
  A sentence they could not act on without opening the code is not plain English, and a symbol, a constant or a test named as though they already know it is that sentence.
- Say which option you would take, where you would take one, by marking that option and no other.
  You have read the finding and the person has not, and a list where every option looks equally likely makes them work out what you already know.
- Where a finding is purely technical, show a small example of it going wrong.
  A reader should not have to work out for themselves what the term costs them.
- Where what a finding costs is what a person sees, show them it rather than describing it.
  They already know the screen and are being asked to recognise it, and a panel retyped in a monospace block is the one form they cannot recognise it in.
- Draw a diagram only where a workflow is what makes a finding hard to follow, and not otherwise.
- Spend tokens as though they were the person's, because they are.

A decision may also carry one illustration.
It travels beside the evidence rather than inside it, in a field of its own:

    "fragment": {"markup": "<div class=\"panel\"><div class=\"row\">...</div></div>"}

An illustration takes one of two forms.
A drawing you write yourself as inline SVG is for a workflow: the surface renders the drawing you supplied and lays out nothing of its own.
Markup the document's stylesheet paints is for a screen the person already knows, because what makes them recognise a panel is its structure, its controls and the words actually on it.

The look is stated once for the whole document rather than once per decision:

    "stylesheet": ".panel{border:1px solid #d8d2c4;padding:0.75rem}.row{display:flex;justify-content:space-between}"

One document is one product, and one stylesheet paints every illustration in it.

Both the markup and the stylesheet are reduced to an allowlist before they are embedded, so do not write what will be dropped.
Out of the markup go script, event handlers, remote references and form controls.
Out of the stylesheet go imports, font faces, every property the allowlist does not name, and every value naming a resource.

The surface scopes the stylesheet itself, and every selector it keeps is confined to the illustration.
Write your selectors for the names your own markup carries, leave the scoping to the surface, and do not expect a rule you aimed at the page to reach the page.

The findings a person answers go on the page in full.
The rest of the round goes on it as a short list, so the person can see the whole of it without reading the whole of it.
That list is the `source` you cite rather than prose you write: the surface renders every entry of it, whether or not a decision cites that entry, so give it every finding `socrates findings --json` reported and retype none of them.

## The loop

1. Implement what the person asked for.
2. Start the review, passing what the person asked for in their own words:

       socrates review --intent "count the characters of a document"

   It returns as soon as the run exists and prints the run identifier, which every later command takes as its argument.
3. Wait for the review to end: poll `socrates status` every five seconds or so until it reports one of the states below.
   Poll rather than follow what the reviewer is saying, because everything you read lands in your own context and you are about to need that context for the patching.
   `socrates status` prints the same handful of lines whether the reviewer was silent or said a great deal.
   Do not set a timeout of your own. socrates gives the reviewer a deadline, a reviewer that answers nothing before it passes ends the run as failed, and so every run reaches one of the states below without you bounding the wait.
4. Act on the state the run ended in. Each of them is below, and a run rests in no other.
5. Report what happened. Say how many rounds ran and what the last one found, not just the verdict.
   Two rounds ending with nothing blocking and one round ending the same way are different results, and the verdict alone does not say which of them the person got.

### clean

An independent reviewer read the change and found nothing blocking.
Hand the round over as a page all the same, because a result reads better on a page than in a scrollback whether or not anybody is being asked to answer it.
Draft it here, in this session: nothing follows a run that is over, so the drafting is paid once and a fresh subagent would cost its own start-up for no benefit.
Serve it with `agora serve`, and hand the person the address it prints unless an earlier round of this run already gave them one: then this is written to the page they have.
Nothing on it is theirs to answer, so `socrates findings` says nothing about waiting and you wait on nothing: say what was reviewed, and stop.
It means nobody found a problem, and not that the tests pass: socrates runs none.

### failed

Say that no review completed, and why. `socrates status` prints the reason.
Nothing has been checked, so this is not a clean result and not a pass.
Report the change as unreviewed, and offer to start another review.

### escalated

Read `socrates findings` and hand every finding to the person, as a page.
Draft it here, in this session, for the same reason: nothing follows an escalated run either.
Serve it with `agora serve`, hand the person the address it prints unless an earlier round of this run already gave them one, and wait on `agora poll`, which is what `socrates findings` says to do about this round.
These are theirs to decide, and socrates has nothing further it can do about them.
A run escalates for one of two reasons, and `socrates status` prints which: a finding holds the change while its own action asks for nothing, so no answer of yours addresses it, or the last permitted round still found something blocking.
No further round is possible. A run has two rounds at most, and no flag gives it a third.
Looking again means a fresh `socrates review` of the code as it now stands, and the reviewer of a new run is told nothing about this run, nothing about what it found and nothing about the answers recorded against it. It reads the change cold.
Offer that, and start one only where the person asks for it.

### awaiting-response

Read `socrates findings`. It says of every finding whether it can be patched or needs a person, and that line is what splits the work in two.

A finding it reports as one that can be patched is yours. Patch it, or say why it is wrong, one finding per command:

       socrates respond --finding f1 --fixed --note "an empty document now answers minus one"
       socrates respond --finding f2 --rejected --reason "there is no retry loop in this function at all"

A finding it reports as needing a person is the person's. Put it in front of them, wait for them, and record what they decided in the words they used:

       socrates respond --finding f3 --decided --reason "the endpoint stays, because the migration reads it"

Hand the round over as a page whichever of the two it holds, and delegate the drafting of that document to a subagent rather than writing it here.
You are about to patch, record one answer per finding, continue, and read another round, and everything you add to this session now is read again on every one of those turns.
That is as true of a round you answer alone as of one holding a question, because what costs you is the finding prose rather than who answers it.
Give the subagent the run identifier and let it read `socrates findings --json` itself, so what comes back to you is the document and a sentence.
Serve that document yourself with `agora serve` and hand the person the address, unless an earlier round of this run already gave them one.
A subagent that served it would exit owning the surface you are the one waiting on.
Then do what `socrates findings` said about this round: where it says the round awaits a person, wait on `agora poll` before recording anything; where it says nothing, get on with the patching and leave the page for them to read whenever they like.

You do not answer such a finding on your own behalf. It asks what the change is for, you are the one who wrote the change, and a recorded decision that no person made is a claim a later reader will find and a reviewer will judge.

Once every finding that takes an answer has one, run `socrates continue`, and go back to waiting on `socrates status`.
Do not ask the person whether the round is worth spending. A patch is unreviewed new code you wrote just after being told you were wrong, which is the code most worth reading, and the round is what the findings earned.
The next round is a fresh reviewer that reads your answers and judges them, and there are two rounds at most.

### cancelled

Somebody stopped the review, so nothing was decided about the change.
Say that, and do not report the change as reviewed.
Whoever cancelled it did so on purpose, so ask before starting another one.

## What socrates will not do for you

socrates patches nothing.
Every fix is yours, made in your own working tree, and recorded afterwards.

There is no way to tell the reviewer about the change.
It gets the diff, the repository as it was frozen, and the one line of intent you passed to `socrates review`.
That is the design and not a limitation to work around: an account written by the session that wrote the code is the thing an independent reviewer exists to do without.

A rejection is a claim, and it goes to the next reviewer, which judges it.
It is not a way to make a finding go away, and a rejection without a reason is refused.
