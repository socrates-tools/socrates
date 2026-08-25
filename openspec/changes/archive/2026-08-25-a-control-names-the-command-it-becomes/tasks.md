**The skill half asserts on the rendered document**, never on the constants it is assembled from.

**The loop half is driven end to end**, through the real `socrates` and the real `agora` against the simulator, and each journey asserts that nothing socrates ran wrote in the repository and that no command outside what the skill names was used.

**Nothing here asserts what a reviewer found.** The journeys require that a round arrived, that the handover happened and that the answer reached the run. A severity, an action or a count is the reviewer's and not this change's to require.

**Each guard is mutated on purpose before it is trusted:** the test that proves the new behaviour is run against the old code first.

**`make lint test` after every task, not at the end of a group.**

## 1. The surface check learns there are two binaries

- [x] 1.1 Attribute each command the skill names to the binary that offers it, rather than taking the word after one product name
- [x] 1.2 Resolve each against that binary and confirm every name in the current skill still passes. Covers `Every command named is real`
- [x] 1.3 Remove a flag from the reviewing binary and confirm the check fails naming it. Covers `A removed flag fails the suite`
- [x] 1.4 Name a surface-binary flag that does not exist and confirm the check fails the same way. Covers `A flag of the second binary is checked too`
- [x] 1.5 Mutate the attribution so the surface binary is unrecognised, and confirm 1.4 fails rather than skipping quietly

## 2. The bar and the page section

- [x] 2.1 Rewrite the page section so it names the local surface and names no host, upload or publish. Covers `The page is served locally and published nowhere`
- [x] 2.2 State that every control names the command it becomes and that nothing is recorded until the session runs it. Covers `Every control names the command it becomes`
- [x] 2.3 Keep the identifier, severity and audience sentences, and say they come from `socrates findings --json`. Covers `The page says what the person does back in the terminal`
- [x] 2.4 Confirm `PageBar` is still stated once and rendered from there, and that its five requirements are unchanged
- [x] 2.5 Confirm the decisions are still given in full and the rest listed briefly

## 3. The branches

- [x] 3.1 Rewrite the `escalated` branch to draft the document inline, serve it, hand over the address and wait, with the reason for drafting inline. Covers `The escalated branch drafts inline`
- [x] 3.2 Rewrite the `awaiting-response` branch to delegate the drafting, with the reason. Covers `The parked branch delegates`
- [x] 3.3 Say the delegate returns the document and the driving session serves and waits. Covers `The driving session owns the surface`
- [x] 3.4 Confirm both branches say to serve the page and hand over the address. Covers `Both branches say to serve the page`
- [x] 3.5 Confirm the session still takes the condition from what socrates reports. Covers `The session does not decide the condition itself`
- [x] 3.6 Confirm a parked round the author answers alone still says nothing about a page. Covers `A round the author answers alone gets no page`
- [x] 3.7 Confirm the bar still appears once, before both branches, and neither restates it. Covers `The bar is stated before the branches that refer to it`
- [x] 3.8 Confirm nothing about the page entered the waiting step. Covers `The page is not in the polling loop`

## 4. How an answer reaches the run

- [x] 4.1 Say the person answers on the page and the session records it with `socrates respond`. Covers `The skill says how an answer reaches the run`
- [x] 4.2 Confirm the skill does not ask the person to restate in the terminal what they answered on the page
- [x] 4.3 Confirm the skill still says the session does not answer a person's finding on its own behalf. Covers `A decision is still the person's to make`
- [x] 4.4 Confirm the identifier sentences survived, so a person naming a finding is still unambiguous. Covers `The skill says the page and the terminal use the same identifiers`
- [x] 4.5 Regenerate `.claude/skills/socrates/SKILL.md` with `make skill`, and confirm the suite fails while the two differ

## 5. The loop, driven

- [x] 5.1 Drive the `escalated` branch end to end: the round arrives, the document is served, the address is handed over, and the findings are not pasted. Covers `A round a person reads is handed over as a page`
- [x] 5.2 Drive the parked branch, answer a decision on the served page, and confirm the run holds nothing until `socrates respond` runs and holds it afterwards. Covers `A decision made on the page is recorded in the terminal`
- [x] 5.3 Confirm the reason recorded is the person's chosen option or their own words, and not a paraphrase. Covers `The person's own words reach the record`
- [x] 5.4 Confirm each finding served carries the identifier, severity and audience `socrates findings --json` reported. Covers `Every finding keeps the identifier socrates gave it`
- [x] 5.5 Confirm every journey left the repository under review byte-identical and used no command outside what the skill names
- [x] 5.6 Confirm nothing socrates spawns is asked to draft the page. Covers `Nothing socrates spawns drafts the page`

## 6. Closing out

- [x] 6.1 Confirm `socrates` gained no command, no flag, no dependency and no network reach, and that no run record carries anything about a page
- [x] 6.2 Confirm `internal/dialogue`, `internal/findings`, `internal/reviewer`, `internal/runstore`, `internal/report`, `internal/agent` and `internal/prompt` are untouched
- [x] 6.3 Run the full sequence: `make fmt`, `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop`, `make corpus-outcomes`
- [x] 6.4 Close issue 1, because the document no longer leaves the machine and the asymmetry it records is resolved rather than restated
