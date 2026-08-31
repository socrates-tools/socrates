package outcomes

import (
	"fmt"
	"slices"
	"strings"

	"github.com/socrates-devtool/socrates/internal/corpus"
	"github.com/socrates-devtool/socrates/internal/dialogue"
	"github.com/socrates-devtool/socrates/internal/findings"
	"github.com/socrates-devtool/socrates/internal/runstate"
)

const DistributionPath = "internal/corpus/outcomes/distribution.txt"

type Reach struct {
	State    runstate.State
	Recorded []string
}

type Report struct {
	Entries int
	Reaches []Reach
}

func Of(entries []corpus.Entry) (Report, error) {
	reached := map[runstate.State][]string{}

	for _, entry := range entries {
		document, err := findings.Parse(entry.Document)
		if err != nil {
			return Report{}, fmt.Errorf(answerDoesNotParse, entry.Name, err)
		}
		if entry.Round < 1 {
			return Report{}, fmt.Errorf(answerNamesNoRound, entry.Name)
		}

		state := dialogue.After(findings.Classify(document), entry.Round)
		reached[state] = append(reached[state], entry.Name)
	}

	report := Report{Entries: len(entries)}
	for _, state := range Reachable() {
		names := reached[state]
		slices.Sort(names)
		report.Reaches = append(report.Reaches, Reach{State: state, Recorded: names})
	}
	return report, nil
}

func (r Report) Distribution() string {
	var built strings.Builder

	built.WriteString(whatSocratesWouldDo + "\n\n")
	for _, reach := range r.Reaches {
		fmt.Fprintf(&built, reachedBy+"\n", reach.State, len(reach.Recorded))
		names(&built, reach)
	}
	fmt.Fprintf(&built, "\n"+recordedAnswers+"\n", r.Entries)

	return built.String()
}

func (r Report) Render(naming map[runstate.State][]string) string {
	var built strings.Builder

	built.WriteString(whatSocratesWouldDo + "\n\n")
	for _, reach := range r.Reaches {
		fmt.Fprintf(&built, reachedByAndNamed+"\n", reach.State, len(reach.Recorded), len(naming[reach.State]))
		names(&built, reach)
		if len(reach.Recorded) == 0 {
			fmt.Fprintf(&built, onlyReachedByHand+"\n", reach.State)
		}
	}
	fmt.Fprintf(&built, "\n"+recordedAnswers+"\n", r.Entries)

	return built.String()
}

func names(built *strings.Builder, reach Reach) {
	for _, name := range reach.Recorded {
		fmt.Fprintf(built, answerLine+"\n", name)
	}
}

func Reachable() []runstate.State {
	pairs := shapes()

	seen := map[runstate.State]bool{}
	for held := range 1 << len(pairs) {
		var document findings.Document
		for position, pair := range pairs {
			if held&(1<<position) != 0 {
				document.Findings = append(document.Findings, pair)
			}
		}
		for round := 1; round <= dialogue.Cap; round++ {
			seen[dialogue.After(findings.Classify(document), round)] = true
		}
	}

	var reachable []runstate.State
	for _, state := range runstate.States() {
		if seen[state] {
			reachable = append(reachable, state)
		}
	}
	return reachable
}

func shapes() []findings.Finding {
	var pairs []findings.Finding
	for _, severity := range findings.Severities() {
		for _, action := range findings.Actions() {
			pairs = append(pairs, findings.Finding{Severity: severity, Action: action})
		}
	}
	return pairs
}
