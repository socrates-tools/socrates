**The suite must be green after every task, not only at the end of a group.**
The ordering below exists for that reason: the check on the toolchain floor is written before the floor it describes is corrected, so that correcting it is what makes the check pass.

**This change follows `a-version-somebody-can-install` and `what-continuous-integration-actually-proves`.**
The install section names a release and the contributing guide names one verification command. A task that documents either before it exists is documenting something a reader cannot do.

**The "what it deliberately is not" section is not edited by this change.**
It names the residual risk on a managed machine and it is the most trust-building prose in the repository. A task that trims, softens or reorders it is in the wrong change.

## 1. What a reader meets first

- [x] 1.1 Add a requirements block to the top of `README.md` naming the reviewer CLI, its version floor, the account it needs, the supported platforms and the toolchain, and verify the em dash check and the user-facing string check still pass
- [x] 1.2 Add a test that the front document states those requirements before it explains how socrates works. Covers `The requirements come before the explanation`
- [x] 1.3 Add a convention check that the toolchain floor a document states and the floor `go.mod` requires agree, and verify it fails against the tree as it stands, naming both. Covers `A stated toolchain floor that disagrees with the module fails`

  The disagreement the proposal names was already gone: `README.md` states Go 1.26 and `go.mod` requires 1.26. The check was verified against a stated floor of 1.25 instead, and it names both.
- [x] 1.4 Correct the stated floor to agree with `go.mod`, and verify 1.3 passes
- [x] 1.5 Add a section to `README.md` saying what `agora` is and what it is for, and add a test that every binary `make install` puts on a path is explained. Covers `The second binary is explained`

## 2. Where what socrates reads goes

- [x] 2.1 State in `README.md` that socrates sends nothing itself, and that the reviewer CLI it drives sends the change to that CLI's vendor
- [x] 2.2 Add a non-affiliation note naming the vendor's trademarks, and verify the em dash check still passes
- [x] 2.3 Add a test that both halves of the reach are stated, and that neither is stated without the other. Covers `The reach is stated in both directions`

## 3. A vulnerability has somewhere to go

- [x] 3.1 Add `SECURITY.md` naming a private reporting channel, the versions supported, and a response window
- [x] 3.2 State in it what socrates defends against and what it does not, pointing at the residual `README.md` already names rather than restating it. Covers `The policy states what is out of scope`
- [ ] 3.3 Enable private vulnerability reporting on the repository, and verify the channel the policy names is the one that is enabled. Covers `A reporter has a private channel`

  Blocked while the repository is private: the endpoint answers 404 because the feature is available for public repositories only. Recorded under "Continuous integration" in `docs/gotchas.md`, and it belongs with making the repository public.

## 4. How a change is proposed here

- [x] 4.1 Add `CONTRIBUTING.md` leading with the fact that a change starts as a proposal in `openspec/`, and verify a reader reaching only the first screen learns it. Covers `The guide leads with the proposal`
- [x] 4.2 State the toolchain, the single verification command, and the conventions the source is held to, naming the entry point rather than listing the gates again. Covers `The guide names one verification command`
- [x] 4.3 State that pinned dependencies produce a steady stream of small changes, so a long queue reads as the mechanism working
- [x] 4.4 Add `make help`, and verify it names every target a contributor is told to run

## 5. The rest of the door

- [x] 5.1 Add `CODE_OF_CONDUCT.md` with a working contact, `SUPPORT.md`, and `.github/CODEOWNERS`

  The contact is GitHub rather than an address: a public issue where a report can be public, and a draft security advisory where it cannot. No inbox is published.
- [x] 5.2 Add issue templates and a pull-request template, and verify the pull-request template asks which proposal the change implements
- [x] 5.3 Add a test that the documents a first-time contributor needs are present. Covers `A first-time contributor finds the documents they need`
- [x] 5.4 Add `.editorconfig` and `.gitattributes` matching the conventions the source is held to, and verify a checkout on another platform produces no whitespace or line-ending diff. Covers `A contributor's editor does not fight the conventions`
