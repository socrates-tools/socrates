## Context

The core layer shipped with two mechanisms for keeping a reviewer read-only, and named the second as the load-bearing one:

1. The vendor's own tool restriction, so a write is refused at the point of the attempt.
2. A comparison of the snapshot worktree against the snapshot commit when the round ends.

The first live run proved that the first mechanism does not cover what actually writes.
`--tools Read,Grep,Glob` restricts the tools the model may call.
It says nothing about the tool servers the CLI starts, which are separate processes that write when they feel like it, and one of them writes a project directory into its working directory at startup.
The reviewer's working directory is the snapshot worktree, so the second mechanism fired, correctly, on every run.

Asking the CLI to load no tool server closes the hole where the CLI allows it.
On a machine with an enterprise tool server configuration the CLI refuses the request outright, and that configuration is at a path the person does not own.
socrates then has one mechanism left, and it is a report rather than a defence.

## Goals / Non-Goals

**Goals:**

- A reviewer that cannot write in the snapshot worktree, whatever the vendor's CLI loads and whatever it refuses to stop loading.
- A property socrates proves at the start of every run rather than assumes from a `chmod` call that returned no error.
- The comparison kept exactly as it is, as the check on anything that defeats the seal.
- socrates usable on the machine it is written on, without editing that machine's enterprise configuration.

**Non-Goals:**

- Naming paths socrates agrees not to notice. An exception list for a third party's bookkeeping directory is a rule that depends on that third party, which is the thing the core layer refused to do.
- Running the reviewer somewhere other than the snapshot worktree. A scratch working directory with the worktree added as a readable path would also stop the write, and it costs the reviewer the thing every model does by reflex, which is to open a relative path in the directory it was started in.
- Containers, user namespaces, or a read-only bind mount. They achieve the same property with a dependency socrates does not otherwise have, and the permission bits are already there.
- Windows support beyond the probe deciding whether socrates runs at all.

## Decisions

### The worktree is sealed after checkout and unsealed when the round ends

After `git worktree add --detach`, socrates walks the worktree and clears every write bit: directories become `0o555` and files `0o444`.
A new file cannot be created in a directory that is not writable, and an existing file cannot be modified, so both shapes of write are refused by the kernel rather than by anybody's good behaviour.

The seal covers the worktree only.
The run directory around it stays writable, because the prompt, the stream and the findings are written there while the round runs.
The worktree's `.git` file is sealed with everything else; git writes the worktree's index and administrative state under the main repository's `.git/worktrees/<name>/`, which is untouched, and `git status --porcelain` in a sealed worktree was verified to work.

The seal is lifted when the round ends, on every path:

- The worktree is removed on a terminal verdict, and `os.RemoveAll` cannot unlink from a directory it may not write.
- The worktree is kept when the comparison failed, and a person inspecting it should not have to run `chmod -R` before they can read a diff or delete the directory.
- `Checkout` already removes a partially created worktree on failure, and that path needs the same treatment.

### The seal is proved by one write, before any reviewer exists

`os.Chmod` returning nil is not evidence.
A filesystem can ignore permission bits, and a `root` process ignores them by design.

So socrates creates one probe file in the sealed worktree root.
If the write is refused, the seal holds, and the review proceeds.
If the write succeeds, socrates removes the probe and refuses the review, naming the worktree's path and that its filesystem does not enforce the seal.

The refusal happens in `scrat review`, before a reviewer is spawned, so the person gets a usage error at their terminal rather than a failed run to investigate.
A refusal is the honest answer: on such a filesystem socrates cannot promise that the change under review is the change that was reviewed, and every verdict it produced would be a claim it could not support.

This also answers the platform question without a platform check.
Wherever the bits work, socrates runs.
Wherever they do not, socrates says so in one sentence and does nothing.

### The comparison stays, and the simulator learns to defeat the seal

The comparison after the round is what catches a write the seal could not stop: a process running as another user, a filesystem that started enforcing nothing after the probe, a future mechanism nobody has thought of.
Keeping it costs one `git status` and removes a whole class of quiet failure.

It also becomes unreachable by an ordinary reviewer, which is a problem for the tests that hold it.
Today the simulator writes a file into its working directory on demand, and that is how three scenarios are driven.
After this change that write is refused, so the simulator gains a knob that clears the way first: it makes the path writable, writes, and leaves the worktree dirty exactly as a more privileged writer would.

That is the simulator doing its job.
It exists to produce on demand what the real thing produces rarely, and "something wrote here despite the permissions" is precisely such a case.

### What the reviewer is told

Nothing changes in the prompt.
The reviewer is not told the worktree is read-only, because it has no reason to try to write and no useful response to being told it cannot.
The tools it is given do not write, and a tool server that fails to write says so in its own log rather than in the review.

## Risks / Trade-offs

- **A tool server that cannot write may fail to start.** Verified against the real CLI on the affected machine: the session continued and answered normally. The reviewer loses whatever tools that server would have offered, which is the intended outcome, because a server socrates cannot restrict is not a tool socrates wants in a reviewer.
- **The seal is a third mechanism, and mechanisms accumulate.** Accepted, and bounded: this one is fifteen lines of `filepath.WalkDir` plus a probe, it depends on nothing outside the standard library, and it replaces no existing check. The alternative on offer was an exception list keyed to a third party's directory name.
- **A person who reads a kept worktree while the round is still running sees a read-only directory.** That is a curiosity rather than a cost, and the state is legible with `ls -l`.
- **`chmod` on a large worktree costs a walk.** It is one pass over files git has just written, so the cost is bounded by what the checkout already paid.
- **A filesystem that ignores permission bits turns socrates off rather than degrading it.** That is the point, and the refusal names what to do about it.
