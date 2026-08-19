## Context

The tree being reviewed is alive.
The person's session may keep editing while the reviewer reads, and without a freeze the two race: the reviewer reports findings anchored to lines in a file that has since changed, and every one of them is wrong in a way that is expensive to discover.

The obvious freeze is wrong.
`git stash create` silently ignores untracked files, and a freshly implemented feature is mostly untracked files.
Run against a real repository whose only changes were untracked, it returned nothing at all, and the review that followed was a confident review of an empty diff.

Everything else in this change follows from two facts: the reviewer is a separate process that must not write, and a review that did not happen must be distinguishable from a review that found nothing.

## Goals / Non-Goals

**Goals:**

- One round, end to end, on a real repository: freeze, review, findings, verdict.
- The person's working tree and index are provably untouched.
- A verdict that lies is impossible: `clean` requires positive evidence that a review completed.
- An end-to-end suite that drives the real binary and substitutes only the agent and the clock.
- The pure decision cores are small enough to test exhaustively.

**Non-Goals:**

- Rounds, rebuttals and the dialogue. That is `a-rejection-is-a-claim`, and this change deliberately reaches only the point where findings exist and belong to the person.
- `scrat watch` and `scrat logs`. That is `one-window-three-needs`. This change ships `scrat status`, because the waiting session needs it and it is three lines.
- The generated skill. That is `the-loop-an-agent-can-drive`.
- Running any test or lint command. Deferred by the scope rule in `openspec/config.yaml`.
- A second agent implementation, a hook, a TUI, a daemon, a database.

## Decisions

### The snapshot is built through a temporary index

```sh
export GIT_INDEX_FILE=$(mktemp -u)
git read-tree HEAD
git add -A
tree=$(git write-tree)
snap=$(git commit-tree "$tree" -p HEAD -m "socrates snapshot")
git worktree add --detach "$wt" "$snap"
```

`GIT_INDEX_FILE` pointing at a path that does not exist yet is what keeps the person's own index out of it: `read-tree` creates the temporary one from `HEAD`, and `add -A` stages into that.
`add -A` honours `.gitignore` and includes untracked files, which is the entire reason this recipe exists.
Verified against a live repository: 1708 insertions captured including untracked files, and `git status` afterwards byte-identical.

The snapshot commit is unreferenced, so ordinary garbage collection reclaims it once the worktree is removed.
socrates does not create a ref for it, because a ref is a thing to clean up and a thing to leave behind.

**A repository with no `HEAD` is refused rather than special-cased.** A repository with no commit has no base to diff against, and the review socrates would produce is a review of everything, which is not what anybody meant. The refusal names the condition.

### Run identity is a directory name, not a lookup

A run directory is `<state>/runs/<UTC timestamp>-<repo key>`, where the repo key is the first eight hex characters of the SHA-256 of the repository root's absolute path.

The consequence is the point: "the most recent run for this repository" is a `readdir` and a suffix match, sorted by name, reading no files.
The alternative, a flat run identifier plus a repository recorded inside each `run.json`, answers the same question by opening every run.
The command that asks it most is `scrat status` called in a loop by a waiting session.

### Run state is files, and every write is atomic

`run.json` is written to a temporary file in the same directory and renamed over the target.
A reader that opens a run mid-write gets the previous complete state or the new complete state, never a truncated one.
This matters more here than the file count suggests: the writer is a detached process and the reader is a session polling it, and they have no other synchronisation.

The event log is JSONL rather than rendered text, so that `status` summarises the same bytes `watch` will later render and a consumer after that parses rather than scrapes.
It is bounded, and rotation truncates the current inode in place, because the reviewer holds a descriptor on it and a rotation that renames leaves the reviewer writing to a file nobody reads.

### Three terminal verdicts, and the third is not a failure of the code

- `clean`: a review completed and reported nothing blocking.
- `escalated`: a review completed and what it found is the person's.
- `failed`: no review completed.

`failed` is separate from `escalated` because they call for different actions.
Escalated means read the findings; failed means socrates could not do its job and nothing has been checked.
Collapsing them is how "empty is not green" gets reintroduced by somebody simplifying the exit codes.

In this change every blocking finding escalates, because there is nobody to hand an auto-fixable finding to yet.
`a-rejection-is-a-claim` inserts the dialogue before the escalation without changing what the three verdicts mean, which is why the exit-code contract is stated now rather than later.

### The reviewer's read-only property is enforced structurally

Two mechanisms, and only the second is load-bearing.

1. The reviewer is invoked with the vendor's own tool restrictions, so writing is refused at the point of the attempt.
2. When the round ends, socrates compares the snapshot worktree against the snapshot commit. A difference fails the round, names the paths, and keeps the worktree for inspection.

