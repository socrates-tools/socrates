package claudesim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type usage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
}

type content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type message struct {
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Role    string    `json:"role"`
	Model   string    `json:"model"`
	Content []content `json:"content"`
	Usage   usage     `json:"usage"`
}

type event struct {
	Type             string          `json:"type"`
	Subtype          string          `json:"subtype,omitempty"`
	SessionID        string          `json:"session_id"`
	Model            string          `json:"model,omitempty"`
	Tools            []string        `json:"tools,omitempty"`
	CurrentWorking   string          `json:"cwd,omitempty"`
	Message          *message        `json:"message,omitempty"`
	IsError          *bool           `json:"is_error,omitempty"`
	Result           string          `json:"result,omitempty"`
	StructuredOutput json.RawMessage `json:"structured_output,omitempty"`
	Usage            *usage          `json:"usage,omitempty"`
	DurationMs       int             `json:"duration_ms,omitempty"`
	NumTurns         int             `json:"num_turns,omitempty"`
	TotalCostUSD     float64         `json:"total_cost_usd,omitempty"`
}

func writeStream(out io.Writer, response Response, invocation Invocation) error {
	tools := strings.Split(invocation.Tools, ",")
	if len(invocation.Schema) > 0 {
		tools = append(tools, "StructuredOutput")
	}
	tools = append(tools, response.ExtraTools...)
	if response.Behaviour == OmitToolList {
		tools = nil
	}

	if err := writeEvent(out, event{
		Type:           "system",
		Subtype:        "init",
		SessionID:      response.SessionID,
		Model:          response.Model,
		Tools:          tools,
		CurrentWorking: invocation.Directory,
	}); err != nil {
		return err
	}

	spent := usage{
		InputTokens:  response.InputTokens,
		OutputTokens: response.OutputTokens,
	}

	for turn := range response.Events {
		if turn > 0 && response.interval() > 0 {
			time.Sleep(response.interval())
		}
		if err := writeEvent(out, event{
			Type:      "assistant",
			SessionID: response.SessionID,
			Message: &message{
				ID:      fmt.Sprintf("msg_%02d", turn+1),
				Type:    "message",
				Role:    "assistant",
				Model:   response.Model,
				Content: []content{{Type: "text", Text: response.Text}},
				Usage:   spent,
			},
		}); err != nil {
			return err
		}
	}

	if response.Behaviour == OmitResult {
		return nil
	}

	if response.Behaviour == TruncateStream {
		_, err := io.WriteString(out, `{"type":"result","subtype":"success","session_id":"`+response.SessionID+`","structured_ou`)
		return err
	}

	failed := false
	result := event{
		Type:         "result",
		Subtype:      "success",
		SessionID:    response.SessionID,
		IsError:      &failed,
		Result:       response.Text,
		Usage:        &spent,
		DurationMs:   1234,
		NumTurns:     response.Events,
		TotalCostUSD: 0.01,
	}
	if response.Structured != "" && invocation.Schema != "" {
		result.StructuredOutput = json.RawMessage(response.Structured)
	}
	return writeEvent(out, result)
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
