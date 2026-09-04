## Why

Everything a person reads before they decide to trust socrates, run it, or contribute to it is currently one file, `README.md`, and one file written for an agent, `AGENTS.md`.

The README is good, and what it is missing is at the top rather than in the body.
It opens with what socrates is for and reaches its requirements in the eighteenth line.
The hardest fact about the project is stated there once, in passing: socrates does not work without the `claude` CLI at 2.1.169 or later and the paid account behind it.
That is the first thing a prospective user needs, because for many of them it is the end of the evaluation, and burying it wastes their time to make the project look more available than it is.
The same line states a Go floor of 1.25 that `go.mod` contradicts with 1.26, so a reader who takes the README at its word cannot build.

`agora` is named once, at line 123, in a sentence about which directory it writes to.
`make install` puts it on a person's `PATH`.
A second binary that appears with no explanation is the kind of thing that makes somebody uninstall a tool.

Then there is what is absent.
No security policy, so a person who finds something in a tool that spawns detached processes and rewrites permission bits has a public issue as their only channel.
No contributing guide, no code of conduct, no templates, no ownership.
And no statement anywhere of how a change is proposed in this repository, which matters more here than in most projects.
27 archived changes and 17 promoted capabilities mean a pull request that edits code without a proposal will be turned away by rules the author had no way to read.

socrates already does the hardest part of this well.
The "Out of scope" section names the residual risk on a managed machine rather than hiding it, which is the most trust-building thing in the repository.
What this change adds sits above and around that, and does not touch it.

## What Changes

- `README.md` SHALL state what socrates needs before it works, including the reviewer CLI and the account behind it, before it explains anything else.
- The stated build floor SHALL agree with what the module requires, held by a test rather than by care.
- `README.md` SHALL say what `agora` is, since installing socrates installs it.
- socrates SHALL state that it sends nothing anywhere, because a tool that reads a person's repository and is silent about where that goes is asking to be assumed the worst of.
- The project SHALL state that it is not affiliated with the vendor of the CLI it drives.
- A security policy SHALL state where to report privately, what is in scope, and what socrates does and does not defend against.
- A contributing guide SHALL state the toolchain, the one verification command, the conventions, and above all that a change starts as a proposal in `openspec/`.
- A code of conduct, a support document, ownership and issue and pull-request templates SHALL exist.
- The repository SHALL carry the editor and attribute settings that keep a contributor's tooling from fighting its conventions.

## Capabilities

### New Capabilities

- `project-documents`: what a person reads before they trust the tool, install it, or propose a change to it.

### Modified Capabilities

None.
`cli-foundation` already requires that what a person reads is plain and says what to do next; this change applies that standard to the documents around the tool rather than restating it.

## Impact

- `README.md`: a requirements block and an install pointer at the top, an agora section, a telemetry statement, a non-affiliation note, and the corrected Go floor. The "Out of scope" section is untouched.
- `SECURITY.md`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `SUPPORT.md`: new.
- `.github/CODEOWNERS`, `.github/ISSUE_TEMPLATE/`, `.github/PULL_REQUEST_TEMPLATE.md`: new.
- `.editorconfig`, `.gitattributes`: new.
- `Makefile`: `help`.
- `internal/convention`: a check that the documented Go floor and `go.mod` agree.
- Private vulnerability reporting enabled on the repository.
- Nothing reaches the reviewer, the run, the dialogue, the surface or the page.

## Dependencies

This change SHALL follow `a-version-somebody-can-install` and `what-continuous-integration-actually-proves`.
The install section names a release, and the contributing guide names one verification command; writing either before it exists documents something a reader cannot do.
