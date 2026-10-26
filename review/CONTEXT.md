# Context: review

`socrates review` puts an independent reviewer in front of unreviewed local changes.
These are review's words, on top of the suite's shared ones in `/CONTEXT.md`.
Use them rather than synonyms.
`docs/adr/` here holds review's decisions and `AGENTS.md` its invariants.

## The freeze

- **run**: one review of one change, from `socrates review begin` to a terminal state. Identified by a name every later command takes as its argument.
- **round**: one reviewer reading one snapshot, inside a run. A run has at most two.
- **base**: the commit a round is diffed against. Recorded once per run when it begins, and unchanged for the life of the run, so both rounds are diffed against the same thing.
- **fork point**: `merge-base(<ref>, HEAD)`, what a base named with `--base` resolves to, so a diff is the branch's own work rather than the branch's work plus the inverse of everything the named ref did since.
- **gap**: how far a named ref is ahead of the fork point. Zero after a rebase onto that ref, and review says nothing about it.
- **head**: what a round's snapshot tree is built from. Always the working tree, so it holds committed and uncommitted work alike.
- **snapshot**: the commit a round freezes the head into. Built through a temporary index so the person's own index is never touched, content-addressed, and named by no ref.
- **worktree**: the detached checkout of a snapshot that a reviewer is given as its working directory. Never the person's own repository directory.
- **seal**: the write permission removed from every file and directory of a worktree while a round runs. Not a sandbox, and never described as one.
- **guide**: `.socrates/review-guide.md`, the one file read out of the repository under review. Read at the run's base commit, so a guide a change writes governs the next run rather than this one.

## The reviewer

- **reviewer**: the model process a round spawns to read the snapshot. It wrote none of the code and inherits no session from a previous round.
- **reviewer CLI**: the vendor command review drives to spawn a reviewer. One qualifies on stated facts rather than by being listed.
- **seam**: the boundary at `review/agent`, the only package that starts a model process. Nothing above it names a vendor, a flag or an event shape.
- **tier**: whether a vendor constrains output against the schema. A tier and not a criterion, because the document is validated either way.
- **intent**: one line of what the author says they were doing, passed to `socrates review begin`. A claim the reviewer judges against the code, and not a statement of what the change does.
- **prompt**: what is delivered to a reviewer, assembled from named sections. Tests assert on it rather than on review's own source text.
- **hygiene**: model output neutralised on the way into a prompt, per field, and never on the way to disk.

## What comes back

- **findings document**: what a reviewer returns. A document that is missing, dead or unparseable is a failure and never a clean verdict.
- **finding**: one claim about the change, carrying a severity, an action and a location.
- **severity**: how much the finding matters. Anything but informational holds the change.
- **action**: what the finding asks for. `auto-fix` asks for a patch, `ask-user` asks a person, `no-op` asks for nothing.
- **effective action**: the action after the one rule that reads a finding's raw action. The only reader of it.
- **inside the change**: whether a finding's location falls in a hunk of the diff between the run's base and its snapshot. Reported per finding, and the one place the placement of a comment is decided.
- **verdict**: what a run concluded. A run rests in `clean`, `failed`, `escalated`, `awaiting-response` or `cancelled`, and in no other state.
- **parked**: a round that stopped because it waits on a person, which is `awaiting-response`.

## The dialogue

- **response**: what the author records against a finding. Appended, never rewritten.
- **resolution**: which kind of response it is. `fixed` claims a patch, `rejected` claims the finding is false, `decided` records what a person decided, `settled` records what the session that started the run decided where nobody was there to ask, `forwarded` hands the finding to whoever wrote the change, `dropped` records it as true and not raised.
- **settled**: the resolution of a finding that asks a person, on a run with nobody to ask: the session that started the run decided it, and the reason says what it decided. A claim a fresh reviewer judges, as a decision is. Never means a person answered, which design's and agora's glossaries use the same word for.
- **nobody to ask**: what a run records when it begins with no person present to answer a finding. Such a run takes `settled` and never `decided`, and every other run takes `decided` and never `settled`.
- **disposition**: what a resolution of a finding does. A patch in the author's working tree, a comment in somebody else's repository, or nothing at all.
- **reason**: why, required on a resolution that leaves no artifact behind. A rejection, a decision and a drop are claims a fresh reviewer or a later reader judges, never suppressions.
- **handover round**: a round of a change nobody here can patch, so every finding of it is forwarded, dropped or rejected. It waits on a person, and no second round follows it. As a noun, handover is this round and nothing else: handing a round over as a page is a verb every round takes.
- **round cap**: two. Nothing configures it, and a configuration file naming one is refused.

## The settings

- **source**: whether a setting's value is `set` or `default`.

## The suite

- **simulator**: a stand-in for a reviewer CLI, one per supported vendor, reproducing its half of the wire contract.
- **corpus**: what real reviewers answered, kept verbatim, embedded and replayed.
- **fixture**: a recorded shape of a real CLI's answer. Recording one reaches the real CLI, and nothing in continuous integration does.
- **forge check**: the one hand-run check that posts what a review forwarded to a real pull request and deletes it again. It reaches a forge, its reviewer is a simulator because the forge is what it reads, and nothing in continuous integration runs it.
