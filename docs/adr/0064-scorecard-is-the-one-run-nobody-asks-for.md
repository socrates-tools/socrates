# 64. Scorecard is the one run nobody asks for

- Status: accepted
- Date: 2026-10-14
- Part of: scorecard

## Decision

`.github/workflows/scorecard.yml` runs OpenSSF Scorecard once a week on a schedule, and whenever somebody asks for it.
It publishes its results to the OpenSSF API, which is what the README's badge reads, and uploads them to the repository's code scanning.

It is the one workflow that starts without anybody asking.
Every other run happens when a person asks for one or a release gates on one, and nothing starts on a push, because a run nobody asked for is one this account pays for without anybody wanting it.
Scorecard is cheap where those are dear: it runs the suite on no platform, starts no browser, builds nothing, and is one job on one Linux runner for about a minute.

What it buys is a regression nobody would otherwise see.
The hardening Scorecard measures, pinned actions, read-only tokens, a vulnerability check and reviewed dependency updates, is mostly held by convention checks already, but the settings it also reads, such as branch protection, live on the forge rather than in the tree, and nothing in the suite can see them change.

Scorecard publishes only from a scheduled run or a push to the default branch, and treats a run started on request as experimental, so the schedule is what makes the badge true rather than an accident of when somebody last asked.

## Rejected

- Running it on request and before every release, like code scanning.
  It keeps the rule without an exception, and the badge would then say whatever the last release said, or nothing, since publishing from a run on request is experimental.
- Running it on every push to `main`, as Scorecard's own example does.
  It is the trigger the rule exists to refuse, and a weekly reading is enough for something that changes when a setting does.
- Running it without publishing, so the results reach code scanning and nothing else.
  It answers the question for the maintainer and not for a person deciding whether to trust the release, which is most of what the badge is for.
