package agent

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/socrates-tools/socrates/internal/process"
	"github.com/socrates-tools/socrates/socrates-review/runlog"
)

const (
	DefaultDeadline = 30 * time.Minute

	inspectingTools  = "Read,Grep,Glob,WebSearch,WebFetch"
	schemaTool       = "StructuredOutput"
	noCustomisations = "--safe-mode"
	settingSources   = "user"
	streamFormat     = "stream-json"
	waitDelay        = 5 * time.Second
)

const (
	EventStarted           = "agent-started"
	EventCompleted         = "agent-completed"
	EventFailed            = "agent-failed"
	EventUnexpectedTools   = "agent-unexpected-tools"
	EventUnexpectedSources = "agent-unexpected-sources"
	EventToolsUnreported   = "agent-tools-unreported"
)

type Options struct {
	Binary          string
	Routing         Configuration
	Deadline        time.Duration
	Recorder        runlog.Recorder
	SchemaDirectory string
}

type claudeAgent struct {
	vendor   Vendor
	binary   string
	routing  map[Role]Routing
	deadline time.Duration
	recorder runlog.Recorder
}

func NewClaude(options Options) (Agent, error) {
	return newClaude(TheDefault(), options)
}

func newClaude(vendor Vendor, options Options) (Agent, error) {
	routing, err := resolveRouting(vendor, options.Routing)
	if err != nil {
		return nil, err
	}

	agent := &claudeAgent{
		vendor:   vendor,
		binary:   options.Binary,
		routing:  routing,
		deadline: options.Deadline,
		recorder: options.Recorder,
	}
	if agent.binary == "" {
		agent.binary = vendor.Binary
	}
	if agent.deadline <= 0 {
		agent.deadline = DefaultDeadline
	}
	return agent, nil
}

func (a *claudeAgent) Run(ctx context.Context, request Request) (Response, error) {
	routing, known := a.routing[request.Role]
	if !known {
		return Response{}, fmt.Errorf(notARoleAtTheCall, request.Role)
	}

	bounded, release := context.WithTimeout(ctx, a.deadline)
	defer release()

	answer, err := a.attempt(bounded, ctx, request, routing)

	a.recordToolSurface(request, answer, err == nil)
	if err != nil {
		return answer, a.failed(request, routing, err)
	}
	return answer, nil
}

func (a *claudeAgent) recordToolSurface(request Request, answer Response, completed bool) {
	if len(answer.ToolSources) > 0 {
		a.record(runlog.Event{
			Kind:    EventUnexpectedSources,
			Subject: string(request.Role),
			Fields: map[string]any{
				"sources": strings.Join(answer.ToolSources, ", "),
			},
		})
	}

	if len(answer.Tools) == 0 {
		if completed {
			a.record(runlog.Event{
				Kind:    EventToolsUnreported,
				Subject: string(request.Role),
				Fields: map[string]any{
					"asked_for": strings.Join(ToolsAskedFor(), ", "),
				},
			})
		}
		return
	}

	unexpected := UnexpectedTools(answer.Tools)
	if len(unexpected) == 0 {
		return
	}

	a.record(runlog.Event{
		Kind:    EventUnexpectedTools,
		Subject: string(request.Role),
		Fields: map[string]any{
			"tools": strings.Join(unexpected, ", "),
		},
	})
}

func ToolsAskedFor() []string {
	return append(strings.Split(inspectingTools, ","), schemaTool)
}

func UnexpectedTools(held []string) []string {
	asked := ToolsAskedFor()

	var unexpected []string
	for _, tool := range held {
		if !slices.Contains(asked, tool) {
			unexpected = append(unexpected, tool)
		}
	}
	return unexpected
}

func (a *claudeAgent) attempt(bounded, parent context.Context, request Request, routing Routing) (Response, error) {

	arguments := a.arguments(routing, request)
	command := exec.CommandContext(bounded, a.binary, arguments...)
	command.Dir = request.Dir
	command.Stdin = strings.NewReader(request.Prompt)
	process.OwnGroup(command)
	command.Cancel = func() error { return process.KillGroup(command.Process.Pid) }
	command.WaitDelay = waitDelay

	failures := NewBoundedText()
	command.Stderr = failures

	out, err := command.StdoutPipe()
	if err != nil {
		return Response{}, fmt.Errorf(readingWhatItWrites, a.binary, err)
	}

	if err := command.Start(); err != nil {
		return Response{}, fmt.Errorf(startingTheAgent, a.binary, err)
	}
	defer func() {
		_ = process.KillGroup(command.Process.Pid)
	}()
	a.record(runlog.Event{
		Kind:    EventStarted,
		Subject: string(request.Role),
		Fields: map[string]any{
			"model":   routing.Model,
			"effort":  routing.Effort,
			"process": command.Process.Pid,
		},
	})

	var reader io.Reader = out
	if request.Transcript != nil {
		reader = io.TeeReader(out, request.Transcript)
	}

	read, parseErr := parseStream(reader)
	if parseErr != nil {
		_, _ = io.Copy(io.Discard, out)
	}
	waitErr := command.Wait()

	reached := Response{Tools: read.tools, ToolSources: read.toolSources}

	if deadlineErr := deadlineReached(a.binary, a.deadline, bounded, parent); deadlineErr != nil {
		return reached, deadlineErr
	}
	if parseErr != nil {
		return reached, parseErr
	}
	if waitErr != nil {
		return reached, ExitedSaying(a.binary, waitErr, AccountOf(read.accountGiven(), failures.String()))
	}

	answer := Response{
		Text:        read.text,
		Model:       read.model,
		Effort:      routing.Effort,
		SessionID:   read.sessionID,
		Usage:       read.usage,
		Tools:       read.tools,
		ToolSources: read.toolSources,
	}

	if read.result == nil {
		return reached, NoTerminalEvent(a.binary)
	}
	if read.result.reportsAnError() {
		return reached, TerminalReportedSaying(a.binary, read.accountGiven(), read.result.Subtype)
	}

	document, err := Document(a.vendor, request, a.binary, read.result.StructuredOutput, read.text)
	if err != nil {
		return reached, err
	}
	answer.Structured = document

	if err := Honoured(request, answer.Structured); err != nil {
		return answer, err
	}

	a.record(runlog.Event{
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

	return answer, nil
}

func (a *claudeAgent) arguments(routing Routing, request Request) []string {
	arguments := []string{
		"-p",
		"--verbose",
		"--output-format", streamFormat,
		noCustomisations,
		"--setting-sources", settingSources,
		"--tools", inspectingTools,
		"--dangerously-skip-permissions",
	}
	if routing.Model != "" {
		arguments = append(arguments, "--model", routing.Model)
	}
	if routing.Effort != "" {
		arguments = append(arguments, "--effort", routing.Effort)
	}
	if len(request.Schema) > 0 {
		arguments = append(arguments, "--json-schema", string(request.Schema))
	}
	return arguments
}

func (a *claudeAgent) failed(request Request, routing Routing, reason error) error {
	a.record(runlog.Event{
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

func (a *claudeAgent) record(event runlog.Event) {
	if a.recorder == nil {
		return
	}
	_ = a.recorder.Record(event)
}
