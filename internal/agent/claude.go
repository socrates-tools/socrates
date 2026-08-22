package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/rui-ktei/socrates/internal/process"
	"github.com/rui-ktei/socrates/internal/runlog"
)

const (
	DefaultBinary   = "claude"
	DefaultDeadline = 30 * time.Minute

	readingTools         = "Read,Grep,Glob"
	toolServerPrefix     = "mcp__"
	noToolServersAllowed = `{"allowedMcpServers":[]}`
	settingSources       = "user"
	streamFormat         = "stream-json"
	waitDelay            = 5 * time.Second
	failureTailBytes     = 8 * 1024
)

const (
	EventStarted           = "agent-started"
	EventCompleted         = "agent-completed"
	EventFailed            = "agent-failed"
	EventToolServersLoaded = "agent-tool-servers-loaded"
)

type Options struct {
	Binary   string
	Routing  Configuration
	Deadline time.Duration
	Recorder runlog.Recorder
}

type claudeAgent struct {
	binary   string
	routing  map[Role]Routing
	deadline time.Duration
	recorder runlog.Recorder
}

func NewClaude(options Options) (Agent, error) {
	routing, err := resolveRouting(options.Routing)
	if err != nil {
		return nil, err
	}

	agent := &claudeAgent{
		binary:   options.Binary,
		routing:  routing,
		deadline: options.Deadline,
		recorder: options.Recorder,
	}
	if agent.binary == "" {
		agent.binary = DefaultBinary
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

	a.recordToolServers(request, answer)
	if err != nil {
		return answer, a.failed(request, routing, err)
	}
	return answer, nil
}

func (a *claudeAgent) recordToolServers(request Request, answer Response) {
	contributed := ToolsFromServers(answer.Tools)
	if len(answer.ToolServers) == 0 && len(contributed) == 0 {
		return
	}

	named := answer.ToolServers
	if len(named) == 0 {
		named = []string{whateverNamedThem}
	}
	a.record(runlog.Event{
		Kind:    EventToolServersLoaded,
		Subject: string(request.Role),
		Fields: map[string]any{
			"servers": strings.Join(named, ", "),
			"tools":   len(contributed),
		},
	})
}

func ToolsFromServers(tools []string) []string {
	var contributed []string
	for _, tool := range tools {
		if strings.HasPrefix(tool, toolServerPrefix) {
			contributed = append(contributed, tool)
		}
	}
	return contributed
}

func (a *claudeAgent) attempt(bounded, parent context.Context, request Request, routing Routing) (Response, error) {

	arguments := a.arguments(routing, request)
	command := exec.CommandContext(bounded, a.binary, arguments...)
	command.Dir = request.Dir
	command.Stdin = strings.NewReader(request.Prompt)
	process.OwnGroup(command)
	command.Cancel = func() error { return process.KillGroup(command.Process.Pid) }
	command.WaitDelay = waitDelay

	failures := &head{limit: failureTailBytes}
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

	reached := Response{Opened: read.opened, ToolServers: read.toolServers, Tools: read.tools}

	if deadlineErr := a.deadlineReached(bounded, parent); deadlineErr != nil {
		return reached, deadlineErr
	}
	if parseErr != nil {
		return reached, parseErr
	}
	if waitErr != nil {
		return reached, fmt.Errorf(theAgentExited, a.binary, waitErr, failures.String())
	}

	answer := Response{
		Text:        read.text,
		Model:       read.model,
		Effort:      routing.Effort,
		SessionID:   read.sessionID,
		Usage:       read.usage,
		Opened:      read.opened,
		ToolServers: read.toolServers,
		Tools:       read.tools,
	}

	if read.result == nil {
		return reached, fmt.Errorf(noResultEvent, a.binary)
	}
	if read.result.IsError || read.result.Subtype != "success" {
		return reached, fmt.Errorf(theResultReported, a.binary, read.result.Subtype)
	}

	if len(request.Schema) > 0 {
		if len(read.result.StructuredOutput) == 0 {
			return reached, fmt.Errorf(noStructuredOutput, a.binary)
		}
		if !json.Valid(read.result.StructuredOutput) {
			return reached, fmt.Errorf(structuredOutputIsNotJSON, a.binary)
		}
	}
	answer.Structured = read.result.StructuredOutput

	if request.Honours != nil && len(answer.Structured) > 0 {
		if err := request.Honours(answer.Structured); err != nil {
			return answer, fmt.Errorf(schemaNotSatisfied, err)
		}
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
		"--setting-sources", settingSources,
		"--settings", noToolServersAllowed,
		"--tools", readingTools,
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

func (a *claudeAgent) deadlineReached(bounded, parent context.Context) error {
	if !errors.Is(bounded.Err(), context.DeadlineExceeded) {
		return nil
	}
	if errors.Is(parent.Err(), context.DeadlineExceeded) {
		return errors.New(theDeadlinePassed)
	}
	return fmt.Errorf(silenceUntilTheDeadline, a.binary, a.deadline)
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

type head struct {
	limit   int
	written []byte
	elided  bool
}

func (h *head) Write(bytes []byte) (int, error) {
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

func (h *head) String() string {
	written := strings.TrimSpace(string(h.written))
	if h.elided {
		return written + "..."
	}
	return written
}
