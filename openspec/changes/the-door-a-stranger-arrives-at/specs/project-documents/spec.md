## ADDED Requirements

### Requirement: A reader learns what socrates needs before they learn what it does

What a person reads first SHALL state what has to be true before socrates works, before it explains how socrates works.

That SHALL include the reviewer CLI, its version floor, and the account behind it.
For many readers that fact ends the evaluation, and a document that reaches it in passing after several screens has spent their time to look more available than it is.

It SHALL include the supported platforms and the toolchain a source install needs.

#### Scenario: The requirements come before the explanation

- **WHEN** the front document is read from the top
- **THEN** it states the reviewer CLI, its version floor, the account it needs, the supported platforms and the toolchain, before it explains how socrates works

### Requirement: What a document states about the build is checked, not maintained

A check SHALL fail while the toolchain floor a document states and the floor the module requires disagree.

This repository has already had that disagreement, with a document naming one version and the module requiring a later one, which means a reader following the document could not build.

The shape SHALL be the one already used for the port a person is told to forward: the document reads the fact from where it is defined rather than writing a second copy of it.

#### Scenario: A stated toolchain floor that disagrees with the module fails

- **WHEN** the floor a document states and the floor the module requires disagree
- **THEN** the check fails and names both

### Requirement: Every binary a person installs is explained

Every binary that installing socrates puts on a person's path SHALL be explained in what a person reads.

A second binary that appears with no explanation is indistinguishable from something a person did not consent to installing.

#### Scenario: The second binary is explained

- **WHEN** the front document is read
- **THEN** it says what the second binary is and what it is for, and not only which directory it writes to

### Requirement: The project says where what it reads goes

What a person reads SHALL state that socrates itself sends nothing anywhere, and SHALL state that the reviewer CLI it drives sends the change to that CLI's vendor.

Both halves SHALL be stated.
Stating only the first is a claim that reads as complete and is not, and a reader who later discovers the second has learned that the document flatters the project.

Silence is not neutral here: a tool that reads the repository somebody is working in and says nothing about where that goes will be assumed the worst of by exactly the careful readers it most wants.

#### Scenario: The reach is stated in both directions

- **WHEN** the front document and the security policy are read
- **THEN** both state that socrates sends nothing itself, and that the reviewer CLI it drives sends the change to its vendor

### Requirement: A vulnerability has somewhere private to go

A security policy SHALL state where to report privately, and private reporting SHALL be enabled on the repository.

It SHALL state what socrates defends against and what it does not, so that a reporter can tell whether what they found is in scope before spending their time and the maintainer's.

What it does not defend against SHALL include the residual already stated in what a person reads: a hook an administrator's policy configures, and whatever that hook's command reaches outside the sealed worktree.

#### Scenario: A reporter has a private channel

- **WHEN** somebody finds a vulnerability
- **THEN** the security policy names a private channel, and that channel is enabled on the repository

#### Scenario: The policy states what is out of scope

- **WHEN** the security policy is read
- **THEN** it names what socrates defends against and what it does not, including the residual on a machine whose policy configures hooks

### Requirement: A contributor learns how a change is proposed here before they write one

The contributing guide SHALL state, before anything else, that a change starts as a proposal in the specifications directory.

A contributor who learns the toolchain first writes code first, and the change that arrives has to be turned away for a rule nothing told them.
This project has 27 archived changes and 17 promoted capabilities, so the rule is not incidental to it.

The guide SHALL also state the toolchain, the single command that runs every gate, and the conventions the source is held to.

It SHALL state that pinned dependencies produce a steady stream of small changes, so that a long queue of them reads as the mechanism working rather than as neglect.

#### Scenario: The guide leads with the proposal

- **WHEN** the contributing guide is read from the top
- **THEN** it states that a change starts as a proposal in the specifications directory, before it states the toolchain

#### Scenario: The guide names one verification command

- **WHEN** the contributing guide states how to check a change
- **THEN** it names the single entry point that runs every gate, rather than listing the gates a second time

### Requirement: The documents a contributor and a maintainer need exist

A code of conduct with a working contact, a support document, stated ownership, and issue and pull-request templates SHALL exist.

The repository SHALL carry the editor and attribute settings that keep a contributor's tooling from fighting its conventions, so that a first change is not a diff of line endings and whitespace.

#### Scenario: A first-time contributor finds the documents they need

- **WHEN** the repository is read by somebody about to open their first change
- **THEN** a code of conduct with a contact, a support document, stated ownership, and issue and pull-request templates are present

#### Scenario: A contributor's editor does not fight the conventions

- **WHEN** a contributor edits a file with tooling that reads the repository's editor and attribute settings
- **THEN** those settings match the conventions the source is held to, so no change arrives as a whitespace or line-ending diff
