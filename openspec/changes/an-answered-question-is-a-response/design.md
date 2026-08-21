## Context

`dialogue.After` is the most load-bearing rule in socrates and the one with the strongest test: its whole input space, not a sample.
Today an `ask-user` finding dominates that rule and ends the run, which is what makes the author's decision unreachable.

The reason recorded beside the rule is that a machine must not argue with intent.
This change does not weaken that reason.
It separates two things the rule currently conflates: *nothing in the dialogue can close this finding* and *nobody can*.
A person can, and the dialogue already has a shape for what a person says about a finding, because a rejection is exactly that: a claim, recorded verbatim, judged by the next reviewer.

What socrates cannot do is verify that a person, rather than the session, answered.
It never could for a rejection either, and the design's answer there is not to try: record the claim, label who made it, and hand it to a reviewer that owes it nothing.

## Goals / Non-Goals

**Goals:**

- A question the reviewer asked can be answered, and the answer is judged.
- Answering a question is visibly different from claiming a patch, in the recorded response and in the next round's prompt.
- A caller that must stop and ask a person can tell from `status` or `--json` without reading prose.
- `escalated` keeps one meaning: socrates is finished and nothing more it can do.

**Non-Goals:**

- Proving a human answered. Unprovable, and the corpus of rejections already shows the design's alternative.
- Raising the round cap. A question answered at round one gets exactly one judging round, which is what the cap already allows.
- Letting socrates apply any patch. The author still writes every line, and this change adds no writer.
- A new exit code.

## Decisions

### An `ask-user` finding parks the round rather than ending the run

`After` returns `awaiting-response` for a round whose blocking findings include ones needing a decision, and `Answerable` admits them.
`escalated` is then reached by the round cap and by failure, which is what it already meant everywhere else.

The alternative is what ships today: end the run and let the author start a fresh review after deciding.
It loses the linkage that makes the second look worth anything.
The new reviewer reads patched code with no knowledge that its predecessor asked whether the change should exist, so nothing ever judges the answer, and the author's own account of the decision reaches nobody.

### A decision is a third resolution, not a reuse of `--fixed`

`scrat respond --decided --reason "<what was decided and by whom>"` answers a finding actioned `ask-user`.
`--fixed` and `--rejected` are refused for such a finding, and `--decided` is refused for one that is not.

Reusing `--fixed` was the cheaper option and it is the one that breaks the design.
A session that wants to move on would answer an intent question by asserting a patch, the response document would record a fix, and the next reviewer would be told a mechanical change closed a question about purpose.
Separating the verbs means bypassing the rule requires stating something false rather than choosing the loosest available word, and it gives the prompt a label to carry.

### The exit code stays `4`, and the distinction lives in the report

A parked run reports how many unanswered findings need a person, in the lines a person reads and in the JSON object.
The exit-code table is unchanged.

A sixth code was considered and rejected on two grounds.
The table is a five-entry contract that a later hook layer keys on, and `4` is factually correct for a parked run: it is not terminal, and something is expected to happen next.
More decisively, a round can hold both a patchable finding and a question, and a single code would have to choose which one to report; the counts report both without choosing.

### The prior-round block labels a decision as a decision

The next round's prompt already carries prior findings and the author's responses in one labelled block, with the preamble that everything in it is a claim rather than evidence.
A decision joins it as a third kind of claim, attributed to a person, so a reviewer that disagrees with the decision can say so and a reviewer that finds the decision unimplemented can say that instead.

The reason is model output's mirror image: text written by the session under review, re-embedded in a prompt, so it goes through the same per-field neutralisation every other such field does.

### The skill states what the terminal state used to enforce

Until now a session could not answer an intent question because the run was over.
After this change nothing mechanical stops it, so the skill's parked branch splits: findings that can be patched are the session's, findings that need a decision are the person's, and the session presents and waits.

This is a real loss of enforcement, stated rather than hidden.
What replaces it is that the bypass is now a false claim in a durable record that a fresh reviewer reads, which is the same protection the project accepts for rejections.

## Risks / Trade-offs

- **The guard against a machine answering intent weakens from a state machine to a flag and a sentence.** → `--fixed` and `--rejected` are refused on a finding needing a decision, so the bypass requires `--decided` and a reason that a fresh reviewer judges; the skill says whose the finding is and a test asserts it says so.
- **BREAKING: a caller that branched on `1` for an intent question now sees `4`.** → `review --wait` already returns `4` for parked runs, so the shape is not new; the counts say what is waited on; README, the skill and a release note state it.
- **A round can now wait forever on a person who never answers.** → Already true of patchable rounds: silence does not advance a round, and retention sweeps the run. No new failure mode, and no new process is left running.
- **`continue` becomes stricter.** A round holding an unanswered question refuses to advance, which is the old rule's protection in a new place: the loop cannot walk past an intent question, it can only wait for one.
- **The corpus reclassifies.** Every recorded answer holding an `ask-user` finding moves from `escalated` to `awaiting-response` in `internal/corpus/outcomes/distribution.txt`. That is a diff somebody reads, and it turns `awaiting-response` from an outcome only fixtures reach into one that real reviewers reach, which is the report's whole purpose.

## Migration Plan

No stored run changes shape, and no run in flight is read differently: the routing rule is applied when a round completes, so runs that already ended keep their recorded state and verdict.
A `decided` resolution is additive in `responses.jsonl`, which is append-only and already tolerates a resolution it does not recognise being absent from older files.

Rollback is the rule and the flag: reverting `After` and `Answerable` restores the old routing, and a `decided` response recorded in the meantime becomes a response to a finding the rule no longer parks on, which the report already handles by counting only what `Answerable` names.

## Open Questions

- Should `--decided` be accepted for a finding actioned `auto-fix`, where the person decides to leave it unpatched? Today that is `--rejected` with a reason, which is honest, and this change keeps the mapping strict rather than guessing.
- Should a decision recorded at the round cap escalate with the decision attached, or fail? The cap's existing behaviour escalates with everything kept, and this change does not alter it, but the wording of the verdict for that case is worth reading once written.
