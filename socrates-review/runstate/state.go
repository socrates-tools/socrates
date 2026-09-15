package runstate

import "slices"

type State string

const (
	Starting         State = "starting"
	Reviewing        State = "reviewing"
	AwaitingResponse State = "awaiting-response"
	Clean            State = "clean"
	Escalated        State = "escalated"
	Failed           State = "failed"
	Cancelled        State = "cancelled"
)

var states = []State{Starting, Reviewing, AwaitingResponse, Clean, Escalated, Failed, Cancelled}

var terminalStates = []State{Clean, Escalated, Failed, Cancelled}

func States() []State {
	return slices.Clone(states)
}

func (s State) Valid() bool {
	return slices.Contains(states, s)
}

func (s State) Terminal() bool {
	return slices.Contains(terminalStates, s)
}
