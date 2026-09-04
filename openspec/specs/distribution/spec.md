# distribution Specification

## Purpose
TBD - created by archiving change a-version-somebody-can-install. Update Purpose after archive.

## Requirements

### Requirement: A build says which release it is, or says it is not one

A binary SHALL report a semantic version alongside the commit stamp it already reports.

The version SHALL be recorded when the binary is linked, from the tag, and SHALL NOT be a constant in the source.
A constant is a second place the version lives, and it is wrong for the whole window between a release and the next time somebody remembers to change it.

A build from a tree that is not at a tag SHALL say that it is a development build rather than report the tag nearest to it.
Reporting the nearest tag makes a build from an arbitrary commit indistinguishable from the release whose name it borrowed, which is precisely the distinction a bug report needs.

#### Scenario: A tagged build reports its version

- **WHEN** a binary built from a tagged commit is asked what it is
- **THEN** it reports that version, and the commit stamp it reported before

#### Scenario: An untagged build does not claim a version

- **WHEN** a binary built from a commit that is not at a tag is asked what it is
- **THEN** it says it is a development build, and names no release version

#### Scenario: The version is read where the binary was linked, not where it runs

- **WHEN** a binary is moved to another machine and asked what it is
- **THEN** it reports the version it was linked with, having read nothing where it runs

### Requirement: A release is verifiable by somebody who trusts nobody

Tagging SHALL produce, for each supported platform, an archive, a checksum for it, a signature over the checksums, a bill of materials, and build provenance.

Each answers a different question and none substitutes for another.
They are: whether the download is intact, whether the checksums are the ones this project published, what is inside a binary that embeds fonts and a client bundle, and which workflow from which commit produced it.

The signing SHALL require no private key held by a person, so that a release depends on no secret that can be lost or asked for.

A stranger SHALL be able to check a downloaded binary against this repository without extending trust to whoever pushed the tag.

#### Scenario: A release carries what is needed to verify it

- **WHEN** a tag is pushed and the release completes
- **THEN** every archive has a checksum, the checksums carry a signature, and a bill of materials and build provenance are published with them

#### Scenario: A release that did not pass verification produces nothing

- **WHEN** a tag is pushed on a commit whose verification does not pass
- **THEN** no archive, checksum, signature or provenance is published

#### Scenario: Nothing in the release path reaches the real reviewer CLI

- **WHEN** the release runs
- **THEN** it compiles and runs nothing that reaches the real CLI, exactly as continuous integration does not

### Requirement: A release builds the supported platforms and no others

A release SHALL build for every platform the supported set names, and for no platform it does not.

A cross-compiled target SHALL NOT be treated as a tested one.
What is built from one machine and what the suite ran on are different claims, and the second is what supported means.

#### Scenario: The release matrix is the supported set

- **WHEN** what a release builds is compared against the stated supported platforms
- **THEN** they agree, and the check that holds this fails when they do not

### Requirement: An archive gives a person what installing from source gives them

An archive SHALL contain both names the tool answers to and the second binary that ships with it, so that extracting one archive and installing from source arrive at the same place.

It SHALL also contain the terms the work is published under and the notices for everything it embeds, because those are conditions on the copy and an archive is a copy.

The two names SHALL come from the same build in the same release, so that they cannot be at different commits.

#### Scenario: An archive carries both names and the second binary

- **WHEN** a release archive is extracted
- **THEN** it contains the tool under both of its names and the second binary, and the two names report the same version and the same commit

#### Scenario: An archive carries the terms and the notices

- **WHEN** a release archive is extracted
- **THEN** it contains the licence the work is published under and the notices for everything embedded in it

### Requirement: What changed in a release is written for the person upgrading

Each release SHALL be recorded in a changelog, in terms of what it means to somebody using the tool.

The changelog SHALL be written rather than generated.
Commit messages in this repository say what a change is, by convention and deliberately.
A release note says what that change means to a person who has the previous version installed, which is a different document for a different reader.

#### Scenario: A release is recorded in the changelog

- **WHEN** a release is published
- **THEN** the changelog names that version and says what changed in it for somebody using the tool
