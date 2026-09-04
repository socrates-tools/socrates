# Security policy

## Reporting something privately

Use GitHub's private vulnerability reporting: [open a draft advisory](https://github.com/ktei/socrates/security/advisories/new).
It is private to you and the maintainer, and it is the only channel for anything you would not put in a public issue.

Do not open a public issue for a vulnerability.

You should get a first answer within 7 days, and a decision on whether it is in scope within 14 days.
If nothing comes back in 14 days, say so in a public issue without describing what you found, and that will be enough to find it.

## Which versions get a fix

The latest release, and the tip of `main`.
While the major version is `0` there is no long-term support branch, so a fix arrives as a new release rather than as a patch to an old one.

## Where your code goes

socrates sends nothing anywhere.
It opens no network connection outside a build tag the ordinary suite does not compile, and a test holds that.

The reviewer CLI it drives is the other half.
socrates spawns `claude`, and that CLI sends the change to its vendor, Anthropic.
Your diff, the file contents a reviewer reads and the one line of intent you chose to state go there, under the terms your account carries.
A finding about that channel is a finding about the CLI rather than about socrates, and belongs with its vendor.

## What socrates defends against

These are the promises, and each is held by a test that fails when it stops holding rather than by care:

- **A reviewer that writes.**
  The reviewer's tools cannot write, the snapshot worktree it reads is sealed so the filesystem refuses a write from anything in that session, and the round fails if the worktree differs by one byte afterwards.
  A review refuses to start where that seal cannot be proved, including where socrates is running as root.
- **A session carrying your customisations.**
  The reviewer runs under the vendor's own flag for disabling them, loads user-level settings only, and never reads the project's agent instructions, which are authored by the thing under review.
  socrates then reads the session's tool list rather than trusting the request, and records anything it did not ask for.
- **A review inside a review.**
  A marker in the reviewer's environment refuses one, so a repository whose own instructions say to run socrates cannot spawn recursively.
- **A process nobody reaps.**
  Every child is reaped by process group, on clean exit and on error, not only on cancel.
- **A page reachable from another machine.**
  The surface binds a loopback address, an address that is not one is refused before any server starts, and every state-changing route carries a token.
- **A prompt injected from model output.**
  Model output is neutralised per field on the way into a prompt.

## What socrates does not defend against

- **A hook your administrator's policy configures.**
  The vendor's flag for disabling your customisations does not take those with it.
  Two of the three write defences still apply, because the hook's command is a child of the reviewer's session, whose working directory is the sealed worktree: a write it makes there is refused, and the comparison afterwards catches one the seal did not.
  What nothing socrates holds bounds is where that command reaches **outside** the worktree.
  This is stated in `README.md` under "Out of scope" as a residual rather than as a covered case, and it is out of scope here for the same reason.
- **Anything the reviewer CLI or its vendor does with the change.**
  See above.
- **A repository you do not trust.**
  socrates is invoked by the person whose machine it is, on their own repository. It is not a sandbox for reading somebody else's code.
- **Whether a model found the bug.**
  `clean` means an independent reviewer found nothing blocking. It is not a proof of anything, and socrates does not claim it is.
- **Your own machine.**
  socrates changes permission bits across a snapshot worktree and spawns detached processes. That is what it is for.

A report that turns one of the residuals above into something socrates could hold and does not is in scope and worth sending.
