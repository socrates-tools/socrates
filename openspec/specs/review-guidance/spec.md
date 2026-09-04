# review-guidance Specification

## Purpose

A project's own conventions are review criteria, and a reviewer that never hears them reports the project's house style as defects and passes what the project actually forbids.
This capability is the one channel by which a reviewed repository tells socrates' reviewer what it holds itself to, on terms that keep the reviewer independent of the change it is reading.

## Requirements

### Requirement: One file in the repository under review carries the project's guidance

socrates SHALL read `.socrates/review-guide.md`, relative to the repository root, and SHALL deliver what it holds to the reviewer.

That path SHALL be fixed.
Nothing SHALL configure it, name a second file, name a directory of files, or scope a file to a subset of the change.
A configurable path would be a way to point socrates at a file the repository's own tooling generates, which is reading the working tree wearing a hat.
A repository that wants a reviewer to know something writes it in the one file, and a reviewer is told the same thing about every change to that repository.
`no-mistakes` scopes its guidance by glob because a contributor there can delete a maintainer's rule from their own review by widening an ignore list, and socrates has no ignore list for that attack to use.

It SHALL be the only file socrates reads out of a reviewed repository on the reviewer's behalf.
socrates SHALL NOT write in that directory, or anywhere else in the repository under review, and its presence there SHALL NOT change what `git status` reports.

Where the file is absent, or holds nothing once it has been read and neutralised, the review SHALL proceed with no guidance and SHALL NOT be refused for it.

#### Scenario: A project's guidance reaches the reviewer

- **WHEN** a review runs against a repository whose base commit carries `.socrates/review-guide.md`
- **THEN** the delivered prompt carries what that file says

#### Scenario: A repository with no guide

- **WHEN** a review runs against a repository that carries no such file
- **THEN** the review proceeds, and the delivered prompt carries no guidance section and no heading for one

#### Scenario: A guide holding nothing

- **WHEN** the file is present at the base commit and holds only whitespace
- **THEN** the review proceeds as though there were no file, and the delivered prompt carries no guidance section

#### Scenario: Nothing configures the path

- **WHEN** a configuration file names a guidance path or a second guidance file
- **THEN** socrates refuses it, naming the field it does not have

#### Scenario: Reading the guide writes nothing

- **WHEN** a review runs against a repository carrying a guide
- **THEN** `git status` in that repository is byte-identical before and after, and nothing under `.socrates/` in it was written

### Requirement: The guidance is the version at the run's base commit

socrates SHALL read the file at the run's base commit.
It SHALL NOT read it from the working tree, from the snapshot commit, or from the snapshot worktree.

The guidance that reviews a change is therefore the guidance as it stood before the change.
A guide the change under review adds, edits or deletes governs the next run rather than this one.
Without this the session that wrote the code also writes the rules the code is judged by, in the same change, and the independence the whole design pays for is spent in one file.
The base commit is the anchor `no-mistakes` gets from a trusted default branch, and socrates has it already.
It is the state the change is measured against, and by construction no part of the change authored it.

Every round of a run SHALL read the same base, which is the run's original base.

The ceiling a run is held to SHALL be resolved when the run begins, recorded on the run, and read from there by every later round, rather than loaded again from configuration each round.
A fixed base alone does not hold the guidance still.
A ceiling read afresh would let a person raise it while a run is parked and hand the next round guidance the first round was never given, while the run's own record still said the guide was over the ceiling.
Both the text and the bound on it therefore come from the run.

The guidance is fixed for the whole of a run and cannot change between the rounds of a dialogue.

#### Scenario: A guide the change adds does not govern that change

- **WHEN** a review runs against a repository whose base commit carries no guide and whose working tree adds one
- **THEN** the delivered prompt carries no guidance section, and the guide is not read from the working tree

#### Scenario: A guide the change edits is delivered as the base had it

- **WHEN** the base commit carries a guide and the working tree changes what it says
- **THEN** the delivered prompt carries what the base commit's version says, and none of what the change made it say

#### Scenario: The guidance does not move between rounds

