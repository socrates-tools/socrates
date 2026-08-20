package dialogue

import (
	"github.com/rui-ktei/socrates/internal/findings"
	"github.com/rui-ktei/socrates/internal/runstate"
)

const Cap = 2

const CapIsForContextRatherThanCost = "The session being handed findings is already large from implementing the change, and each round adds findings, patches and rebuttals to it. Round three is where that session starts compacting mid-review and quietly loses the thing it was defending. A cap chosen for a reason that is not cost is a cap nobody raises the moment they meet a run that nearly finished."

func After(classified findings.Classification, round int) runstate.State {
	if len(classified.NeedingAPerson) > 0 {
		return runstate.Escalated
	}
	if classified.Verdict() == runstate.Clean && len(classified.Patchable) == 0 {
		return runstate.Clean
	}
	if round >= Cap {
		return runstate.Escalated
	}
	return runstate.AwaitingResponse
}

func CapReached(round int) bool {
	return round >= Cap
}
