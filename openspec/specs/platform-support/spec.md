# platform-support Specification

## Purpose
TBD - created by archiving change the-platforms-the-seal-holds-on. Update Purpose after archive.

## Requirements

### Requirement: The platforms the guarantees hold on are stated in one place

The platforms socrates supports SHALL be stated in one place in the source, and everything that needs to know SHALL read it from there.

A platform SHALL be supported when reaping by process group is available on it, the filesystem enforces the seal on it, and the suite runs on it.
Compiling SHALL NOT be the basis, because both guarantees compile everywhere and hold in fewer places than that.

The statement SHALL be read by the refusal, by the release, and by what a person reads, rather than written once for each.
A set written three times is a set that disagrees with itself, and this repository has already had one such disagreement between the Go version its build requires and the Go version its README states.

#### Scenario: The supported set is stated once

- **WHEN** the supported platforms are needed by the refusal, by what a release builds, and by what a person reads
- **THEN** all three read one stated set, and no second copy of it exists

#### Scenario: What a person reads matches what the source states

- **WHEN** the documented supported platforms are compared against the stated set
- **THEN** they agree, and the check that holds this fails when they do not

### Requirement: A review is refused where a guarantee cannot hold

socrates SHALL refuse to start a review on a platform where reaping by process group is not available.

The refusal SHALL happen before any reviewer process is spawned and before any snapshot is taken, so that nothing has been started that cannot then be reaped.

The refusal SHALL name the guarantee that is missing and the way forward, so that a person who meets it is not left to search for what to do.

A verdict socrates cannot stand behind is worse than no verdict.
That is the reasoning the seal already refuses on, and the reason the guarantee is absent does not change it.
A platform that cannot reap is as unable to hold the promise as a filesystem that ignores permission bits.

#### Scenario: An unsupported platform refuses the review before anything starts

- **WHEN** a review is started on a platform where reaping by process group is not available
- **THEN** socrates refuses, names the missing guarantee and the way forward, spawns no reviewer, and takes no snapshot

#### Scenario: Reading a run is not refused on an unsupported platform

- **WHEN** a command that neither spawns nor reaps is run on an unsupported platform
- **THEN** it runs, because it makes neither promise, and a finished run's record is readable wherever it is read

### Requirement: A release cannot build a platform that is not supported

A test SHALL fail while what a release builds and the stated supported set disagree.

The test SHALL fail in both directions: a platform a release builds and the statement does not name, and a platform the statement names and no release builds.
A one-directional check passes on the day somebody adds a target and forgets the statement, which is the exact failure it exists to catch.

#### Scenario: A platform added to a release without being stated fails the suite

- **WHEN** a platform is added to what a release builds and is not added to the stated supported set
- **THEN** the suite fails and names the platform

#### Scenario: A platform stated but never built fails the suite

- **WHEN** the stated supported set names a platform that no release builds
- **THEN** the suite fails and names the platform
