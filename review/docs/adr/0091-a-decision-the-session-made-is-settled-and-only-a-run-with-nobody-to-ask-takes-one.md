# 91. A decision the session made is settled, and only a run with nobody to ask takes one

- Status: accepted
- Date: 2026-10-26
- Part of: last-read-review

## Decision

A run begun with `socrates review begin --nobody-to-ask` records that no person is present to answer a finding.
On such a run a finding that asks a person accepts one resolution, `settled`: the session that started the run decided it, and `--reason` says what it decided.
On every other run that finding accepts `decided` alone, as before.
So the record cannot claim either party falsely: `decided` always means a person answered, and `settled` always means nobody did.

A `settled` response is a claim a fresh reviewer judges, like `decided`, and never a suppression.
The second round's prompt names it as the session's decision rather than a person's.

The word is settled because the implement skill already uses it for exactly this act, "Settle everything a capable engineer would decide by your own best judgement".
Review's glossary defines it on its own terms, per `CONTEXT-MAP.md`, and it gives up the meaning design's and agora's glossaries give the same word, a question a person has answered.

The surprise is that a run states who is present when it begins, as `--not-my-change` states who can patch, and that the socrates-review skill, which runs with a person present, never passes it and so can never settle a question that was the person's.

## Rejected

- Keep `decided` and record who decided it with a second flag, such as `--decided --by session`.
  Declined without a word.
- A run marked unattended at `begin`, which reclassifies every finding that asks a person as one answered by a fix or a rejection, with no new resolution.
  Declined without a word.
- No change to review, so a finding that asks a person parks the run and waits for the person at hand-back.
  Declined without a word.
- `settled` accepted on any run that is mine, with the socrates-review skill told in words never to use it.
  Declined without a word.
- The name `judged`.
  Declined without a word.
- The name `assumed`.
  Declined without a word.
