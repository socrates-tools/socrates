# 26. The menu is bubbletea, and what that costs

- Status: accepted
- Date: 2026-09-13
- Part of: socrates-review configure

## Decision

`socrates-review configure` draws its menu with `charm.land/bubbletea/v2`.

The two binaries linked five modules before this change.
bubbletea brings sixteen more, and every one of them is MIT or BSD-3, so `NOTICE` grows by sixteen entries each carrying its licence in full.
`github.com/charmbracelet/ultraviolet` and `github.com/xo/terminfo` publish no tag and are pinned to a commit, which their entries state.

That is the price, and it is paid once.
What it buys is a terminal program that reads keys, redraws a frame and restores the terminal on the way out, on both supported platforms, rather than a hand-rolled raw-mode loop that this repository would then own.
A raw-mode loop is not hard to start and is hard to finish: bracketed paste, a resize mid-frame, a signal, a terminal that answers a capability query, and the restoration of the terminal when the process is killed are each a bug a person meets rather than a test.

Nothing else comes with it.
No `bubbles` and no `lipgloss`: `internal/configure` builds its view out of plain strings, so the two widget libraries and the styling library stay out of the module graph and out of `NOTICE`.

`internal/convention/wording.go` finds user-facing strings by reading `fmt`, `errors` and `refusal` call sites, and a view that returns a string reaches none of them.
The check is therefore widened to read every prose literal in `internal/configure`, quoted or raw, so the menu's labels are held to the same registry rule as every other package rather than escaping it because of how a view is built.
The one thing it lets through is a literal passed to `regexp`, because a pattern is matched on rather than read.

## Rejected

- a hand-rolled raw-mode loop, which keeps the module graph at five and moves the cost from `NOTICE`, which is read once, to the source, which is read forever.
- `bubbles` for the list and the text field, which is the obvious convenience and adds a dependency for two screens whose whole logic is a cursor and a string.
- `lipgloss` for the view, which would make the menu's output a styled string and the suite's assertions about it a test of a renderer.
- exempting `internal/configure` from the message registry because its strings reach nobody through `fmt`, which is true of how they are written and false of what they are.