- **WHEN** a second round is invoked on a run whose repository changed the guide after the run began
- **THEN** the second round's delivered prompt carries the same guidance the first round's did

#### Scenario: The ceiling a run began with is the one every round is held to

- **WHEN** a run's first round is held to a ceiling that leaves the guide out, the person raises that ceiling while the run is parked, and a later round is invoked
- **THEN** the later round is held to the ceiling the run began with and is given what the first round was given

### Requirement: The guidance adds to socrates' rules and removes none of them

The guidance SHALL be delivered as something the project says about its own code, added to the rules socrates delivers.

It SHALL NOT be able to withdraw an instruction the prompt states, redefine a severity, an action or a risk level, alter the verdict rule or the round cap, exempt a file or a directory from review, or ask that a finding be withheld.
The delivered prompt SHALL say so in socrates' own words, in the guidance section itself, so that a guide asking for any of those is asking the reviewer for something the prompt has already refused.
A guide that could switch a rule off would make `clean` mean a different thing in each repository, and the first rule written under a deadline is the one that removes the finding nobody wants.

socrates SHALL NOT decide anything from the guidance.
It is not parsed, not matched and not compared; no state, verdict, exit code, severity, action or classification depends on a word in it.
What socrates does with it is carry it, bounded and neutralised, into one section of one prompt.

#### Scenario: The section says what the guidance cannot do

- **WHEN** a review runs against a repository carrying a guide
- **THEN** the delivered prompt states, in socrates' own text beside the guidance, that it adds to the rules above and withdraws none of them, and that it does not decide severity, action, or whether a finding is reported

#### Scenario: The guidance comes after the rules it cannot remove

- **WHEN** the delivered prompt's sections are read in order
- **THEN** the guidance section appears after socrates' own reviewing rules and the test-quality block, and before the reporting contract

#### Scenario: The guidance decides nothing socrates decides

- **WHEN** two reviews run against the same change with the same reviewer answer, one against a repository whose guide demands a severity, demands an action and asks that a finding be withheld, and one against a repository carrying no guide
- **THEN** the verdict, the classification of every finding, what a person is asked and the exit code are the same in both

### Requirement: The guidance is neutralised before it enters the prompt

The guidance SHALL be neutralised on the way into the prompt, on the mechanism socrates already applies to model output, and SHALL NOT be given a second mechanism of its own.

It is text from outside socrates entering a prompt, the same class of input as a finding socrates re-embeds.
The reason is the same: without it the text can present itself as socrates' own instruction.
Neutralisation SHALL leave a guide unable to open a section of the delivered prompt, because a section is how the reviewer tells socrates' instructions apart from everything socrates is quoting.

The file on disk SHALL NOT be rewritten, because socrates never writes in the repository under review.

What is recorded about the guidance and what is delivered SHALL be decided from the same text, and that text SHALL be the neutralised one.
Neutralising can reduce a guide that is not empty to one that is, because a line of conflict markers becomes spaces and a bare heading marker becomes nothing.
Two emptiness tests, one on what was read and one on what is delivered, therefore let a run record that a project's guidance was carried to a reviewer given none of it.
That is the same shape as a review that returned nothing because it errored.

#### Scenario: A guide that neutralises to nothing is no guidance

- **WHEN** the base commit carries a guide whose every line neutralises away
- **THEN** the delivered prompt carries no guidance section, and the run records that there was no guidance rather than that guidance was carried

#### Scenario: A guide cannot open a section

- **WHEN** the guide's text carries lines that would begin a prompt section
- **THEN** the delivered prompt's sections are the ones socrates assembled, and the guide's text sits inside the guidance section

#### Scenario: A guide carrying conflict markers

- **WHEN** the guide's text carries conflict markers
- **THEN** the delivered prompt carries the guidance with them removed, on the same mechanism a re-embedded finding goes through

### Requirement: A guide over the ceiling is not carried, and the review still runs

The guidance SHALL be held to a fixed byte ceiling, measured on the size the base commit records for the file and before any of it is read.

