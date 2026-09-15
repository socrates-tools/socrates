package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/socrates-tools/socrates/socrates-review/findings"
	"github.com/socrates-tools/socrates/socrates-review/runlog"
)

const failureTailBytes = 8 * 1024

func Document(vendor Vendor, request Request, binary string, structured json.RawMessage, text string) (json.RawMessage, error) {
	if len(request.Schema) == 0 {
		return structured, nil
	}

	if vendor.Tier == Parsed {
		extracted, err := findings.Extract(text)
		if err != nil {
			return nil, fmt.Errorf(noDocumentInTheAnswer, binary, err)
		}
		return extracted, nil
	}

	if len(structured) == 0 {
		return nil, fmt.Errorf(noStructuredOutput, binary)
	}
	if !json.Valid(structured) {
		return nil, fmt.Errorf(structuredOutputIsNotJSON, binary)
	}
	return structured, nil
}

func Honoured(request Request, document json.RawMessage) error {
	if request.Honours == nil || len(document) == 0 {
		return nil
	}
	if err := request.Honours(document); err != nil {
		return fmt.Errorf(schemaNotSatisfied, err)
	}
	return nil
}

func AccountOf(given ...string) string {
	var said []string
	for _, account := range given {
		if trimmed := strings.TrimSpace(account); trimmed != "" {
			said = append(said, trimmed)
		}
	}
	return strings.Join(said, beside)
}

func Bounded(said string) string {
	bounded := NewBoundedText()
	_, _ = io.WriteString(bounded, said)
	return bounded.String()
}

func NewBoundedText() *BoundedText {
	return &BoundedText{limit: failureTailBytes}
}

func deadlineReached(binary string, deadline time.Duration, bounded, parent context.Context) error {
	if !errors.Is(bounded.Err(), context.DeadlineExceeded) {
		return nil
	}
	if errors.Is(parent.Err(), context.DeadlineExceeded) {
		return errors.New(theDeadlinePassed)
	}
	return fmt.Errorf(silenceUntilTheDeadline, binary, deadline)
}

func ExitedSaying(binary string, waited error, said string) error {
	if said == "" {
		return fmt.Errorf(theAgentExitedSayingNothing, binary, waited)
	}
	return fmt.Errorf(theAgentExited, binary, waited, said)
}

func NoTerminalEvent(binary string) error {
	return fmt.Errorf(noResultEvent, binary)
}

func TerminalReportedSaying(binary, said, otherwise string) error {
	if said != "" {
		return fmt.Errorf(theResultReportedSaying, binary, said)
	}
	return fmt.Errorf(theResultReported, binary, otherwise)
}

type BoundedText struct {
	limit   int
	written []byte
	elided  bool
}

func (h *BoundedText) Write(bytes []byte) (int, error) {
	room := h.limit - len(h.written)
	if room <= 0 {
		h.elided = true
		return len(bytes), nil
	}
	if len(bytes) > room {
		h.written = append(h.written, bytes[:room]...)
		h.elided = true
		return len(bytes), nil
	}
	h.written = append(h.written, bytes...)
	return len(bytes), nil
}

func (h *BoundedText) String() string {
	written := strings.TrimSpace(string(h.written))
	if h.elided {
		return written + "..."
	}
	return written
}

func recordStarted(recorder runlog.Recorder, request Request, routing Routing, process int) {
	recordOn(recorder, runlog.Event{
		Kind:    EventStarted,
		Subject: string(request.Role),
		Fields: map[string]any{
			"model":   routing.Model,
			"effort":  routing.Effort,
			"process": process,
		},
	})
}

func recordCompleted(recorder runlog.Recorder, request Request, routing Routing, answer Response) {
	recordOn(recorder, runlog.Event{
		Kind:    EventCompleted,
		Subject: string(request.Role),
		Fields: map[string]any{
			"model":                 answer.Model,
			"effort":                routing.Effort,
			"input_tokens":          answer.Usage.InputTokens,
			"output_tokens":         answer.Usage.OutputTokens,
			"cache_read_tokens":     answer.Usage.CacheReadTokens,
			"cache_creation_tokens": answer.Usage.CacheCreationTokens,
		},
	})
}

func recordFailure(recorder runlog.Recorder, request Request, routing Routing, reason error) error {
	recordOn(recorder, runlog.Event{
		Kind:    EventFailed,
		Subject: string(request.Role),
		Fields: map[string]any{
			"model":  routing.Model,
			"effort": routing.Effort,
			"reason": reason.Error(),
		},
	})
	return reason
}

func recordSessionHeldNothingStated(vendor Vendor, recorder runlog.Recorder, request Request, completed bool) {
	if !completed || vendor.ReportsItsToolInventory {
		return
	}
	recordOn(recorder, runlog.Event{
		Kind:    EventToolsUnreported,
		Subject: string(request.Role),
		Fields: map[string]any{
			"asked_for": strings.Join(ToolsAskedFor(), ", "),
		},
	})
}

func recordOn(recorder runlog.Recorder, event runlog.Event) {
	if recorder == nil {
		return
	}
	_ = recorder.Record(event)
}
