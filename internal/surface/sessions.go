package surface

import (
	"slices"
	"sync"
	"time"

	"github.com/rui-ktei/socrates/internal/document"
)

type held struct {
	token    string
	exchange *Exchange
}

type Sessions struct {
	mu       sync.Mutex
	standing map[string]*held
	emptied  time.Time
}

func NewSessions(at time.Time) *Sessions {
	return &Sessions{standing: map[string]*held{}, emptied: at}
}

func (s *Sessions) Open(name, token string, read document.Document, at time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.standing[name] = &held{token: token, exchange: NewExchange(read, at)}
	s.emptied = time.Time{}
}

func (s *Sessions) Exchange(name string) (*Exchange, bool) {
	one, standing := s.find(name)
	if !standing {
		return nil, false
	}
	return one.exchange, true
}

func (s *Sessions) Token(name string) (string, bool) {
	one, standing := s.find(name)
	if !standing {
		return "", false
	}
	return one.token, true
}

func (s *Sessions) Names() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	names := make([]string, 0, len(s.standing))
	for name := range s.standing {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func (s *Sessions) Collect(now time.Time, period time.Duration) ([]string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var collected []string
	for name, one := range s.standing {
		if one.exchange.IdleSince(now, period) {
			delete(s.standing, name)
			collected = append(collected, name)
		}
	}
	slices.Sort(collected)

	if len(s.standing) > 0 {
		s.emptied = time.Time{}
		return collected, false
	}
	if s.emptied.IsZero() {
		s.emptied = now
		return collected, false
	}
	return collected, now.Sub(s.emptied) >= period
}

func (s *Sessions) find(name string) (*held, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	one, standing := s.standing[name]
	return one, standing
}