A guide over the ceiling SHALL NOT be carried, and SHALL NOT be shortened, summarised or delivered in part.
Half a maintainer's rules delivered silently is a reviewer told rules that do not exist as stated, so the guide is carried whole or not at all.

The review SHALL proceed with no guidance rather than being refused.
Refusing would stop every review in that repository until somebody committed a smaller guide unreviewed, including the review of the change that shrinks it.
One person's oversized guide would then hold up everybody else's work in progress.
The ceiling exists to bound what the reviewer's attention is spent on, not to stop reviews.

The run SHALL record that the guide was over the ceiling, distinctly from there being no guide, and the person SHALL be told its size, the ceiling, and that the next review carries it once it is smaller.

The ceiling SHALL default below the diff ceiling, because guidance competes with the change for the reviewer's attention and the change is what the review is for.

It SHALL be settable by the person, in socrates' own configuration file, and SHALL NOT be settable by a repository.
The file that states it lives outside every repository, so raising the ceiling is the decision of whoever runs the review rather than of the thing being reviewed, and the anchor is untouched by it.
What a stated ceiling buys is a project whose conventions genuinely need more room, and the worst it costs is a person spending more of their own reviewer's attention on their own conventions.

A stated ceiling SHALL apply to every repository and every run, on the footing the diff ceiling already has, and a repository file naming one SHALL reach nothing.
A ceiling that is not a size SHALL be refused with the rest of the configuration rather than quietly replaced by the default.

#### Scenario: The ceiling is the person's to state, not the repository's

- **WHEN** the person's configuration states a guidance ceiling and a file in the repository under review states a different one
- **THEN** the guide is held to what the person stated, and what the repository states reaches nothing

#### Scenario: An unstated ceiling is the default

- **WHEN** no configuration states a guidance ceiling
- **THEN** the guide is held to the default, and a guide over it is not carried

#### Scenario: A guide over the ceiling is not carried

- **WHEN** a review runs against a repository whose base commit carries a guide larger than the ceiling
- **THEN** the review proceeds, the delivered prompt carries no guidance section, and no part of the guide reaches the reviewer

#### Scenario: A guide over the ceiling is reported to the person

- **WHEN** a guide over the ceiling is met
- **THEN** the run records that it was over the ceiling rather than that there was none, and a person reading the run is told the file, its size, the ceiling, and that the next review carries it

#### Scenario: A guide at the ceiling is carried whole

- **WHEN** the guide is exactly the ceiling in size
- **THEN** the review proceeds and the delivered prompt carries all of it

#### Scenario: An oversized guide is never shortened

- **WHEN** a guide over the ceiling is met
- **THEN** no prompt is delivered carrying part of it

### Requirement: What happened with the guidance is recorded in the run

The run SHALL record which of three things happened: guidance was carried, there was none, or the change carries a guide the base commit does not.

The third SHALL be recorded distinctly from the second, and SHALL reach the person reading the run.
It is what somebody sees the first time they write the file, and without it the base-commit anchor is a silent surprise rather than a rule.
It SHALL NOT fail the round, because writing the guide and reviewing the change that writes it is an ordinary thing to do.

The third SHALL be answered only where the base commit carried nothing at that path.
A change that removes a guide, or empties one, also names that path among the files it changes, and telling the person to commit a guide so the next review reads it is advice about a file they deleted.

#### Scenario: Guidance carried is recorded

- **WHEN** a round runs with guidance delivered
- **THEN** the run records that the project's guidance was carried

#### Scenario: A newly written guide is reported to the person

- **WHEN** a review runs against a repository whose change adds a guide the base commit does not carry
- **THEN** the run records that a guide is present in the change and did not reach this reviewer, a person reading the run is told so, and the round is not failed for it

#### Scenario: A guide the change removes is not a guide it adds

- **WHEN** the base commit carries a guide and the change deletes it
- **THEN** the run does not record it as a guide present in the change, and nobody is told to commit a file they removed

#### Scenario: No guide is recorded as no guide

- **WHEN** a review runs against a repository carrying no guide at the base and none in the change
- **THEN** the run records that there was no guidance, distinctly from a guide that did not reach the reviewer
