---
name: release
description: Cut a release of socrates. Use when asked to release, tag, publish or ship a version, to write the changelog entry for one, or to find out what happened to a release already under way.
---

# Cut a release

A release of socrates is a tag and nothing else.
Pushing a tag matching `v*` starts `.github/workflows/release.yml`, which runs the whole verification matrix against the tagged commit and publishes only if it passes.
Nothing is published by hand, and no step of this is done in the GitHub web interface.

Everything before the tag is yours: deciding the version, writing what the release means to somebody who has the previous one, and proving locally that the archives build.

## Before you tag

1. Be on `main`, with a clean tree, level with `origin/main`.

   ```sh
   git status --porcelain && git fetch origin && git log --oneline origin/main..HEAD
   ```

   All three should print nothing.
   A tag on a commit that is not on `origin/main` publishes a release nobody can check out.

2. Decide the version.
   It is semantic, and the tag carries a `v` prefix while `CHANGELOG.md` does not.
   The major version is `0`, so a minor release may change the surface half in ways a script has to follow, and a patch release may not.
   Read what is under `## [Unreleased]` and pick from that rather than from the commit log.

   ```sh
   gh release view --repo socrates-tools/socrates-review --json tagName --jq .tagName
   ```

3. Write the changelog entry.
   `CHANGELOG.md` is written by hand and is the one document a person reads to decide whether to upgrade, so write it in terms of what the release means to somebody who has the previous one, not in terms of what the commits did.
   The commit messages already say what each change is.

   - Rename `## [Unreleased]` content into a new `## [X.Y.Z] - YYYY-MM-DD` section, dated the day you publish, and leave `## [Unreleased]` above it holding nothing.
   - Group under `### Added`, `### Changed`, `### Fixed`, `### Removed` and `### Notes`, following the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format the file already states.
   - One sentence per line, and no em dashes.
   - Update the link definitions at the bottom: point `[Unreleased]` at `compare/vX.Y.Z...HEAD` and add `[X.Y.Z]` pointing at `releases/tag/vX.Y.Z`.

   A section with nothing under it fails the suite, and so does a version that is no semantic version or an entry with no date.

4. Verify, and prove the archives build before a tag exists.

   ```sh
   make fmt
   make verify
   make release-snapshot
   ```

   `make release-snapshot` runs the same releaser as the release path with publishing, signing and the bill of materials skipped, and leaves the archives under `dist/`.
   Extract one and check that `socrates`, `scrat`, `agora`, `LICENSE` and `NOTICE` are all in it.

   Never run `make release` locally.
   That is the publishing path, and it is the workflow's to run.

5. Commit the changelog and push it.

   ```sh
   git commit -am 'Record vX.Y.Z in the changelog'
   git push origin main
   ```

## Tag and push

The tag is what starts everything, so push the commit first and the tag second.

```sh
git tag -a vX.Y.Z -m 'vX.Y.Z'
git push origin vX.Y.Z
```

The version a binary reports comes from `git describe --tags --exact-match` at build time, so the tag has to be on the commit being built.
A build that is not at a tag says it is a development build rather than borrowing the nearest one.

## Watch it through

```sh
gh run list --repo socrates-tools/socrates-review --workflow release.yml --limit 1
gh run watch --repo socrates-tools/socrates-review <run-id>
```

The run has two jobs, and they mean different things when they fail.

- **`verification` failed.**
  Nothing was published, and the tag names a commit that does not pass.
  Fix it on `main`, retire the tag, and tag again at the new commit.

- **`publication` failed.**
  Fix the cause, and whether the version can be reused depends on what was published before it failed.

Either way, ask before touching the tag.

## Retiring a tag

A tag is retired in two steps, and the first one decides whether the second is allowed at all.

1. Establish that the version is unspent.

   ```sh
   gh release view vX.Y.Z --repo socrates-tools/socrates-review --json tagName,isDraft,assets
   ```

   A tag with a published release is spent.
   Do not delete it, do not re-tag it, and do not edit the release by hand: somebody may already have installed it, and a tag that moves makes their checksum a lie.
   Cut the next patch version instead.
   Only a command that fails with a release not found is a version that can be tagged again.

2. Ask the person, then delete.

   Deleting a remote tag is not yours to decide.
   Say which tag, say what step 1 found, and wait for them to say so.

   ```sh
   git tag -d vX.Y.Z && git push origin :refs/tags/vX.Y.Z
   ```

   Never widen this to a pattern, a glob or a loop over tags.
   One named tag per invocation, and a tag you were not asked about is a tag you leave alone.

## Confirm what shipped

```sh
tag=$(gh release view --repo socrates-tools/socrates-review --json tagName --jq .tagName)
gh release view "$tag" --repo socrates-tools/socrates-review
```

The release should carry one `tar.gz` per supported platform, `checksums.txt`, a `checksums.txt.bundle` signature, and a bill of materials per archive.
Download one and run the verification the release footer states, then check that the binary names the tag:

```sh
socrates version
```

It should read `socrates vX.Y.Z, built from <commit>`.
A line saying it is a development build means the archive was not built at the tag, which is a release to throw away rather than one to explain.

Finally, update `## [Unreleased]` links if anything about the compare range changed, and tell the person what the release means rather than that the workflow went green.

## What this skill does not do

- It does not decide the version for a person who has not said what the release is for.
  Where the change list makes minor and patch genuinely arguable, ask.
- It does not write the changelog from `git log`.
  The changelog says what a release is for, and the commits already say what each change is.
- It does not edit a released section of `CHANGELOG.md`.
  A past release is what somebody already installed.
- It does not publish, sign, or attach anything by hand.
  If a release needs something the workflow does not produce, that is a change to `.goreleaser.yaml` or `.github/workflows/release.yml`, and it starts as a spec.
- It does not delete a tag, a release or an asset on its own say-so.
  Every one of those is somebody else's checkout, so it is the person's call and not this skill's.
