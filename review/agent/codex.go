package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/socrates-tools/socrates/internal/process"
	"github.com/socrates-tools/socrates/review/runlog"
)

const (
	codexCommand         = "exec"
	codexEventStream     = "--json"
	codexSandbox         = "--sandbox"
	codexReadOnly        = "read-only"
	codexSchemaFile      = "--output-schema"
	codexNoUserConfig    = "--ignore-user-config"
	codexOverride        = "-c"
	codexNoProjectDoc    = "project_doc_max_bytes=0"
	codexModel           = "-m"
	codexEffort          = "model_reasoning_effort="
	codexPromptFromStdin = "-"

	codexSchemaName = "schema.json"

	codexThreadStarted = "thread.started"
	codexTurnCompleted = "turn.completed"
	codexTurnFailed    = "turn.failed"
	codexItemCompleted = "item.completed"
	codexErrorEvent    = "error"
	codexAgentMessage  = "agent_message"
)

type codexAgent struct {
	vendor   Vendor
	binary   string
	routing  map[Role]Routing
	deadline time.Duration
	recorder runlog.Recorder
	schemaIn string
}

func newCodex(vendor Vendor, options Options) (Agent, error) {
	routing, err := resolveRouting(vendor, options.Routing)
	if err != nil {
		return nil, err
	}

	agent := &codexAgent{
		vendor:   vendor,
		binary:   options.Binary,
		routing:  routing,
		deadline: options.Deadline,
		recorder: options.Recorder,
		schemaIn: options.SchemaDirectory,
	}
	if agent.binary == "" {
		agent.binary = vendor.Binary
	}
	if agent.deadline <= 0 {
		agent.deadline = DefaultDeadline
	}
	return agent, nil
}

func (a *codexAgent) Run(ctx context.Context, request Request) (Response, error) {
	routing, known := a.routing[request.Role]
	if !known {
		return Response{}, fmt.Errorf(notARoleAtTheCall, request.Role)
	}

	bounded, release := context.WithTimeout(ctx, a.deadline)
	defer release()

	answer, err := a.attempt(bounded, ctx, request, routing)

	recordSessionHeldNothingStated(a.vendor, a.recorder, request, err == nil)
	if err != nil {
		return answer, recordFailure(a.recorder, request, routing, err)
	}
	return answer, nil
}

func (a *codexAgent) attempt(bounded, parent context.Context, request Request, routing Routing) (Response, error) {
	arguments, written, err := a.arguments(routing, request)
	if err != nil {
		return Response{}, err
	}
	if written != "" {
		defer func() { _ = os.Remove(written) }()
	}

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
	recordStarted(a.recorder, request, routing, command.Process.Pid)

	var reader io.Reader = out
	if request.Transcript != nil {
		reader = io.TeeReader(out, request.Transcript)
	}

	read, parseErr := parseCodexStream(reader)
	if parseErr != nil {
		_, _ = io.Copy(io.Discard, out)
	}
	waitErr := command.Wait()

	if deadlineErr := deadlineReached(a.binary, a.deadline, bounded, parent); deadlineErr != nil {
		return Response{}, deadlineErr
	}
	if parseErr != nil {
		return Response{}, parseErr
	}
	if waitErr != nil {
		return Response{}, ExitedSaying(a.binary, waitErr, AccountOf(read.account, failures.String()))
	}

	if !read.completed {
		if read.account != "" || read.failed {
			return Response{}, TerminalReportedSaying(a.binary, AccountOf(read.account, failures.String()), codexTurnFailed)
		}
		return Response{}, NoTerminalEvent(a.binary)
	}

	answer := Response{
		Text:      read.text,
		Model:     routing.Model,
		Effort:    routing.Effort,
		SessionID: read.threadID,
		Usage:     read.usage,
	}

	document, err := Document(a.vendor, request, a.binary, json.RawMessage(strings.TrimSpace(read.text)), read.text)
	if err != nil {
		return Response{}, err
	}
	answer.Structured = document

	if err := Honoured(request, answer.Structured); err != nil {
		return answer, err
	}

	recordCompleted(a.recorder, request, routing, answer)
	return answer, nil
}

