package codexsim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ktei/socrates/internal/simulation"
)

const (
	defaultRejection = "error: unexpected argument '--settings' found"

	promptFromStdin  = "-"
	defaultThreadID  = "01a06ab9-0000-7000-8000-000000000000"
	sandboxReadOnly  = "read-only"
	effortOverride   = "model_reasoning_effort="
	projectDocCapped = "project_doc_max_bytes=0"
)

type wire struct{}

func Run(arguments []string, input io.Reader, out, failures io.Writer, environment []string) int {
	return simulation.Run(wire{}, arguments, input, out, failures, environment)
}

func (wire) Rejection() string {
	return defaultRejection
}

func (wire) Read(arguments []string, invocation *simulation.Invocation) error {
	parsed, err := parseArguments(arguments)
	if err != nil {
		return err
	}

	invocation.Model = parsed.model
	invocation.Effort = parsed.effort
	invocation.Schema = parsed.schema
	return nil
}

func (wire) Write(out io.Writer, response simulation.Response, invocation simulation.Invocation) error {
	if err := writeEvent(out, event{Type: "thread.started", ThreadID: threadOf(response)}); err != nil {
		return err
	}
	if err := writeEvent(out, event{Type: "turn.started"}); err != nil {
		return err
	}
	if response.Behaviour == simulation.ReportAnError {
		return writeEvent(out, event{Type: "turn.failed", Error: &failure{Message: response.Text}})
	}

	for turn := range response.Events {
		if turn > 0 && response.EveryInterval() > 0 {
			time.Sleep(response.EveryInterval())
		}
		if err := writeEvent(out, event{
			Type: "item.completed",
			Item: &item{ID: fmt.Sprintf("item_%d", turn+1), Type: "agent_message", Text: answerOf(response, turn+1 == response.Events, invocation)},
		}); err != nil {
			return err
		}
	}

	if response.Behaviour == simulation.OmitResult {
		return nil
	}
	if response.Behaviour == simulation.TruncateStream {
		_, err := io.WriteString(out, `{"type":"turn.completed","usage":{"input_toke`)
		return err
	}

	return writeEvent(out, event{
		Type: "turn.completed",
		Usage: &usage{
			InputTokens:       response.InputTokens,
			CachedInputTokens: 0,
			CacheWriteTokens:  0,
			OutputTokens:      response.OutputTokens,
		},
	})
}

func answerOf(response simulation.Response, last bool, invocation simulation.Invocation) string {
	if !last || invocation.Schema == "" || response.Structured == "" {
		return response.Text
	}
	return response.Structured
}

func threadOf(response simulation.Response) string {
	if response.SessionID != "" {
		return response.SessionID
	}
	return defaultThreadID
}

type usage struct {
	InputTokens       int `json:"input_tokens"`
	CachedInputTokens int `json:"cached_input_tokens"`
	CacheWriteTokens  int `json:"cache_write_input_tokens"`
	OutputTokens      int `json:"output_tokens"`
	ReasoningTokens   int `json:"reasoning_output_tokens"`
}

type item struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	Message string `json:"message,omitempty"`
}

type failure struct {
	Message string `json:"message"`
}

type event struct {
	Type     string   `json:"type"`
	ThreadID string   `json:"thread_id,omitempty"`
	Item     *item    `json:"item,omitempty"`
	Usage    *usage   `json:"usage,omitempty"`
	Error    *failure `json:"error,omitempty"`
	Message  string   `json:"message,omitempty"`
}

func writeEvent(out io.Writer, emitted event) error {
	var encoded bytes.Buffer
	encoder := json.NewEncoder(&encoded)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(emitted); err != nil {
		return fmt.Errorf("encoding a %s event: %w", emitted.Type, err)
	}
	if _, err := out.Write(encoded.Bytes()); err != nil {
		return fmt.Errorf("writing a %s event: %w", emitted.Type, err)
	}
	return nil
}

type request struct {
	model   string
	effort  string
	schema  string
	sandbox string
}

func parseArguments(arguments []string) (request, error) {
	if len(arguments) == 0 || arguments[0] != "exec" {
		return request{}, fmt.Errorf("codexsim answers the exec command only")
	}

	parsed := request{}
	streaming := false
	promptFromInput := false
	ignoresUserConfig := false
	suppressesProjectDoc := false

	for position := 1; position < len(arguments); position++ {
		switch argument := arguments[position]; argument {
		case "--json":
			streaming = true
		case "--ignore-user-config":
			ignoresUserConfig = true
		case promptFromStdin:
			promptFromInput = true
		case "--sandbox", "-s":
			value, err := valueAfter(arguments, position)
			if err != nil {
				return request{}, err
			}
			position++
			parsed.sandbox = value
		case "--output-schema":
			value, err := valueAfter(arguments, position)
			if err != nil {
				return request{}, err
			}
			position++
			held, readErr := os.ReadFile(value)
			if readErr != nil {
				return request{}, fmt.Errorf("reading the schema file %s: %w", value, readErr)
			}
			if !json.Valid(held) {
				return request{}, fmt.Errorf("the schema file %s is not JSON", value)
			}
			parsed.schema = string(held)
		case "-m", "--model":
			value, err := valueAfter(arguments, position)
			if err != nil {
				return request{}, err
			}
			position++
			parsed.model = value
		case "-c", "--config":
			value, err := valueAfter(arguments, position)
			if err != nil {
				return request{}, err
			}
			position++
			if effort, stated := strings.CutPrefix(value, effortOverride); stated {
				parsed.effort = effort
			}
			if value == projectDocCapped {
				suppressesProjectDoc = true
			}
		default:
			return request{}, fmt.Errorf("unexpected argument '%s' found", argument)
		}
	}

	if !streaming {
		return request{}, fmt.Errorf("codexsim answers a JSON event stream, and --json was not given")
	}
	if parsed.sandbox != sandboxReadOnly {
		return request{}, fmt.Errorf("codexsim runs read-only, and the sandbox was stated as '%s'", parsed.sandbox)
	}
	if !ignoresUserConfig {
		return request{}, fmt.Errorf("codexsim is not given the machine's own configuration, and --ignore-user-config was not given")
	}
	if !suppressesProjectDoc {
		return request{}, fmt.Errorf("codexsim does not read the reviewed repository's own instructions, and %s was not given", projectDocCapped)
	}
	if !promptFromInput {
		return request{}, fmt.Errorf("codexsim reads the prompt from standard input, and '%s' was not given", promptFromStdin)
	}
	return parsed, nil
}

func valueAfter(arguments []string, position int) (string, error) {
	if position+1 >= len(arguments) {
		return "", fmt.Errorf("option '%s' argument missing", arguments[position])
	}
	return arguments[position+1], nil
}
