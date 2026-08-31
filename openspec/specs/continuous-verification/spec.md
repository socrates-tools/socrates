# continuous-verification Specification

## Purpose
TBD - created by archiving change what-continuous-integration-actually-proves. Update Purpose after archive.

## Requirements

### Requirement: Verification runs on every platform the project claims

An automated verification run SHALL run the suite on every supported platform.

A platform SHALL be claimed because the suite passed on it, never because it compiled for it.
Compiling proves that the code is well formed; the two mechanisms this tool depends on, reaping by process group and a seal the filesystem enforces, differ between platforms in ways only running finds.

The platforms run SHALL be read from the stated supported set rather than written into the workflow, so that no third copy of that set exists to disagree with the other two.

#### Scenario: The suite runs on each supported platform

- **WHEN** a verification run completes
- **THEN** the suite has run on every platform the supported set names, and on no platform it does not

#### Scenario: A platform in the matrix but not in the supported set fails

- **WHEN** a platform is added to what verification runs and is not named by the supported set
- **THEN** the suite fails and names the platform

### Requirement: What a workflow does changes only when this repository changes

Every action a workflow uses SHALL be named by an immutable commit, with its human-readable version beside it for a reader.

A tag is a pointer somebody outside this repository can move.
An unpinned action means the behaviour of a verification run can change with no commit here, which makes a passing run evidence about a workflow nobody in this repository chose.

A check SHALL fail while any action is named by a tag, because a pin held only by habit lapses on the first hurried change.

Pins SHALL be kept current by something that opens a change for a person to review, so that pinning is a decision to see every update rather than a decision to fall behind.

#### Scenario: An action named by a tag fails the check

- **WHEN** a workflow names an action by a tag rather than by a commit
- **THEN** the check fails and names the workflow and the action

#### Scenario: A dependency update arrives as a change somebody reviews

- **WHEN** an action, a Go module or a client dependency has a newer version
- **THEN** a change is opened proposing it, rather than the update arriving silently or not at all

### Requirement: A workflow asks for the least authority it needs

Every workflow SHALL state its token's permissions, and SHALL state read-only at the top.

A job that needs more SHALL ask for it in that job, so that what a workflow may do is readable in the file rather than inherited from a repository setting nobody reviews.

#### Scenario: Verification runs read-only

- **WHEN** a verification workflow runs
- **THEN** its token has read access and no more, stated at the top of the workflow

#### Scenario: A superseded run does not report

- **WHEN** a new commit arrives on a branch whose verification run is still going
- **THEN** the superseded run is cancelled, so what is reported is the state of the branch now

### Requirement: The suite's reach is checked, and the page half cannot vanish

The module graph SHALL be checked for known vulnerabilities on every verification run and on a schedule, so that an advisory published between runs becomes visible without anybody looking.

The browser that drives the page SHALL be installed by the run at a known version rather than found on the runner.
A browser resolved from whatever the runner image ships is a dependency this repository does not control, and its removal would take the half of the suite that drives a real page with it.

An absent browser SHALL remain a failure, never a skip.

#### Scenario: A known vulnerability in a dependency is reported

- **WHEN** a dependency in the module graph has a published advisory
- **THEN** the verification run reports it, and reports it on a schedule as well as on a run

#### Scenario: The page half runs against a known browser

- **WHEN** a verification run drives the surface page
- **THEN** it drives a browser the run installed at a known version, and a run with no browser fails rather than passing

### Requirement: A contributor and a verification run do the same thing

The full verification sequence SHALL be one entry point, and both a contributor and an automated run SHALL use it.

A workflow that lists the gates and a document that lists them again is a list that will disagree, and the disagreement appears as a gate added in one place and missed in the other.

The cheapest gate SHALL stay the first to fail, so that a formatting error does not wait behind the end-to-end suite.

#### Scenario: One command runs every gate

- **WHEN** the full verification is run
- **THEN** one entry point runs every gate, and the automated run invokes that same entry point

#### Scenario: Formatting fails before the suite does

- **WHEN** verification is run against a tree that is not formatted
- **THEN** it fails on formatting before the end-to-end suite is run

### Requirement: Nothing in continuous integration reaches the real reviewer CLI

No automated run SHALL compile or execute anything that reaches the real reviewer CLI.

This is held by build tags rather than by workflow configuration, so that a workflow written later inherits it rather than having to remember it.

A run that reached the real CLI would make the suite depend on a third party's availability, on a paid account, and on model output, and would make a green run evidence about something other than this repository.

#### Scenario: A verification run reaches no reviewer CLI

- **WHEN** a verification run completes
- **THEN** nothing it compiled or executed reached the real reviewer CLI
