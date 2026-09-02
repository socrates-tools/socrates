package dialogue

import (
	"github.com/ktei/socrates/internal/findings"
	"github.com/ktei/socrates/internal/runstate"
)

const Cap = 2

const CapIsForContextRatherThanCost = "The session being handed findings is already large from implementing the change, and each round adds findings, patches and rebuttals to it. Round three is where that session starts compacting mid-review and quietly loses the thing it was defending. A cap chosen for a reason that is not cost is a cap nobody raises the moment they meet a run that nearly finished."

func After(classified findings.Classification, round int) runstate.State {
	if classified.Verdict() == runstate.Clean {
		return runstate.Clean
	}
	if len(BeyondAnswering(classified)) > 0 || round >= Cap {
		return runstate.Escalated
	}
	return runstate.AwaitingResponse
}

func AwaitsAPerson(classified findings.Classification, state runstate.State) bool {
	switch state {
	case runstate.Escalated:
		return true
	case runstate.AwaitingResponse:
		return len(AnsweredByAPerson(classified)) > 0
	default:
		return false
	}
}

func Answerable(classified findings.Classification) []findings.Finding {
	return append(AnsweredByAPatch(classified), AnsweredByAPerson(classified)...)
}

func AnsweredByAPatch(classified findings.Classification) []findings.Finding {
	return blocking(classified.Patchable)
}

func AnsweredByAPerson(classified findings.Classification) []findings.Finding {
	return blocking(classified.NeedingAPerson)
}

func BeyondAnswering(classified findings.Classification) []findings.Finding {
	return blocking(classified.NoOps)
}

func blocking(held []findings.Finding) []findings.Finding {
	var blocks []findings.Finding
	for _, finding := range held {
		if finding.Severity.Blocks() {
			blocks = append(blocks, finding)
		}
	}
	return blocks
}

func EscalatedByTheCap(classified findings.Classification, round int) bool {
	return After(classified, round) == runstate.Escalated && len(BeyondAnswering(classified)) == 0 && len(Answerable(classified)) > 0
}

func CapReached(round int) bool {
	return round >= Cap
}
