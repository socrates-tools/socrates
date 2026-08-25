## Context

The first page a real run served was read by the person it was built for, and it failed in five places at once.
Four of them are the specification being followed rather than broken, so the specification is what moves.

The one thread through all five: every decision about that page was made by asking what the session composing it needed to be sure of, and none was made by asking what the person reading it was about to do.
The session needs the command to be exact, the identifier to match, the severity to be copied and the prose to be faithful to the reviewer.
The person needs to know what they are choosing between.

## Goals / Non-Goals

**Goals:**

- Nothing on the page requires the person to know that a terminal exists.
- A decision is named in words, and the return to the terminal stays one exchange.
- The bar the page is held to is specific enough that a page failing it fails visibly.
- The two rendering defects, the hidden element that paints and the borrowed palette, are fixed by rules rather than by care.

**Non-Goals:**

- Changing what the caller receives. `poll` returns the same answer, with the same chosen labels and the same commands.
- Letting a caller rename a decision to a title of its own. That is the failure the copy rule exists for, and it is unchanged.
- Teaching the surface anything about what the words on it mean. It renders a document; the caller writes one.
- Any new asset. The palette and the lettering are the ones already embedded.

## Decisions

**The command travels, and it is not rendered.**
`a-control-names-the-command-it-becomes` permitted a control exactly where it named the command it becomes, so that nobody decides something whose destination they cannot see.
The premise stands and the rendering was the wrong conclusion drawn from it.
The destination a person needs to see is what the answer does to the change, said in words. The exact argument vector is what the *caller* needs, and the caller gets it back from `poll` verbatim, which is where it was always going to be used.
So the rule keeps its force and loses its display: an option whose command the session cannot name is still one it does not offer, because that is a session that does not know where the answer goes.
What replaces the command under each option is nothing. The option's own sentence is the whole of what the person needs.
The same reading applies to what the page renders around the evidence.
A diff heading read `a/src/api/windows.py` beside `b/src/api/windows.py`, which is the patch format naming the two inputs of a tool about to apply it.
A person knows that file as `src/api/windows.py`, once, so that is what it says.

Alternative rejected: rendering the command behind a disclosure control. A control the page offers is a control somebody opens, and the page would still be teaching terminal syntax to a person who is not in one.

**A name is copied like an identifier, not accepted like prose.**
`f1` is what the reviewer calls a finding. It is a key, and keys are for looking things up.
The obvious fix, letting the caller write whatever heading it likes, is exactly the failure `The surface copies the identifiers rather than accepting them` exists to prevent: a page about "the endpoint the migration reads" leaves the person unable to say which of six findings they mean.
So the name enters through the same door as the severity: the source entry a decision cites carries it, and the surface copies it onto the decision.
A caller that writes a different name on the decision itself loses to the source, mechanically.
What stops the name being an invention is stated where the session that writes it reads: derive it from the identifier and never title it. `f1` becomes "Finding 1", and a person saying "finding 1" and a session running `socrates respond --finding f1` are talking about the same thing without either of them translating.
The surface itself learns nothing from this. It does not know what a name is for, and could not: it may not say the word.

**The bar gains the sentence that would have caught the prose that was shipped.**
The bar already said plain English, and the page shipped "the smoke assertions on the 25 cap come out".
That passed because "plain English" is a standard a writer grades themselves against, and a session that has just read the reviewer's findings believes the reviewer's vocabulary is plain.
The addition is a test the writer can actually fail: an option says what the person gets if they choose it, in the words of the thing being built, and a sentence a person could not act on without opening the code is not plain English.
Naming a symbol, a constant or a test as though the reader knows it is the specific move being ruled out, because it is the one that happened.

**A hidden element paints nothing, and that is a rule rather than a discipline.**
`.committed` is hidden until an answer is committed, and the stylesheet gives it `display: grid`, which outranks what `hidden` does.
The result was an empty bordered box beside the commit control that a person clicked and nothing happened.
The two halves of the bug are in different files and neither is wrong on its own, so the fix is a rule at the top of the stylesheet that no later declaration can beat, plus a check that renders the page and looks for anything hidden that still occupies space.
Alternative rejected: removing `display: grid` from that one rule. It fixes today's instance and leaves the next one to be found by a person again.

**The palette is quarried, not borrowed.**
The page was cream paper and rust ink, which is the palette of the tool that composed the document, and a person reading it said so within a minute of seeing it.
The product is named for a stonemason's son, so the palette is stone: the white of cut marble for the page, the grey of the block for its rules and quiet ground, slate for the ink, and the blue of the sea between the islands for the one accent that carries a decision.
The settled ground is olive rather than green-grey, because that is what grows on the hill the marble came from.
Small labels are set as an inscription is: capitals, letterspaced, cut into the stone rather than written on it.
None of this is enforceable by a test and none of it is left to care either: what the test holds is that every colour is a token stated once, and that every text on the page meets a contrast ratio against what it is painted on, under both preferences.
A palette that is stated once can be replaced in one place; a contrast floor is what stops "quiet" from becoming "unreadable", which is the direction a stone palette fails in.

## Risks / Trade-offs

**A person can no longer read the identifier off the page.**
They read "Finding 1" and the session holds `f1`.
The mitigation is that the name is derived from the identifier rather than invented, so the mapping is obvious in both directions, and the person's answer travels back through `poll` keyed by the identifier anyway.
The case this costs is a person reading the page and typing `socrates respond` themselves, which is not the loop the skill teaches and never was.

**The contrast floor is a gate that can fail on a palette somebody prefers.**
That is the point of it. The floor is stated in the specification rather than in the stylesheet, so raising or lowering it is a decision somebody makes on purpose.

## Migration Plan

None. Nothing persists between exchanges, no document already written stops validating, and a document with no name on its source entries renders under its identifiers exactly as it does today.
