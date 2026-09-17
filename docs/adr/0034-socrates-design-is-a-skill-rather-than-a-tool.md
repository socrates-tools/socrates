# 34. socrates design is a skill rather than a tool

- Status: accepted
- Date: 2026-09-17
- Part of: socrates-design

## Decision

`socrates design` is a skill the binary ships and a namespace holding one operation, `socrates design skill`, which installs it.

Review is a tool because a round spawns a model process against a frozen snapshot, and a run outlives the command that started it.
None of that is true of a design.
The interview is the session itself, the machinery is agora, and nothing runs between one round and the next.
So there is no run store, no state machine, no retention policy and no state to supervise.

What the skill still needs from the binary is the check that it names only commands the binary offers, which is why it ships from `design/skill` rather than as a loose file.
That check is the reason a skill belongs to a binary at all, and it applies here exactly as it applies to review's two.

A consequence taken deliberately: because nothing is stored, a design that loses its page loses its address and no work, since every round's document is composed fresh from the session.

## Rejected

- A namespace with sessions, records and retention, on the model of review.
  Its one real pull was a verbatim record of what the person answered, readable back when a spec is later disputed.
  That record was bought with a store, a state machine and a retention policy invented for a workflow that spawns no process, and the ADRs already carry the part of it that matters, which is what was decided and what it beat.
  It stays additive: if the verbatim record proves worth having, the skill is regenerated and the option commands change shape.
- A personal skill file, driving agora over the CLI with socrates untouched.
  Nothing would check that the invocations it writes are ones the binary offers, which is the failure the suite has a gate against.
