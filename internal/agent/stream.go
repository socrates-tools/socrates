package agent

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	streamBufferSize  = 64 * 1024
	streamLineCeiling = 64 * 1024 * 1024
	elidedLength      = 120
)

type wireUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
}

type wireContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type wireMessage struct {
	Model   string        `json:"model"`
	Usage   wireUsage     `json:"usage"`
	Content []wireContent `json:"content"`
}

type wireEvent struct {
	Type             string          `json:"type"`
	Subtype          string          `json:"subtype"`
	SessionID        string          `json:"session_id"`
	Model            string          `json:"model"`
	Message          json.RawMessage `json:"message"`
	IsError          bool            `json:"is_error"`
	Result           string          `json:"result"`
	StructuredOutput json.RawMessage `json:"structured_output"`
	Usage            *wireUsage      `json:"usage"`
}

type streamed struct {
	result    *wireEvent
	text      string
	model     string
	sessionID string
	usage     Usage
}

func parseStream(input io.Reader) (streamed, error) {
	watched := &lastByteReader{inner: input}
	scanner := bufio.NewScanner(watched)
	scanner.Buffer(make([]byte, 0, streamBufferSize), streamLineCeiling)

	var read streamed
	malformed := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var event wireEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			malformed = line
			continue
		}
		if malformed != "" {
			return streamed{}, notAnEvent(malformed)
		}

		if event.SessionID != "" {
			read.sessionID = event.SessionID
		}
		if event.Model != "" {
			read.model = event.Model
		}

		switch event.Type {
		case "assistant":
			var message wireMessage
			if err := json.Unmarshal(event.Message, &message); err != nil {
				return streamed{}, fmt.Errorf(noMessageInAnAssistantEvent, err)
			}
			if message.Model != "" {
				read.model = message.Model
			}
			read.usage.add(message.Usage)
			for _, part := range message.Content {
				if part.Type == "text" {
					read.text += part.Text
				}
			}
		case "result":
			captured := event
			read.result = &captured
			if event.Usage != nil {
				read.usage = Usage{}
				read.usage.add(*event.Usage)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return streamed{}, fmt.Errorf(readingTheStream, err)
	}
	if malformed != "" {
		if watched.truncated() {
			return streamed{}, truncatedStream(malformed)
		}
		return streamed{}, notAnEvent(malformed)
	}
	if watched.truncated() {
		return streamed{}, truncatedStream("")
	}

	return read, nil
}

func truncatedStream(line string) error {
	if line == "" {
		return errors.New(truncatedMidLine)
	}
	return fmt.Errorf(truncatedMidLineAt, elide(line))
}

func notAnEvent(line string) error {
	return fmt.Errorf(notAnEventLine, elide(line))
}

func (u *Usage) add(reported wireUsage) {
	u.InputTokens += reported.InputTokens
	u.OutputTokens += reported.OutputTokens
	u.CacheReadTokens += reported.CacheReadInputTokens
	u.CacheCreationTokens += reported.CacheCreationInputTokens
}

type lastByteReader struct {
	inner io.Reader
	last  byte
	any   bool
}

func (r *lastByteReader) Read(into []byte) (int, error) {
	read, err := r.inner.Read(into)
	if read > 0 {
		r.last = into[read-1]
		r.any = true
	}
	return read, err
}

func (r *lastByteReader) truncated() bool {
	return r.any && r.last != '\n'
}

func elide(line string) string {
	if len(line) <= elidedLength {
		return line
	}
	return line[:elidedLength] + "..."
}
