# Context

`socrates review` puts an independent reviewer in front of unreviewed local changes.
`socrates agora` is a second namespace, inside the same binary, that serves one document as one page, waits for a person, and prints what they answered.
`socrates design` is a second workflow, inside the same binary, that interviews a person over agora and leaves a spec and its decisions behind.
It is a skill rather than a tool.

This file is the vocabulary.
Use these words rather than synonyms, in a spec, an issue title, a test name or a commit message.
`docs/adr/` holds the decisions and `CLAUDE.md` the invariants.

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
- **resolution**: which kind of response it is. `fixed` claims a patch, `rejected` claims the finding is false, `decided` records what a person decided, `forwarded` hands the finding to whoever wrote the change, `dropped` records it as true and not raised.
- **disposition**: what a resolution of a finding does. A patch in the author's working tree, a comment in somebody else's repository, or nothing at all.
- **reason**: why, required on a resolution that leaves no artifact behind. A rejection, a decision and a drop are claims a fresh reviewer or a later reader judges, never suppressions.
- **handover round**: a round of a change nobody here can patch, so every finding of it is forwarded, dropped or rejected. It waits on a person, and no second round follows it. As a noun, handover is this round and nothing else: handing a round over as a page is a verb every round takes.
- **round cap**: two. Nothing configures it, and a configuration file naming one is refused.

## The settings

- **setting**: one thing a person can configure, named by its YAML key, carrying the value in force, whether that value is theirs or a default, and the default it falls back to. What `socrates review config` reports and what `socrates review configure` offers are the same set.
- **source**: whether a setting's value is `set` or `default`. Not agora's source, which is a different word in a different vocabulary.
- **use the default**: removing a setting's line, rather than writing today's default into the file, so that the file holds only what a person decided.

## The design

- **change**: one unit of work, from the spec that opens it to the diff a reviewer reads. Named by a slug, and `docs/changes/<slug>/` holds its spec.
- **spec**: what a design leaves behind, holding the problem, the solution, the user stories and the decisions. Published as an issue, and committed to the change's own folder.
- **design**: one interview about one change, from its first round to the spec and the decision records it writes. The interview is the session itself, so nothing runs between one round and the next.
- **design round**: one page of questions put to a person, inside a design. Not review's round, which is a different word in a different vocabulary.
- **frontier**: every question whose prerequisites are settled. What one design round asks, and a design ends when it is empty.
- **question**: one thing a design needs decided, carrying the answers it offers and the one the session would take. Identified for the life of the design, so it keeps its name across rounds.
- **settled**: a question the person has answered. It stays on every later page as a card that states the answer and offers nothing.

## The page

These are agora's words, and agora knows nothing about code review.

- **document**: what a caller writes for the surface to render. Not a page.
- **source**: the entries a caller cites, each carrying a name, a severity, an audience and a location.
- **decision**: one card on the page, citing the source entry it came from.
- **option**: one answer to a decision. It carries the command it becomes, or none, and a decision does not mix the two: a person choosing between the options of one decision is choosing between things of the same kind. The surface composes a command and never runs one. The control that takes none of them is not an option: it carries no command, and leaves the decision with nothing selected.
- **illustration**: the one fragment a decision may carry, painted by the document's own stylesheet, contained so nothing in it reaches the page's own controls.
- **session**: what the surface mints so that a page and its routes are served under a path holding it. Kept for the whole review.
- **exchange**: the delivery of a document and the collection of what a person answered. It holds at most one uncollected answer, and a later one replaces it.
- **answer**: what one submit carries: every decision the person answered, anything they wrote, and whether it ends the exchange.
- **confirmation**: the dialog a press raises where any decision is blank, naming each one and taking the person to it. The only thing on the page standing between a press and a send.

## The suite

- **simulator**: a stand-in for a reviewer CLI, one per supported vendor, reproducing its half of the wire contract.
- **scenario**: the behaviour a test is named for. A task is not done until a test named for the scenario failed for the right reason and then passed.
- **corpus**: what real reviewers answered, kept verbatim, embedded and replayed.
- **fixture**: a recorded shape of a real CLI's answer. Recording one reaches the real CLI, and nothing in continuous integration does.
- **forge check**: the one hand-run check that posts what a review forwarded to a real pull request and deletes it again. It reaches a forge, its reviewer is a simulator because the forge is what it reads, and nothing in continuous integration runs it.
