## Why

There is one way to get socrates, and it is to clone the repository and run `make build`.
There are no tags, no archives, no checksums and no changelog.

The version story is the sharper end of this.
`internal/version` stamps the commit and whether the tree was dirty, which answers "which build is this" precisely and answers "which release is this" not at all.
A person reporting a bug can name a commit they may not have, or name nothing.
A person deciding whether to upgrade has nothing to compare.
A person who wants to pin socrates in a script has no version to pin to.

The verification story is why this is more than convenience.
socrates spawns detached processes, changes permission bits across a worktree, and reads a repository the person cares about.
A tool that does that, distributed as an unsigned binary from an account somebody has never heard of, is asking for trust it has given nobody the means to check.
Checksums, signatures and build provenance are how a stranger verifies that the binary they downloaded is the one this repository's tagged commit produced, without having to trust the publisher at all.

The build stamp is not replaced by any of this.
`socrates version` says which commit it was linked from and whether that tree was clean, and that stays true and stays useful.
What is added is a release the stamp belongs to.

## What Changes

- A release SHALL be named by a semantic version, and `socrates version` SHALL report it alongside the commit stamp it already reports.
- The version SHALL be stamped at link time from the tag, in the same way the commit already is, and SHALL NOT be a constant in the source that somebody has to remember to bump.
- A build from an untagged tree SHALL say so rather than claim a version it is not.
- `CHANGELOG.md` SHALL record what changed in each release, written by hand rather than generated, because this repository's commit messages say what a change is and a changelog says what it means to somebody using it.
- Tagging SHALL produce, for every supported platform and no others, an archive, a checksum for it, a signature over the checksums, a bill of materials, and build provenance.
- The archive SHALL contain both names of the tool and the second binary, so that what a person downloads is what `make install` would have given them.
- `README.md` SHALL say how to install a release and how to install from source, with the release first.
- Nothing in the release path SHALL reach the real CLI, which the existing build-tag boundary already holds and which this change must not breach.

## Capabilities

### New Capabilities

- `distribution`: what a release is, what it contains, how a stranger verifies it, and what a build says about which release it is.

### Modified Capabilities

None.
`cli-foundation` already requires that `version` says which build this is; the version reported is added to what that build stamp carries rather than replacing it, and the requirement's shape is unchanged.

## Impact

- `internal/version`: a version alongside the commit and the dirty flag, and what it says when there is no tag.
- `Makefile`: the version in the ldflags beside `STAMP` and `SURFACE_STAMP`.
- `CHANGELOG.md`: new.
- `.github/workflows/release.yml`: new, on tag, gated on the verification matrix.
- `.goreleaser.yaml`: new, building the platforms `platform-support` states and no others.
- `internal/e2e`: scenarios for what a tagged and an untagged build report.
- `README.md`: an install section, above the from-source instructions that exist today.
- `NOTICE` and `LICENSE` travel in every archive, which is what `the-licence-the-work-is-published-under` produced them for.
- Nothing reaches the reviewer, the run, the dialogue, the surface or the page.

## Dependencies

This change SHALL follow `the-platforms-the-seal-holds-on`.
The release matrix is the supported set, and a release built before that set is stated would ship a socrates whose reaping guarantee does not hold on a platform nobody decided to support.

This change SHALL follow `what-continuous-integration-actually-proves`.
A release publishes nothing unless verification passed, so the verification it gates on has to exist and has to run on every supported platform before a tag can mean anything.
