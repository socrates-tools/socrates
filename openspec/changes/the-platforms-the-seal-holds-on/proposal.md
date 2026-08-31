## Why

socrates makes two promises that are not matters of care: the reviewer never writes, and everything the reviewer started is reaped with it.
Both are held by unix mechanisms.
The seal clears write bits and proves itself with a write before a reviewer is spawned; the reaping sets a process group on spawn and kills the group afterwards.

Off unix, one of them is not implemented.
`internal/process/process_other.go` returns "socrates reaps by process group, which this platform does not have", and the seal's permission bits mean something different on a platform whose access control is not permission bits.
So a socrates built for such a platform runs, reviews, and returns a verdict, with one of its two load-bearing guarantees absent and nothing saying so.

The README names no platform at all.
That is defensible while the only user is the author, who knows what they built it on.
It stops being defensible the moment a stranger can download it, because the stranger's reasonable assumption is that a Go program runs everywhere Go runs.

socrates already knows how to handle a machine where the seal cannot hold: it refuses.
`ProveTheSeal` writes into the sealed worktree before any reviewer starts and refuses the review if the write succeeds, which covers a filesystem that ignores permission bits and covers running as root.
The reasoning given for that refusal is that a verdict socrates cannot stand behind is worse than no verdict.
That reasoning does not stop applying because the reason the guarantee is absent is the platform rather than the filesystem.

## What Changes

- The platforms socrates supports SHALL be stated: which are supported, which are refused, and on what basis.
- socrates SHALL refuse to review on a platform where reaping by process group is not available, before any reviewer is spawned, with a refusal that says which guarantee is missing and what to do instead.
- The refusal SHALL name the same remedy the documentation names, so a person who meets it is not left to search.
- `README.md` SHALL state the supported platforms and point Windows users at WSL.
- A platform SHALL NOT be addable to what a release builds without an entry in the support statement, held by a test rather than by care.

## Capabilities

### New Capabilities

- `platform-support`: which platforms the guarantees hold on, what socrates does where they do not, and how a release and the statement are kept in agreement.

### Modified Capabilities

None.
`reviewer-run` already requires that the reviewer never writes and that everything it starts is reaped together. This change does not weaken or restate either; it states where they hold and refuses where they do not.

## Impact

- `internal/process`: the unsupported path becomes something a caller must handle rather than an error a caller may ignore.
- `internal/reviewer`: the refusal, raised where `ProveTheSeal` is already raised, before a reviewer is spawned.
- `internal/refusal`, and the relevant `messages.go`: the refusal's text.
- `internal/platform` or an equivalent single place: the supported set, read by the refusal, by the test and by the release matrix rather than written three times.
- `README.md`: a supported-platforms section.
- `internal/e2e`: a scenario that the refusal happens before anything is spawned.
- Nothing reaches the findings, the dialogue, the corpus, the surface or the page.
