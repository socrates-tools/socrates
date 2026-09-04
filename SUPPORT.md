# Support

## Read these first

- **`README.md`** answers what socrates-review needs, what it does to your repository, what the exit codes mean, and what it deliberately is not.
  Most questions are answered by "What it deliberately is not".
- **`docs/gotchas.md`** holds the constraints somebody would otherwise re-derive: the ones that cost time once and should not cost it twice.
- **`openspec/specs`** says what is required and why, and the design document of each archived change under `openspec/changes/archive/` says what was decided and what was rejected.
  If you are asking "why does it work like that", the answer is in one of those.
- **`CHANGELOG.md`** says what each release means to somebody who has the previous one.

## Then

| what you have | where it goes |
| --- | --- |
| socrates-review did something you did not expect | an issue, with the output of `socrates-review logs <run>` and what `socrates-review version` prints |
| socrates-review refused and you do not know why | an issue. Every refusal names what to do next, and one that does not is a bug in itself |
| you want socrates-review to do something it does not | an issue saying what you are trying to achieve, not a pull request. See `CONTRIBUTING.md`: a change starts as a proposal |
| a vulnerability | `SECURITY.md`. Not an issue |
| a question about conduct | `CODE_OF_CONDUCT.md` |
| the `claude` CLI itself misbehaved | its vendor. socrates-review drives it and does not maintain it |

## What to put in an issue

- What `socrates-review version` prints.
- Your platform.
- What you ran, and what you expected instead of what happened.
- The run's log where there is a run: `socrates-review logs <run>`.
  Read it before you paste it. It quotes the change that was under review.

## What this is not

There is no support contract, no response commitment for anything but a security report, and one maintainer.
An issue that says what happened, on what, with the log attached, gets an answer faster than any amount of asking whether anybody is there.
