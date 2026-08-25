**The name is one constant per place it is needed**, and the tests that matter are the ones that would still pass if a second copy appeared, so each is written to fail on the copy rather than on the value.

**What a person reads is asserted from the binary's own output**, never from the source string it came from.

**The alias is driven through `make install`**, into a directory the test owns, so what is proved is the recipe a person runs.

**Each guard is mutated on purpose before it is trusted.**

**`make lint test` after every task, not at the end of a group.**

## 1. The binary

- [x] 1.1 Move the command package so the toolchain produces `socrates`, and point the Makefile, the harness and the scripts at it
- [x] 1.2 Confirm the surface is unchanged under the new name. Covers `The command surface is what it says`

## 2. The short name

- [x] 2.1 Install the short name beside the full one, in the directory the toolchain reports rather than one guessed here
- [x] 2.2 Drive `make install` into a directory the test owns and read the same commit from both names. Covers `Both names are installed and reach the same binary`
- [x] 2.3 Confirm the short name is a second name for one binary rather than a second binary
- [x] 2.4 Mutate the link out of the recipe and confirm the journey fails rather than skipping

## 3. What it prints

- [x] 3.1 Take the name from one place in each package that prints it, and rename what a person reads
- [x] 3.2 Confirm a refusal under the short name names the full one and names the short one nowhere. Covers `The short name reports the full one`
- [x] 3.3 Confirm no string socrates prints names the short name, over the whole module

## 4. The skill

- [x] 4.1 Name the skill from the one constant the installer and the surface check already share. Covers `One constant names the skill and reads its invocations`
- [x] 4.2 Install it under the product's name, and regenerate the committed copy rather than moving it. Covers `The skill is installed under the product's name`
- [x] 4.3 Confirm every invocation the skill names is the full command and that the surface check still resolves each one. Covers `The skill names the full command`
- [x] 4.4 Confirm the loop, the polling step, the single mention of `socrates watch` and the states it names all survive the rename

## 5. What the documents say

- [x] 5.1 Record both names and the layout change
- [x] 5.2 State the invariant that one binary carries both names and that what it prints names the full one
- [x] 5.3 Correct nothing else

## 6. Verification

- [x] 6.1 `make fmt && make lint test` clean
- [x] 6.2 `make skill` leaving the committed copy where the new name puts it, with the old directory gone
- [x] 6.3 `make scenario-coverage`, `make scenario-drop` and `make corpus-outcomes`
- [x] 6.4 `make install` and both names read back, against `git rev-parse --short=7 HEAD`

## What the verification recorded

`make install` linked one binary and put both names in `$GOPATH/bin`: `socrates`, and `scrat` pointing at it.
Both answered `socrates dfcab69, built from a working tree with uncommitted changes`, from the checkout at `dfcab69`.

Removing the link from the recipe was tried on purpose, and both journeys failed on the missing name rather than skipping.

The rename found one thing prose could not carry: with the product and the command spelled the same, the check that reads the skill's invocations reported `gives`, `runs`, `has`, `will` and `patches` as commands the binary lacks, one for each sentence whose subject is socrates.
That is why an invocation in the skill is now written as code and the check reads code alone.
