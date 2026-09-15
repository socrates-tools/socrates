package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/socrates-tools/socrates/internal/process"
	"github.com/socrates-tools/socrates/internal/runlog"
)

const (
	opencodeCommand      = "run"
	opencodeEventStream  = "--format"
	opencodeJSON         = "json"
	opencodeDirectory    = "--dir"
	opencodeModel        = "--model"
	opencodeVariant      = "--variant"
	opencodeNoProjectDoc = "OPENCODE_DISABLE_PROJECT_CONFIG=1"

	opencodeText       = "text"
	opencodeStepFinish = "step_finish"
	opencodeErrorEvent = "error"
)

type opencodeAgent struct {
	vendor   Vendor
	binary   string
	routing  map[Role]Routing
	deadline time.Duration
	recorder runlog.Recorder
}

func newOpencode(vendor Vendor, options Options) (Agent, error) {
	routing, err := resolveRouting(vendor, options.Routing)
	if err != nil {
		return nil, err
	}

	agent := &opencodeAgent{
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

func (a *opencodeAgent) Run(ctx context.Context, request Request) (Response, error) {
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

func (a *opencodeAgent) attempt(bounded, parent context.Context, request Request, routing Routing) (Response, error) {
	command := exec.CommandContext(bounded, a.binary, a.arguments(routing, request)...)
	command.Dir = request.Dir
	command.Env = append(os.Environ(), opencodeNoProjectDoc)
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

	read, parseErr := parseOpencodeStream(reader)
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
		if read.account != "" {
			return Response{}, TerminalReportedSaying(a.binary, AccountOf(read.account, failures.String()), opencodeErrorEvent)
		}
		return Response{}, NoTerminalEvent(a.binary)
	}

	answer := Response{
		Text:      read.text,
		Model:     routing.Model,
		Effort:    routing.Effort,
		SessionID: read.sessionID,
		Usage:     read.usage,
	}

	document, err := Document(a.vendor, request, a.binary, nil, read.text)
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

func (a *opencodeAgent) arguments(routing Routing, request Request) []string {
	arguments := []string{
		opencodeCommand,
		opencodeEventStream, opencodeJSON,
		opencodeDirectory, request.Dir,
	}
	if routing.Model != "" {
		arguments = append(arguments, opencodeModel, routing.Model)
	}
	if routing.Effort != "" {
		arguments = append(arguments, opencodeVariant, routing.Effort)
	}
	return arguments
}

type opencodeTokens struct {
	Input     int `json:"input"`
	Output    int `json:"output"`
	Reasoning int `json:"reasoning"`
	Cache     struct {
		Write int `json:"write"`
		Read  int `json:"read"`
	} `json:"cache"`
}

type opencodeEvent struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionID"`
	Part      *struct {
		Type   string          `json:"type"`
		Text   string          `json:"text"`
		Reason string          `json:"reason"`
		Tokens *opencodeTokens `json:"tokens"`
	} `json:"part"`
	Error *struct {
		Name string `json:"name"`
		Data struct {
			Message string `json:"message"`
		} `json:"data"`
	} `json:"error"`
}

type opencodeStreamed struct {
	sessionID string
	text      string
	usage     Usage
	completed bool
	account   string
}

func parseOpencodeStream(input io.Reader) (opencodeStreamed, error) {
	watched := &lastByteReader{inner: input}
	scanner := bufio.NewScanner(watched)
	scanner.Buffer(make([]byte, 0, streamBufferSize), streamLineCeiling)

	var read opencodeStreamed
	malformed := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var event opencodeEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			malformed = line
			continue
		}
		if malformed != "" {
			return read, notAnEvent(malformed)
		}

		if event.SessionID != "" {
			read.sessionID = event.SessionID
		}

		switch event.Type {
		case opencodeText:
			if event.Part != nil {
				read.text += event.Part.Text
			}
		case opencodeStepFinish:
			read.completed = true
			if event.Part != nil && event.Part.Tokens != nil {
				read.usage = Usage{
					InputTokens:         event.Part.Tokens.Input,
					OutputTokens:        event.Part.Tokens.Output,
					CacheReadTokens:     event.Part.Tokens.Cache.Read,
					CacheCreationTokens: event.Part.Tokens.Cache.Write,
				}
			}
		case opencodeErrorEvent:
			if event.Error != nil {
				read.account = AccountOf(read.account, Bounded(strings.TrimPrefix(event.Error.Name+": "+event.Error.Data.Message, ": ")))
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
