package surface

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/socrates-devtool/socrates/internal/document"
)

type Sent struct {
	Decisions []Choice `json:"decisions"`
	Message   string   `json:"message"`
	Ends      bool     `json:"ends"`
}

type Choice struct {
	ID     string `json:"id"`
	Chosen []int  `json:"chosen"`
	Words  string `json:"words"`
}

type Answered struct {
	Decisions []Decided `json:"decisions"`
	Message   string    `json:"message"`
	Ended     bool      `json:"ended"`
}

type Decided struct {
	ID       string   `json:"id"`
	Chosen   []string `json:"chosen"`
	Commands []string `json:"commands"`
	Words    string   `json:"words"`
}

type Exchange struct {
	mu       sync.Mutex
	read     document.Document
	revision int
	pending  []Answered
	over     bool
	changed  chan struct{}
	readers  int
	seen     time.Time
}

func NewExchange(read document.Document, at time.Time) *Exchange {
	return &Exchange{
		read:     read,
		revision: 1,
		changed:  make(chan struct{}),
		seen:     at,
	}
}

func (e *Exchange) notify() {
	close(e.changed)
	e.changed = make(chan struct{})
}

func (e *Exchange) Document() (document.Document, int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.read, e.revision
}

func (e *Exchange) Revise(read document.Document) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.read = read
	e.revision++
	e.over = false
	e.notify()
	return e.revision
}

func (e *Exchange) Over() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.over
}

func (e *Exchange) End() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.over {
		return
	}
	e.over = true
	e.notify()
}

func (e *Exchange) Receive(sent Sent) (Answered, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.over {
		return Answered{}, fmt.Errorf(exchangeIsOver)
	}

	answered, err := resolve(e.read, sent)
	if err != nil {
		return Answered{}, err
	}

	e.pending = append(e.pending, answered)
	if sent.Ends {
		e.over = true
	}
	e.notify()
	return answered, nil
}

func (e *Exchange) Await(ctx context.Context, within time.Duration) (Answered, bool) {
	deadline := time.After(within)

	for {
		e.mu.Lock()
		if len(e.pending) > 0 {
			taken := e.pending[0]
			e.pending = e.pending[1:]
			e.mu.Unlock()
			return taken, true
		}
		waiting := e.changed
		e.mu.Unlock()

		select {
		case <-ctx.Done():
			return Answered{}, false
		case <-deadline:
			return Answered{}, false
		case <-waiting:
		}
	}
}

func (e *Exchange) Changed() <-chan struct{} {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.changed
}

func (e *Exchange) Touch(at time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.seen = at
}

func (e *Exchange) Reading(delta int, at time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.readers += delta
	e.seen = at
}

func (e *Exchange) IdleSince(at time.Time, period time.Duration) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.readers == 0 && at.Sub(e.seen) >= period
}

func resolve(read document.Document, sent Sent) (Answered, error) {
	answered := Answered{Message: sent.Message, Ended: sent.Ends}

	for _, choice := range sent.Decisions {
		position := slices.IndexFunc(read.Decisions, func(held document.Decision) bool {
			return held.ID == choice.ID
		})
		if position < 0 {
			return Answered{}, fmt.Errorf(answerNamesNoDecision, choice.ID)
		}

		decision := read.Decisions[position]
		if decision.Mode == document.One && len(choice.Chosen) > 1 {
			return Answered{}, fmt.Errorf(singleChoiceCarriesMoreThanOne, choice.ID, len(choice.Chosen))
		}

		decided := Decided{ID: decision.ID, Words: choice.Words}
		for _, offered := range slices.Sorted(slices.Values(choice.Chosen)) {
			if offered < 0 || offered >= len(decision.Options) {
				return Answered{}, fmt.Errorf(answerNamesNoOption, offered, choice.ID)
			}
			decided.Chosen = append(decided.Chosen, decision.Options[offered].Label)
			decided.Commands = append(decided.Commands, decision.Options[offered].Command)
		}
		answered.Decisions = append(answered.Decisions, decided)
	}

	return answered, nil
}
