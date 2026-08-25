**The skill half asserts on the rendered document**, never on the constants it is assembled from.

**Nothing here asserts that a diagram was drawn.** What a session chooses to illustrate is judgement, and a test that claimed a page carried a picture would read as passing while checking nothing. What is driven is the instruction the skill carries.

**Nothing here re-tests the surface.** That a fragment survives neutralising and that one outside the allowlist is stripped are already required of `surface-document` and already driven by a real browser. This change teaches what is already true, so a new test of it would be a second copy of an existing guard.

**Each guard is mutated on purpose before it is trusted:** the test that proves the new behaviour is run against the old prose first.

**`make lint test` after every task, not at the end of a group.**

## 1. The skill says how an illustration reaches the page

- [ ] 1.1 Say that a decision may carry one illustration, travelling as a fragment of its own beside the evidence rather than inside it. Covers `The skill names the fragment as the way an illustration reaches the page`
- [ ] 1.2 Say the fragment is inline SVG the session writes, and confirm the rendered skill names no diagram engine and implies none the surface would run. Covers `The skill says what form the illustration takes`
- [ ] 1.3 Say the fragment is reduced to a fixed allowlist before it is embedded, naming what does not survive it. Covers `The skill says what is stripped before the fragment is embedded`
- [ ] 1.4 Place the mechanism beside the bar's line about diagrams, and confirm the bar's judgement about when one is warranted still appears exactly once. Covers `How to draw one is stated beside when to draw one`
- [ ] 1.5 Mutate the bar's diagram line, and confirm the test that holds the bar to one statement fails rather than the mechanism's alone

## 2. The generated copy agrees with its source

- [ ] 2.1 Regenerate `.claude/skills/socrates/SKILL.md` with `make skill`, and confirm `make test` fails while the committed copy and `internal/skill` disagree
- [ ] 2.2 Confirm the skill still names nothing the binaries do not offer, and that the prohibitions it is held to are unchanged

## 3. The whole gate

- [ ] 3.1 `make fmt`, `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop`, `make corpus-outcomes`
