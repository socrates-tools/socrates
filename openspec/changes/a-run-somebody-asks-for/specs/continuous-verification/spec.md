## ADDED Requirements

### Requirement: A verification run happens when somebody asks for one

A verification run SHALL be startable on request, and a release SHALL gate on one.

It SHALL NOT start on every push or on every pull request.
A private repository bills a macOS minute as ten, and the supported set names macOS, so a run costs about two orders of magnitude more than the work of pushing; a project that cannot pay for a run per push and schedules one anyway gets no runs at all, the release among them.

What a run proves SHALL be unchanged by this.
The same entry point, the same matrix read from the same statement, the same pins, the same tokens, the same browser, the same boundary. Only when a run happens is decided here.

What this gives up SHALL be stated rather than absorbed: nothing reports that a commit on the default branch is broken, and nothing reports that an advisory has been published, until somebody asks for a run.
The one thing holding that is the entry point a contributor runs before reporting work done, which is the same command a run invokes.

#### Scenario: A run is started by a person or by a release

- **WHEN** what starts a verification run is read
- **THEN** it can be started on request and by a release, and no push or pull request starts one

#### Scenario: A release still gates on the full matrix

- **WHEN** a release runs
- **THEN** it runs the suite on every platform the supported set names, and publishes nothing unless that passed

## MODIFIED Requirements

### Requirement: The suite's reach is checked, and the page half cannot vanish

The module graph SHALL be checked for known vulnerabilities on the same footing as the rest of verification: on request, rather than on a schedule of its own.

An advisory published against the module graph therefore waits for somebody to ask.
That is a loss, and it is the same loss the rest of this change makes for the same reason: a run that cannot be paid for is not a run that happens.

The browser that drives the page SHALL be installed by the run at a known version rather than found on the runner.
A browser resolved from whatever the runner image ships is a dependency this repository does not control, and its removal would take the half of the suite that drives a real page with it.

An absent browser SHALL remain a failure, never a skip.

#### Scenario: A known vulnerability in a dependency is reported

- **WHEN** a dependency in the module graph has a published advisory and a check is asked for
- **THEN** the check reports it

#### Scenario: The page half runs against a known browser

- **WHEN** a verification run drives the surface page
- **THEN** it drives a browser the run installed at a known version, and a run with no browser fails rather than passing
