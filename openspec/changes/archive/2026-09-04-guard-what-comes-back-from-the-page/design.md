## Context

The skill is the only document in socrates that an agent reads as instructions and acts on without a person in between.
Everything else socrates writes is read by a person: the status lines, the findings, the page.
So the skill is the one place where an unstated rule is not a gap in documentation but a gap in a control path.

Three strings reach the session from outside the session, and the skill currently treats all three as though the session had written them.

The first is `Decided.Commands`, the option's own command string, handed back by `agora poll`.
The surface copies it and never runs it, and the invariant that it never runs one is deliberate.
It also cannot judge it: `internal/surface`, `internal/document` and `internal/render` may name no review, no finding, no verdict, no round and no run, and a check holds that.
A surface that validated `socrates respond --finding f3` would have to know what a finding is, which is the one thing agora is built not to know.
So the bound on that string cannot be a refusal in agora.
It is the skill's or it is nowhere.

The second is `Decided.Words`, whatever the person typed.
The page is on loopback and reaches no other machine, which is the property socrates argues for and the one it holds.
It is not the property that makes the words safe to act on: a person typing a sentence that reads as an instruction is the ordinary case, not an attack, and the session has been told to record their words rather than its own.

The third is a finding.
`internal/hygiene` neutralises every reviewer field on the way into a prompt, so the position that a reviewer's output is structurally untrusted is already socrates', already stated, and already tested.
The session reads the same field at the terminal with none of that said to it.

Alongside those, the skill has grown to 17.4 KB and carries the argument for most of its own rules.
`openspec/` is normative and is where the reasoning lives, and `CLAUDE.md` says the map is not the place to duplicate it.
The skill has been quietly exempt from that, and the exemption has no owner and no ceiling.

## Goals / Non-Goals

**Goals**

- A command that reaches the session from the page has a stated shape, and the session checks it before running it.
- The words a person writes are recorded, not obeyed, and reach `socrates respond` as one argument.
- A finding has a stated standing: a claim the session judges.
- Every invocation the skill names is written as code, which is what the skill is already required to do.
- The skill has a stated size, and exceeding it fails the suite.

**Non-Goals**

- agora does not validate a command. It composes one and hands it back, and it learns no noun to judge one by.
  A later change proposing that agora refuse a command is proposing that agora learn what a finding is, and that is what this non-goal exists to make visible.
- socrates does not run the command for the session.
  A `socrates respond` that took a polled answer and applied it would put the recording inside the tool and leave nothing for the check to sit between.
- The skill does not stop naming a command in an option.
  The command under an option is what makes an answer recordable at all, and the fix is a bound on it rather than its removal.
- The reasoning is not deleted. It moves to `openspec/`, which is where a reader who wants it goes.
- No refusal is added to `internal/document` for the command's shape, for the same reason agora cannot hold it: `internal/document` is the surface's.

## Decisions

### The bound on a command is a shape, and the shape is `socrates respond` for that decision's own finding

The narrowest bound that still lets an option record an answer is: the command is a `socrates respond` invocation, and the finding it names is the finding the decision cites.

Anything wider admits a command the session did not intend, and anything narrower stops an option recording anything.
Naming the decision's own finding is the part that matters most and costs nothing: it is already true of every document the skill teaches, and it closes the case where one option records an answer against a different finding than the card the person was reading.

The check happens twice, and the two are not redundant.
Writing the document is where the session is the author and can simply not write a command it cannot name, which the skill already says.
Reading what `agora poll` printed is where the session is not the author, possibly not even the same session, and that is where a shape check is the only thing standing between the page and a command being run.

### A command of another shape is reported rather than run, and rather than corrected

The session does not repair it, and does not fall back to composing its own.

A command that came back in the wrong shape means the document that was served is not the document the session thinks it served, and that is a fact the person needs rather than one the session smooths over.
Repairing it would also make the check unfalsifiable: a session that fixes what it finds never reports finding anything.

### The person's words are one argument, and the skill says so in the two places it says to record them

The rule is that the words go to `socrates respond` as a single argument with no shell interpolation, and that they are recorded rather than acted on.

