## Context

The core writes a structured log and gives the person one way to look at it that was designed for a machine.
The reason it was designed for a machine is sound and stays: everything `status` prints lands in a language model's context, so it is deliberately tiny.

What is missing is the two commands designed for a person, and the reason they are separate from each other is that following and reading are different problems.
Following has to decide when to stop and what to do before there is anything to read; reading has to decide how much there is and whether to page it.

## Goals / Non-Goals

**Goals:**

- Watching a review is pleasant enough that a person does it instead of guessing.
- A follower that cannot lie: it never sits on a finished run, and it never fails because it was started too early.
- An exit code from `watch` that a hook could already be wired to, so the deferred forcing function needs nothing from this layer later.
- Rendering testable with no terminal, no pseudo-terminal and no golden screenshots.

**Non-Goals:**

- A TUI. No full-screen rendering, no key handling, no alternate screen.
- Multiplexer integration. No tmux, no herdr, no pane adapters.
- Colour theming, progress bars or spinners as a configuration surface.
- Changing anything the core writes.

## Decisions

### Three commands, because there are three consumers

| Command | Behaviour | Consumer |
| --- | --- | --- |
| `scrat watch [run]` | Live, follows to a terminal state, then exits with the verdict | the person, in a second tab |
| `scrat logs [run]` | Dump or page a finished run, no follow | the person, afterwards |
| `scrat status [run]` | One shot, a few lines, no stream | the waiting session |

`status` is the one that is easy to miss and the one that gets called most.
It already exists and this change does not touch it beyond writing down that it is not the other two.

### Rendering reads the log and nothing else

The renderer takes a stream of events and writes lines.
It does not read the run directory, does not know about git, and does not decide anything about the run.

The consequence is the test story: rendering is a pure function from a slice of events to a slice of lines, tested by comparison rather than by driving a terminal.
The terminal awareness is one question asked once, whether the output is a terminal, and it selects a renderer.

### Following is a tail with two stopping conditions

`watch` tails the event log and stops when the run is terminal.
Which of the two it noticed first does not matter, so both are checked: the log carrying a terminal event, and `run.json` recording a terminal state.

Checking only the log is the version that hangs.
A run that failed before its reviewer ever wrote an event has a terminal state and an empty log, and a follower waiting for a terminal event in it waits forever, which is exactly the stall the core's own state rule exists to prevent.

### Racing the spawn is waited out, bounded

`scrat review` returns as soon as the run directory exists, and a person who immediately runs `scrat watch` may arrive before the reviewer has written its first byte.

The follower waits for the log to appear, bounded, and reports the run's state while it waits.
Erroring is the behaviour that punishes the person for being quick, and it happens in exactly the window where they are most likely to try.

The bound is real: a run directory with no log and no live reviewer is a failed run, and the follower says so rather than waiting out its patience.

### Paging is offered, never imposed

`logs` writes to standard output.
Where the output is a terminal and the person has a pager configured, it pages; where the output is not a terminal, it dumps.
A command whose output cannot be piped is a command that has to be worked around.

## Risks / Trade-offs

- **Two stopping conditions is two things to keep true.** The mitigation is that both are cheap and the test for the empty-log case is the one that fails if either is dropped.
- **Rendering choices age.** What reads well now will look dated, and every change to it is a change to something with no assertion beyond a comparison. Accepted: the alternative is a log with rendering baked in, which ages the same way and cannot be re-rendered.
- **The pager is somebody else's program.** A misconfigured one turns `scrat logs` into a hang the person blames on socrates. Falling back to a plain dump where the pager cannot be started is the whole mitigation.
