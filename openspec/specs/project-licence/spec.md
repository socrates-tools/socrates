# project-licence Specification

## Purpose
The terms the work is published under, and the notices the work carries for everything a shipped binary embeds.

## Requirements

### Requirement: The work states the terms it is published under

The repository SHALL carry a file at its root stating the terms the work is published under.
Those terms SHALL be the MIT licence, and the copyright line SHALL name the holder.

The terms SHALL be stated in that one file.
No source file SHALL carry a licence header, because a header is a comment and this repository's source carries none, and a header would say on every file what the root file already says once.

What is forbidden is a header stating the work's own terms.
A notice for a work socrates borrowed is required by the requirement below, and the check SHALL NOT catch one.
The repository SHALL recognise a borrowed notice by its citing the notice file, because a rule that cannot tell the two apart makes two requirements contradict each other in whichever file has to carry both.

A reader SHALL be able to learn the terms without running anything, reading a second file, or asking.

#### Scenario: The root states the terms

- **WHEN** the repository is read
- **THEN** a file at its root states the MIT licence in full, with a copyright line naming the holder

#### Scenario: No source file carries a licence header

- **WHEN** every source file in the repository is read
- **THEN** none opens with a licence header, and the check that holds this fails when one is added

#### Scenario: A borrowed work's notice is not the header this forbids

- **WHEN** a source file opens with a comment carrying a borrowed work's notice and citing the notice file
- **THEN** the check passes for it, and still fails for a header stating the work's own terms

### Requirement: Everything embedded carries the notice its own licence requires

Everything embedded in a shipped binary that socrates did not author SHALL be named in a notice file at the root, with its licence and that licence's full text.

A reference SHALL NOT stand in for the text.
The licences these works are published under require the text to accompany the copy, and a link is not the text: it is a request that the reader go and find it, which fails exactly when the copy has travelled furthest from the repository.

The notice SHALL be a condition of what the binary distributes rather than a description of the source tree, so a work present in the source and absent from the binary SHALL NOT need an entry, and a work absent from the source and present in the binary SHALL.

#### Scenario: An embedded font's licence travels with the work

- **WHEN** the notice file is read
- **THEN** it names each font family embedded through the assets package, states its licence, and carries that licence's full text

#### Scenario: An embedded library's licence travels with the work

- **WHEN** the notice file is read
- **THEN** it names each third-party library bundled into the embedded client, states its licence, and carries that licence's full text

### Requirement: The shipped bundle keeps the notices of what it bundles

The bundler SHALL NOT delete legal comments from the client bundle it produces.

Where a bundled work ships no notice of its own, the notice SHALL be written where the bundler collects it from, so that what the bundler is forbidden from deleting exists to be deleted.
A setting that preserves notices proves nothing while there is no notice to preserve.

What is asserted SHALL be the bytes that ship rather than the setting that produced them.
A test on the setting passes for a bundler that was configured correctly and then failed, and the thing the licence attaches to is the artifact, not the configuration.

The notices SHALL be collected rather than scattered, so that a person auditing the bundle has one place in it to read.

#### Scenario: The shipped bundle carries the notice of what it bundles

- **WHEN** the embedded client bundle is read as it ships
- **THEN** it carries the copyright notice of every third-party library bundled into it

#### Scenario: A bundler that deletes notices fails the suite

- **WHEN** the bundler is configured to delete legal comments and the bundle is regenerated
- **THEN** the suite fails, naming the notice that is missing from what ships

### Requirement: A new embedded asset cannot arrive without its notice

A check SHALL fail while anything embedded in a shipped binary is not accounted for in the notice file.

The check SHALL answer whether an entry was written and SHALL NOT attempt to judge whether the entry is legally sufficient.
A check that tried to decide sufficiency would be one nobody could trust and everybody would work around, and the judgement belongs to the person adding the asset.

#### Scenario: An embedded asset with no entry fails the check

- **WHEN** a file is added to what the assets package embeds and no entry naming it is added to the notice file
- **THEN** the check fails and names the file

#### Scenario: An asset in the source but not in the binary needs no entry

- **WHEN** a file is present in the client source tree and is not embedded in any shipped binary
- **THEN** the check passes with no entry for it
