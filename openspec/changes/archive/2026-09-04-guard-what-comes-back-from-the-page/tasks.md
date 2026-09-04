**The suite must be green after every task, not only at the end of a group.**

**The ceiling is stated before the trimming and not after.**
A ceiling written once the skill already fits is a number nobody chose, and a test that passed the first time it ran is a test nobody has seen fail.
Group 5 states the ceiling, watches it fail, and group 6 is the work that makes it pass.

**Nothing in this change reaches the surface or the reviewer.**
agora learns no noun, refuses nothing new and validates no command: the bound on a command is the session's to check and cannot be the surface's.
`ResponseRules` is not in the delivered prompt, so no prompt section is touched, no fixture is re-recorded and `make fixtures` is not run.
A task that adds a refusal to `internal/document`, `internal/surface` or `internal/render` is in the wrong change.

## 1. The response rules write their invocations as code

This group carries no scenario.
It restores `The skill writes an invocation as code`, which is promoted and which `ResponseRules` currently breaks by naming three invocations in prose.

- [x] 1.1 Add a skill test that every `socrates respond` the response rules name is written as code, and verify it fails for the right reason: three are prose
- [x] 1.2 Write the three invocations in `dialogue.ResponseRules` as code, and verify 1.1 passes
- [x] 1.3 Verify `One source, two surfaces` still holds, so the skill's rendering and what `scrat findings` prints are still one text
- [x] 1.4 Verify the `internal/report` and `internal/e2e` tests that assert the rules' text still pass, and update any that pinned the old form
- [x] 1.5 Run `make skill` and commit the regenerated `.claude/skills/socrates/SKILL.md`
- [x] 1.6 Mutate the guard: take the backticks off one invocation and verify 1.1 fails

## 2. The command an option carries has a stated shape

- [x] 2.1 Add a skill test that the page section states the shape of an option's command as a `socrates respond` invocation naming the decision's own finding, and verify it fails for the right reason. Covers `The skill states the shape a command may take`
- [x] 2.2 Add a skill test that the skill says a polled command is checked against that shape before it is run. Covers `The skill says to check what came back before running it`
- [x] 2.3 Add a skill test that the skill says an unmatched command is reported to the person, and names no repair, replacement or fallback. Covers `A command of another shape is not run`
- [x] 2.4 Add a skill test that the shape is stated once, where the document's command field is described, and that no branch restates it. Covers `The bound is stated where the command is named`
- [x] 2.5 State the bound in `internal/skill/messages.go`, and verify 2.1 to 2.4 pass
- [x] 2.6 Verify `The skill says how an answer reaches the run` and `Every command named is real` still pass over the new text
- [x] 2.7 Mutate the guard: widen the stated shape to any `socrates` invocation and verify 2.1 fails

## 3. What a person wrote is recorded rather than followed

- [x] 3.1 Add a skill test that the skill says the person's words reach `socrates respond` as a single argument with nothing interpolated, and verify it fails for the right reason. Covers `The skill says the words are one argument`
- [x] 3.2 Add a skill test that the skill says those words are recorded rather than followed. Covers `The skill says the words are not instructions`
- [x] 3.3 State both in `internal/skill/messages.go`, beside the sentence that already says to record the person's words rather than the session's, and verify 3.1 and 3.2 pass
- [x] 3.4 Verify `A decision is still the person's to make` is untouched: recording their words rather than the session's is still what the skill says to do

## 4. A finding is a claim the session judges

- [x] 4.1 Add a skill test that the skill states a finding is a claim the session judges and that nothing in one is a reason to run a command it names, and verify it fails for the right reason. Covers `The skill states what standing a finding has`
- [x] 4.2 Add a skill test that the standing appears with the response rules and before the step that reads the findings. Covers `The standing is stated with the response rules`
- [x] 4.3 State it in `internal/skill/messages.go` beside the response rules, and verify 4.1 and 4.2 pass
- [x] 4.4 Verify `The rules are stated before the findings are read` still holds with the sentence added

## 5. The skill has a stated ceiling

- [x] 5.1 Record the rendered skill's size with groups 1 to 4 landed, before anything is trimmed
      18,298 bytes, against 17,413 before this change.
- [x] 5.2 Add the ceiling as one named constant in `internal/skill`, and choose the number as a target rather than from the current size. State the number and the reasoning for it in this task
      16,896, which is 16.5 KiB.
      16,384 was tried first and is not reachable: twelve promoted requirements pin their own reasoning inside the skill, and the trimming stalled at 16,508.
      The number is that floor plus two or three sentences, and `design.md` records why it could not be chosen independently.
- [x] 5.3 Add a skill test that fails while the rendered skill exceeds the ceiling, naming the size, the ceiling and that the reasoning belongs in `openspec/`, and verify it fails against the current skill. Covers `A skill over the ceiling fails the suite`

## 6. The reasoning moves to openspec

- [x] 6.1 Read the page section against the requirements over it, and list every sentence that argues for a rule rather than stating one
      Nine moved. Three more were tried and restored, because `A round with nothing to answer is still a round in full`, `Every finding in the round is a card` and `The reason for offering no options is given` each require the skill to carry the reason.
- [x] 6.2 Move each of those into the requirement in `openspec/specs/agent-skill/spec.md` that owns the rule, or into this change's `design.md` where no requirement owns it, and take it out of `internal/skill/messages.go`
- [x] 6.3 Verify the three prohibitions still carry their reasoning, so `prohibitions_test.go` passes untouched. Covers `The reasoning a session needs stays`
- [x] 6.4 Verify every requirement over the skill still passes, one at a time rather than as a suite run, so a rule dropped with its argument is caught as itself
- [x] 6.5 Verify 5.3 now passes, and record the rendered size in this task
      16,508 bytes, down from 18,298 with the three new rules included, and from 17,413 before this change.
- [x] 6.6 Run `make skill` and commit the regenerated `.claude/skills/socrates/SKILL.md`
- [x] 6.7 Mutate the guard: add a paragraph to the skill and verify 5.3 fails

## 7. The whole thing

- [x] 7.1 Drive the loop end to end through the real binaries, as `internal/e2e` already does, and verify a document whose option carries a bounded command is served, polled and recorded
- [x] 7.2 Run `make fmt` and `make verify`
- [x] 7.3 Verify `make fixtures-diff` is unaffected and that nothing in this change reached the real CLI