The shell half is stated because the words routinely contain the characters that break it.
A person explaining a decision writes quotation marks.
Stating it as "one argument, no interpolation" rather than as an escaping recipe keeps the rule true whatever the session's tool for running a command is.

The standing half is stated because the skill already tells the session to prefer the person's words over its own, which is exactly the instruction that makes their words feel authoritative.

### A finding is a claim the session judges, and the skill says it once

The skill already says a rejection is a claim the next reviewer judges.
This is the same sentence pointed the other way, and it belongs beside the response rules rather than in the page section, because it governs how the session reads a finding at all rather than what it puts on a page.

It stops short of a procedure.
There is nothing to check and nothing to refuse: the rule is what standing the text has, and its whole effect is that a sentence inside a finding is not a reason to do anything the skill does not already say to do.

### `ResponseRules` writes its invocations as code, and that is a defect rather than a decision

`The skill writes an invocation as code` is promoted and says every invocation.
`ResponseRules` names three in prose, so the skill breaks a requirement it is held to.

The requirement's own test reads the check rather than the skill, so it proves that prose about the product yields no invocation and never asks whether an invocation was written as prose.
Nothing goes unproven as a result: the three flags appear as code in the `awaiting-response` branch and the check reads them there.
What it costs is that a session parsing the skill for commands meets the same invocation in two forms.

So this carries no delta.
It is the promoted requirement restored, with the test its unenforced direction never had.

`ResponseRules` is also what `socrates findings` prints, so the backticks land in front of the session at the terminal as well as in the skill.
That is the whole of their reach.
The rules are not in the delivered prompt, so no prompt section moves and no fixture is re-recorded.

### The skill is held to a stated ceiling rather than to taste

A byte ceiling is a crude measure of a document and it is the only one that fails on its own.
The alternatives all read prose: a rule about paragraphs, a rule about arguments, a rule about one sentence per rule.
The repository already knows what those are worth, because `make scenario-coverage` and `make scenario-drop` state plainly that neither reads prose and that a stale sentence is something only a reader catches.

So the ceiling is a number, and the reason it exists is the message its failure prints, because that is where a contributor meets it.

It is not chosen freely, and the reason is worth recording.
Twelve promoted requirements over the skill pin their own reasoning inside it: a session told that a resolved card carries no options without being told that a control changing nothing reads as a question talks itself out of the rule.
Those are the carve-out this change states, and together they are most of what looked movable.
What actually moved is the reasoning no requirement claims, which is about nine hundred bytes.

So the ceiling is the floor that leaves plus a margin of two or three sentences, rather than a target chosen against the skill's will.
That is weaker than it sounds only if the margin is large.
At this size the next paragraph added has to displace one, which is the question a reviewer would otherwise have to remember to ask.

## Risks / Trade-offs

**A shape check the session performs is a check the session can skip.**
This is true of every rule in the skill and it is the reason the bound is stated as a shape rather than as advice: a session that ran a command not matching a stated shape did something the skill names, which is reviewable after the fact, where a session exercising judgement is not.
It is also why the non-goals record what enforcing it in code would cost, so a later change can take that cost knowingly.

**A byte ceiling can be met by making prose worse.**
A session under a ceiling could compress rules into a dense sentence that is shorter and harder to follow.
The mitigation is that the ceiling is not the only requirement over the skill: a dozen existing requirements state what it must say and in what order, and the ceiling can only be met by moving reasoning out, not by dropping a rule.

**Moving reasoning to `openspec/` moves it away from the reader who needed it.**
A session reading the skill does not read `openspec/`.
This is accepted rather than mitigated: the session needs the rule and a person deciding whether the rule is right needs the argument, and the two readers are not the same reader.
Where a requirement finds that the session itself needs the reason, it says so and the reason stays, and `prohibitions_test.go` already holds three of those.

**Backticks reach a person's terminal.**
`socrates findings` prints these rules, and a person reading them sees the marks.
This is accepted: the same four sentences cannot be code in one rendering and prose in the other without becoming two copies, which is the drift the requirement over them exists to prevent.

## Migration

Nothing installed changes shape.
`socrates skill` writes the same file to the same place, and a session holding the previous skill is not wrong about any command, flag or state.
The three new rules narrow what a session may do with a string it did not write, and no document that the previous skill taught becomes invalid.