func (a *codexAgent) arguments(routing Routing, request Request) ([]string, string, error) {
	arguments := []string{
		codexCommand,
		codexEventStream,
		codexSandbox, codexReadOnly,
		codexNoUserConfig,
		codexOverride, codexNoProjectDoc,
	}

	written := ""
	if len(request.Schema) > 0 {
		path, err := a.writeSchema(request.Scratch, request.Schema)
		if err != nil {
			return nil, "", err
		}
		written = path
		arguments = append(arguments, codexSchemaFile, path)
	}

	if routing.Model != "" {
		arguments = append(arguments, codexModel, routing.Model)
	}
	if routing.Effort != "" {
		arguments = append(arguments, codexOverride, codexEffort+routing.Effort)
	}

	return append(arguments, codexPromptFromStdin), written, nil
}

func (a *codexAgent) writeSchema(scratch string, schema json.RawMessage) (string, error) {
	directory := scratch
	if directory == "" {
		directory = a.schemaIn
	}
	if directory == "" {
		made, err := os.MkdirTemp("", "socrates-schema-")
		if err != nil {
			return "", fmt.Errorf(writingTheSchema, err)
		}
		directory = made
	}
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", fmt.Errorf(writingTheSchema, err)
	}

	path := filepath.Join(directory, codexSchemaName)
	if err := os.WriteFile(path, schema, 0o644); err != nil {
		return "", fmt.Errorf(writingTheSchema, err)
	}
	return path, nil
}

type codexUsage struct {
	InputTokens       int `json:"input_tokens"`
	CachedInputTokens int `json:"cached_input_tokens"`
	CacheWriteTokens  int `json:"cache_write_input_tokens"`
	OutputTokens      int `json:"output_tokens"`
}

type codexItem struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Message string `json:"message"`
}

type codexEvent struct {
	Type     string      `json:"type"`
	ThreadID string      `json:"thread_id"`
	Message  string      `json:"message"`
	Item     *codexItem  `json:"item"`
	Usage    *codexUsage `json:"usage"`
	Error    *struct {
		Message string `json:"message"`
	} `json:"error"`
}

type codexStreamed struct {
	threadID  string
	text      string
	usage     Usage
	completed bool
	failed    bool
	account   string
}

func parseCodexStream(input io.Reader) (codexStreamed, error) {
	watched := &lastByteReader{inner: input}
	scanner := bufio.NewScanner(watched)
	scanner.Buffer(make([]byte, 0, streamBufferSize), streamLineCeiling)

	var read codexStreamed
	malformed := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var event codexEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			malformed = line
			continue
		}
		if malformed != "" {
			return read, notAnEvent(malformed)
		}

		switch event.Type {
		case codexThreadStarted:
			read.threadID = event.ThreadID
		case codexItemCompleted:
			if event.Item != nil && event.Item.Type == codexAgentMessage {
				read.text += event.Item.Text
			}
		case codexErrorEvent:
			read.account = AccountOf(read.account, Bounded(event.Message))
		case codexTurnCompleted:
			read.completed = true
			if event.Usage != nil {
				read.usage = Usage{
					InputTokens:         event.Usage.InputTokens,
					OutputTokens:        event.Usage.OutputTokens,
					CacheReadTokens:     event.Usage.CachedInputTokens,
					CacheCreationTokens: event.Usage.CacheWriteTokens,
				}
			}
		case codexTurnFailed:
			read.failed = true
			if event.Error != nil {
				read.account = AccountOf(read.account, Bounded(event.Error.Message))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return read, fmt.Errorf(readingTheStream, err)
	}
	if malformed != "" {
		if watched.truncated() {
			return read, truncatedStream(malformed)
		}
		return read, notAnEvent(malformed)
	}
	if watched.truncated() {
		return read, truncatedStream("")
	}

	return read, nil
}
