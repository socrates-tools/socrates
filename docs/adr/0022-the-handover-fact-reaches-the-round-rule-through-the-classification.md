# 22. The handover fact reaches the round rule through the classification

- Status: accepted
- Date: 2026-09-09
- Amends: ADR 7
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

ADR 7 left `dialogue.After` untouched and had `dialogue.AwaitsAPerson` gain the handover fact.
That is one clause short.
`After` escalates a round holding a blocking finding whose own action asks for nothing, because no response addresses such a finding.
On a change nobody here can patch, one does: forwarding it.
A run that escalated at its first round is not awaiting responses, so no disposition can be recorded against it at all, and the finding most worth raising with the author is the one class of finding that could not be raised.

So the handover fact reaches the rule through its input rather than through its clauses.
`dialogue.Classified` reads a document into a classification, and where nobody here can patch the change every finding of it is a person's.
`findings.ClassifyForAPerson` is what it returns there.

`dialogue.After` keeps its signature and its three clauses, and its exhaustive test is unchanged.
`dialogue.AwaitsAPerson`, `dialogue.Answerable` and `dialogue.EscalatedByTheCap` need no handover fact at all, because a classification whose findings all need a person already answers each of them.
One round is therefore derived twice over: from the dispositions, and from a rule nobody added a case to.

What `status` counts changes with it. Such a run reports every finding as one for the person and none as patchable, which is what a run nobody here can patch holds.

## Rejected

- passing the authorship to `After`, `AwaitsAPerson`, `Answerable` and `EscalatedByTheCap`, which is the same fact stated four times; and leaving the escalation, which loses the findings it matters most to forward.
