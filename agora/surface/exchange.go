package surface

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/internal/refusal"
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
	mu          sync.Mutex
	read        document.Document
	revision    int
	uncollected *Answered
	answeredAt  int
	collected   map[string]document.Settled
	over        bool
	changed     chan struct{}
	readers     int
	seen        time.Time
}

func NewExchange(read document.Document, at time.Time) *Exchange {
	return &Exchange{
		read:      read,
		revision:  1,
		collected: map[string]document.Settled{},
		changed:   make(chan struct{}),
		seen:      at,
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

	return e.reviseLocked(read)
}

func (e *Exchange) ReviseAdmitted(read document.Document) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.admitsLocked(read); err != nil {
		return 0, err
	}
	return e.reviseLocked(read), nil
}

func (e *Exchange) Admits(read document.Document) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.admitsLocked(read)
}

func (e *Exchange) admitsLocked(read document.Document) error {
	for _, stated := range read.Settled {
		answered, collected := e.collected[stated.ID]
		if !collected {
			continue
		}
		if stated.Question != answered.Question {
			return refusal.New(settledOtherwiseThanCollected, stated.ID, "question", answered.Question)
		}
		if !slices.Equal(stated.Chosen, answered.Chosen) {
			return refusal.New(settledOtherwiseThanCollected, stated.ID, "chosen", strings.Join(answered.Chosen, ", "))
		}
		if stated.Words != answered.Words {
			return refusal.New(settledOtherwiseThanCollected, stated.ID, "words", answered.Words)
		}
	}
	return nil
}

func (e *Exchange) reviseLocked(read document.Document) int {
	e.read = read
	e.revision++
	e.over = false
	e.notify()
	return e.revision
}

func (e *Exchange) Standing() (int, bool, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.revision, e.over, e.uncollectedLocked()
}

func (e *Exchange) uncollectedLocked() bool {
	return e.uncollected != nil && e.answeredAt == e.revision
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
		return Answered{}, errors.New(exchangeIsOver)
	}

	answered, err := resolve(e.read, sent)
	if err != nil {
		return Answered{}, err
	}

	e.uncollected = &answered
	e.answeredAt = e.revision
	if sent.Ends {
		e.over = true
	}
	e.notify()
	return answered, nil
}

func (e *Exchange) Await(ctx context.Context, within time.Duration, keep func(Answered) error) (Answered, bool, error) {
	deadline := time.After(within)

	for {
		e.mu.Lock()
		if e.uncollectedLocked() {
			taken := *e.uncollected
			if err := keep(taken); err != nil {
				e.mu.Unlock()
				return Answered{}, false, err
			}
			e.uncollected = nil
			e.holdCollectedLocked(taken)
			e.notify()
			e.mu.Unlock()
			return taken, true, nil
		}
		waiting := e.changed
		e.mu.Unlock()

		select {
		case <-ctx.Done():
			return Answered{}, false, nil
		case <-deadline:
			return Answered{}, false, nil
		case <-waiting:
		}
	}
}

func (e *Exchange) holdCollectedLocked(taken Answered) {
	for _, decided := range taken.Decisions {
		position := slices.IndexFunc(e.read.Decisions, func(held document.Decision) bool {
			return held.ID == decided.ID
		})
		if position < 0 {
			continue
		}
		e.collected[decided.ID] = document.Settled{
			ID:       decided.ID,
			Question: e.read.Decisions[position].Question,
			Chosen:   decided.Chosen,
			Words:    decided.Words,
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
			option := decision.Options[offered]
			decided.Chosen = append(decided.Chosen, option.Label)
			if strings.TrimSpace(option.Command) != "" {
				decided.Commands = append(decided.Commands, option.Command)
			}
		}
		answered.Decisions = append(answered.Decisions, decided)
	}

	return answered, nil
}
