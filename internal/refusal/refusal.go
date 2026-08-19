package refusal

import (
	"errors"
	"fmt"
)

type Refusal struct {
	Reason string
}

func New(format string, arguments ...any) *Refusal {
	return &Refusal{Reason: fmt.Sprintf(format, arguments...)}
}

func (r *Refusal) Error() string {
	return r.Reason
}

func Is(err error) bool {
	var refused *Refusal
	return errors.As(err, &refused)
}
