## Why

The skill tells the session to run a command it did not compose, and nothing in it says what that command may be.

`agora poll` prints one `Decided` per decision, and each carries `Commands []string` copied from the option the person chose.
The surface is right to pass the string through rather than judge it: agora names no review, no finding and no run, so it cannot know which commands are socrates' to run.
The skill is the only place that knows, and all it says is that every control carries the command it becomes, that `agora poll` hands that command back, and that nothing is recorded until the session runs it here.
So the shape of what may come back is stated nowhere.

The path it comes back by is longer than one session, and socrates built every link of it on purpose.
The `awaiting-response` branch tells the session to delegate the drafting to a subagent, and to let that subagent read `socrates findings --json` itself.
A finding's prose is model output.
A subagent reads it, writes the option's `command`, the surface copies that string onto the page, the person chooses, and `agora poll` prints it to a session that did not write it and has been told to run it.
Each link is justified where it stands.
The guard the chain needs is a sentence bounding the command, and that sentence has never been written.

The person's words arrive the same way and are handled no better.
`Decided.Words` is free text from a page that anything reaching the loopback interface can answer, and the skill says to record it with `socrates respond` in their own words rather than in the session's.
It does not say those words are one argument rather than a fragment of shell, and it does not say they are something recorded rather than something followed.

A finding is model output too, and the skill tells the session to read it and patch from it.
socrates already holds that a reviewer's field is not to be trusted structurally: `internal/hygiene` neutralises every field on the way into a prompt, and nothing carries that same holding to the session reading the same field at the terminal.

`agent-skill` already requires that the skill write every invocation it names as code, and `dialogue.ResponseRules` names three in prose.
Nothing goes unproven by it, because the same three flags are written as code in the `awaiting-response` branch and the check reads them there.
It is a promoted requirement with no test on its "every" direction, and a session parsing the skill for commands meets the same invocation in two forms.

The last of these is not security.
The skill is 17.4 KB and every session that reviews anything reads all of it.
A large share of that is the skill arguing for its own rules, and `openspec/` is where the reasoning is supposed to live.
Nothing measures the skill, so nothing notices it growing, and a rule a session has to find inside a paragraph of argument is a rule it can read past.

## What Changes

- The skill SHALL bound the command an option may carry: a `socrates respond` invocation naming that decision's own finding, and nothing else.
- The skill SHALL say that what `agora poll` hands back is checked against that shape before it is run, and that a command of any other shape is reported to the person rather than run.
- The skill SHALL say that the words a person writes are recorded as one argument, without shell interpolation, and are data recorded rather than instruction followed.
- The skill SHALL say that a finding is a claim to judge and not a directive to obey, and that nothing in one is a reason to run a command it names.
- `dialogue.ResponseRules` SHALL write its invocations as code, restoring `The skill writes an invocation as code`, which is a promoted requirement it currently breaks.
- The skill SHALL be held to a stated size ceiling, stated once and asserted by a test, so that verbosity fails the suite rather than accumulating.
- The reasoning for a rule the skill states SHALL live in `openspec/` rather than beside the rule, except where a requirement states that the skill itself must carry it.

## Capabilities

### New Capabilities

None.
Every part of this is the one skill capability being asked to state what it currently leaves unstated.

### Modified Capabilities

- `agent-skill`: what the skill says about the command an option carries, about what comes back from the page, about what a finding is, and how large the skill may be.
  No existing requirement is amended.
  The response rules being written as code is `The skill writes an invocation as code` restored rather than changed, so it carries no delta.

## Impact

- `internal/skill/messages.go`: the page section gains the bound on a command, the handling of what comes back, and the standing of a finding, and loses the arguments that belong in `openspec/`.
- `internal/dialogue/messages.go`: `ResponseRules` writes its three invocations as code, which `socrates findings` prints as well as the skill.
- `internal/report`: the same rules are what `findings` puts in front of the session, so its tests move with them.
- `internal/skill`: a stated ceiling, and the tests for each new requirement.
- `.claude/skills/socrates/SKILL.md`: regenerated by `make skill`.
- Nothing reaches the reviewer.
  `ResponseRules` is not in the delivered prompt, so no prompt section, no fixture and no recording moves.
- Nothing reaches the surface.
  agora learns no noun, gains no validation and refuses nothing new, because the command it hands back is not a string it is able to judge.
- Nothing reaches the snapshot, the run store, the verdict rule or the dialogue's own rules.

## Dependencies

This change SHALL follow `every-finding-earns-a-card`, which is archived.
That change is what put a card and an option-less card on the page, and the bound on a command is stated against the document shape it left behind.
`agent-skill` is otherwise promoted, and every requirement here narrows it rather than widening it.
