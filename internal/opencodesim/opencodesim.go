package opencodesim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/socrates-tools/socrates/internal/simulation"
)

const (
	defaultRejection     = "error: Unknown argument: settings"
	defaultSession       = "ses_00000000000000000000000000"
	suppressesProjectDoc = "OPENCODE_DISABLE_PROJECT_CONFIG"
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
	if os.Getenv(suppressesProjectDoc) != "1" {
		return fmt.Errorf("opencodesim does not read the reviewed repository's own instructions, and %s was not set", suppressesProjectDoc)
	}

	invocation.Model = parsed.model
	invocation.Effort = parsed.variant
	return nil
}

func (wire) Write(out io.Writer, response simulation.Response, invocation simulation.Invocation) error {
	session := response.SessionID
	if session == "" {
		session = defaultSession
	}

	if response.Behaviour == simulation.ReportAnError {
		return writeEvent(out, event{
			Type:      "error",
			Timestamp: 1,
			SessionID: session,
			Error:     &failure{Name: "UnknownError", Data: failureData{Message: response.Text, Ref: "err_0"}},
		})
	}

	if err := writeEvent(out, event{Type: "step_start", Timestamp: 1, SessionID: session, Part: &part{Type: "step-start", ID: "prt_0"}}); err != nil {
		return err
	}

	for turn := range response.Events {
		if turn > 0 && response.EveryInterval() > 0 {
			time.Sleep(response.EveryInterval())
		}
		if err := writeEvent(out, event{
			Type:      "text",
			Timestamp: 2,
			SessionID: session,
			Part:      &part{Type: "text", ID: fmt.Sprintf("prt_%d", turn+1), Text: response.Spoken(turn+1 == response.Events, true)},
		}); err != nil {
			return err
		}
	}

	if response.Behaviour == simulation.OmitResult {
		return nil
	}
	if response.Behaviour == simulation.TruncateStream {
		_, err := io.WriteString(out, `{"type":"step_finish","timestamp":4,"sessionID":"`+session+`","part":{"tok`)
		return err
	}

	return writeEvent(out, event{
		Type:      "step_finish",
		Timestamp: 4,
		SessionID: session,
		Part: &part{
			Type:   "step-finish",
			ID:     "prt_finish",
			Reason: "stop",
			Tokens: &tokens{Input: response.InputTokens, Output: response.OutputTokens},
		},
	})
}

type tokens struct {
	Total     int `json:"total"`
	Input     int `json:"input"`
	Output    int `json:"output"`
	Reasoning int `json:"reasoning"`
	Cache     struct {
		Write int `json:"write"`
		Read  int `json:"read"`
	} `json:"cache"`
}

type part struct {
	ID     string  `json:"id"`
	Type   string  `json:"type"`
	Text   string  `json:"text,omitempty"`
	Reason string  `json:"reason,omitempty"`
	Tokens *tokens `json:"tokens,omitempty"`
}

type failureData struct {
	Message string `json:"message"`
	Ref     string `json:"ref"`
}

type failure struct {
	Name string      `json:"name"`
	Data failureData `json:"data"`
}

type event struct {
	Type      string   `json:"type"`
	Timestamp int64    `json:"timestamp"`
	SessionID string   `json:"sessionID"`
	Part      *part    `json:"part,omitempty"`
	Error     *failure `json:"error,omitempty"`
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
	model     string
	variant   string
	directory string
}

func parseArguments(arguments []string) (request, error) {
	if len(arguments) == 0 || arguments[0] != "run" {
		return request{}, fmt.Errorf("opencodesim answers the run command only")
	}

	parsed := request{}
	format := ""

	for position := 1; position < len(arguments); position++ {
		argument := arguments[position]
		if argument == "" || argument[0] != '-' {
			return request{}, fmt.Errorf("opencodesim reads the prompt from standard input, and was given the argument %q", argument)
		}

		value, err := valueAfter(arguments, position)
		if err != nil {
			return request{}, err
		}
		position++

		switch argument {
		case "--format":
			format = value
		case "--dir":
			parsed.directory = value
		case "-m", "--model":
			parsed.model = value
		case "--variant":
			parsed.variant = value
		default:
			return request{}, fmt.Errorf("Unknown argument: %s", argument)
		}
	}

	if format != "json" {
		return request{}, fmt.Errorf("opencodesim answers raw JSON events, and the format was stated as '%s'", format)
	}
	if parsed.directory == "" {
		return request{}, fmt.Errorf("opencodesim reads the directory it is given, and --dir was not given")
	}
	return parsed, nil
}

func valueAfter(arguments []string, position int) (string, error) {
	if position+1 >= len(arguments) {
		return "", fmt.Errorf("Not enough arguments following: %s", arguments[position])
	}
	return arguments[position+1], nil
}