The second is what the specification requires, because it is true whatever the vendor's flags are called next year and whatever a future implementation behind the seam supports.
A rule that depends on a third party's flag name is a rule that quietly stops holding.

### The reviewer loads user settings only

The reviewer is a coding session pointed at a repository whose own agent instructions are, from socrates' point of view, untrusted input authored by the thing under review.
Project instructions can seat an unrelated identity on the reviewer, and in the worst case they tell it to run socrates, which spawns recursively.

So the reviewer loads user-level settings and memory only, and carries an environment marker naming the run it belongs to.
A `scrat review` invoked in a process whose environment already carries that marker refuses, naming the outer run.
Belt and braces on purpose: the marker is one line and the failure without it is unbounded process creation.

The cost is real and accepted: the reviewer does not know the project's conventions.
The "stay in lane" rule in the prompt covers most of what that would otherwise produce, and a socrates-owned channel for repository review guidance is on the roadmap rather than in this change.

### The diff is inlined below a ceiling and read above it

The reviewer's working directory is always the snapshot worktree, so it can inspect a call site the diff does not show, which is exactly where the counterexample-tracing rule earns its keep.

The prompt carries the diff itself when the diff is under a byte ceiling, so the common small change needs no exploration turns at all.
Above the ceiling the prompt carries the changed-file list and the base commit, and says to read what it needs.
Measured `no-mistakes` review turns that read the repository themselves ran ten to sixty minutes, and the ceiling is what keeps the ordinary case away from that.

The ceiling is one setting with a default. It is not per repository, and it is not inferred from a token count, because a byte count is something a person can predict.

### One seam, and the simulator is why it is a seam

```go
type Role string

type Request struct {
	Role   Role
	Prompt string
	Dir    string
	Schema json.RawMessage
}

type Response struct {
	Structured json.RawMessage
	Text       string
	Model      string
	Usage      Usage
}

type Agent interface {
	Run(context.Context, Request) (Response, error)
}
```

The Claude implementation invokes `claude -p --verbose --output-format stream-json --setting-sources user --json-schema <schema>`, parses the JSONL event stream, and takes the structured output from the terminal `result` event.
A schema was supplied and no structured output came back is an error, not an empty answer.

The seam is narrow enough that the simulator implements the same wire contract rather than the Go interface: `cmd/claudesim` is a binary symlinked as `claude` into a temporary `PATH` directory, so the end-to-end suite exercises argument construction, stream parsing, process handling and reaping, all of which are where the bugs are.
Faking at the Go interface would test none of it.

### Testing, and what the simulator is for

The end-to-end suite is primary.
It builds the real `scrat`, creates a real temporary repository, runs real `git`, spawns real processes, and substitutes the agent and the clock.

The simulator earns its place by producing on demand what the real tool produces rarely:

- a stream that ends mid-event
- a `result` event with no structured output when a schema was supplied
- structured output that parses but violates the schema
- a non-zero exit after a complete stream
- a process that writes nothing and never returns
- a process that spawns a child which outlives it
- a reviewer that writes to the worktree

Every one of those is a specification scenario in this change, and none of them is reachable against the real tool on demand.

Its fidelity is a separate question with a separate answer.
Responses from the real `claude` CLI are recorded behind a build tag, and a drift check compares the simulator's output shape against the recordings.
A non-empty diff means the CLI changed.
A third suite, also behind its own tag, runs one real review against the real binary and asks whether the flags are still accepted and the schema still honoured, which no recording can answer.

## Risks / Trade-offs

- **The simulator can drift into fiction.** A simulator nobody compares to the real thing becomes a specification of its author's beliefs, and the suite then manufactures confidence. The recording and drift check are the whole mitigation, and the scheduled job that runs them is what makes them something other than a command nobody types.
- **`git add -A` in a repository with a large ignored tree is slow.** It is bounded by what `git status` already costs on the same repository, which the person pays constantly, so it is a cost they can predict. If it becomes a problem the answer is a pathspec, not a different snapshot mechanism.
- **`--setting-sources user` loses project conventions.** Accepted here and revisited by the roadmap item for repository review guidance, which is a socrates-owned channel rather than a re-enabling of arbitrary project instructions.
- **The detached reviewer outlives a `scrat review` that is killed.** That is deliberate: it is the property that lets the session go back to work. Process-group reaping covers the reviewer's own children; a reviewer orphaned by a machine that lost power is cleaned up by the run's own staleness rule, not by a supervisor socrates does not have.
- **A byte ceiling on the inlined diff is a blunt instrument.** A 200 KiB diff of generated files and a 200 KiB diff of dense logic are not the same review, and one number treats them alike. It is still better than a token estimate nobody can predict, and the reviewer above the ceiling is not blinded, only slower.
