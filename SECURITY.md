# Security policy

## Reporting something privately

Use GitHub's private vulnerability reporting: [open a draft advisory](https://github.com/socrates-tools/socrates/security/advisories/new).
It is private to you and the maintainer, and it is the only channel for anything you would not put in a public issue.

Do not open a public issue for a vulnerability.

You should get a first answer within 7 days, and a decision on whether it is in scope within 14 days.
If nothing comes back in 14 days, say so in a public issue without describing what you found, and that will be enough to find it.

## Which versions get a fix

The latest release, and the tip of `main`.
While the major version is `0` there is no long-term support branch, so a fix arrives as a new release rather than as a patch to an old one.

## Where your code goes

socrates review sends nothing anywhere.
The code a review runs opens no network connection outside a build tag the ordinary suite does not compile, and a test holds that by reading what the review command tree and every package it reaches import, and what the dispatcher and the binary's own main import directly.

The binary that runs reviews now contains an HTTP server, where it once provably did not: `socrates` is one binary, and agora's page server is linked into it beside the reviewer.
Running a review starts no server and reaches no network.
Running `socrates agora serve` binds a loopback address on purpose, to put a decision in front of a person, and a separate test holds that the address it binds cannot be widened.

The reviewer CLI it drives is the other half.
socrates review spawns the reviewer CLI you configured, and that CLI sends the change to its vendor: `claude` to Anthropic, `codex` to OpenAI, and `opencode` to whichever model provider you configured it with.
Your diff, the file contents a reviewer reads and the one line of intent you chose to state go there, under the terms your account carries.
A finding about that channel is a finding about the CLI rather than about socrates review, and belongs with its vendor.

agora's page also draws an illustration, an HTML file a caller wrote, in a frame sandboxed without `allow-same-origin`.
A drawing may load scripts, styles and fonts from three hosts: `cdnjs.cloudflare.com`, `cdn.jsdelivr.net` and `unpkg.com`.
It may load from nowhere else, and a drawing sends nothing: the policy it is served under names nothing for connecting, framing or submitting a form, so what it holds goes to none of those hosts or anywhere else.
Opening a page that carries an illustration may fetch a library from one of those three hosts, and that host sees the request for it.

## What socrates review defends against

These are the promises, and each is held by a test that fails when it stops holding rather than by care:

- **A reviewer that writes.**
  The reviewer's tools cannot write, the snapshot worktree it reads is sealed so the filesystem refuses a write from anything in that session, and the round fails if the worktree differs by one byte afterwards.
  A review refuses to start where that seal cannot be proved, including where socrates review is running as root.
- **A session carrying your customisations.**
  The reviewer runs under the vendor's own flag for disabling them, loads user-level settings only, and never reads the project's agent instructions, which are authored by the thing under review.
  socrates review then reads the session's tool list rather than trusting the request, and records anything it did not ask for.
- **A review inside a review.**
  A marker in the reviewer's environment refuses one, so a repository whose own instructions say to run socrates review cannot spawn recursively.
- **A process nobody reaps.**
  Every child is reaped by process group, on clean exit and on error, not only on cancel.
- **A page reachable from another machine.**
  The surface binds a loopback address, an address that is not one is refused before any server starts, and every state-changing route carries a token.
- **A prompt injected from model output.**
  Model output is neutralised per field on the way into a prompt.

## What socrates review does not defend against

- **A hook your administrator's policy configures.**
  The vendor's flag for disabling your customisations does not take those with it.
  Two of the three write defences still apply, because the hook's command is a child of the reviewer's session, whose working directory is the sealed worktree: a write it makes there is refused, and the comparison afterwards catches one the seal did not.
  What nothing socrates review holds bounds is where that command reaches **outside** the worktree.
  This is stated in `README.md` under "Out of scope" as a residual rather than as a covered case, and it is out of scope here for the same reason.
- **Anything the reviewer CLI or its vendor does with the change.**
  See above.
- **A repository you do not trust.**
  socrates review is invoked by the person whose machine it is, on their own repository. It is not a sandbox for reading somebody else's code.
- **Whether a model found the bug.**
  `clean` means an independent reviewer found nothing blocking. It is not a proof of anything, and socrates review does not claim it is.
- **Your own machine.**
  socrates review changes permission bits across a snapshot worktree and spawns detached processes. That is what it is for.

A report that turns one of the residuals above into something socrates review could hold and does not is in scope and worth sending.
