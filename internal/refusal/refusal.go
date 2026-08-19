package refusal

import (
	"errors"
	"fmt"
)

type Refusal struct {
	Reason string
	Cause  error
}

func New(format string, arguments ...any) *Refusal {
	return &Refusal{Reason: fmt.Sprintf(format, arguments...)}
}

func Wrapping(cause error, format string, arguments ...any) *Refusal {
	return &Refusal{Reason: fmt.Sprintf(format, arguments...), Cause: cause}
}

func (r *Refusal) Error() string {
	return r.Reason
}

func (r *Refusal) Unwrap() error {
	return r.Cause
}

func Is(err error) bool {
	var refused *Refusal
	return errors.As(err, &refused)
}
