package runstate

import "github.com/socrates-tools/socrates/internal/exit"

func (s State) Code() int {
	switch s {
	case Clean:
		return exit.NothingAsked
	case Escalated:
		return exit.SomethingAsked
	case Failed, Cancelled:
		return exit.Failed
	default:
		return exit.Working
	}
}
