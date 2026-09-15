package exit

import "github.com/socrates-tools/socrates/internal/runstate"

const (
	Clean     = 0
	Escalated = 1
	Failed    = 2
	Usage     = 3
	Working   = 4
)

func For(state runstate.State) int {
	switch state {
	case runstate.Clean:
		return Clean
	case runstate.Escalated:
		return Escalated
	case runstate.Failed, runstate.Cancelled:
		return Failed
	default:
		return Working
	}
}
