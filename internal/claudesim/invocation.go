package claudesim

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	ScenarioVariable  = "CLAUDESIM_SCENARIO"
	LogVariable       = "CLAUDESIM_LOG"
	RecordEnvVariable = "CLAUDESIM_RECORD_ENV"
)

var recordedPrefixes = []string{"SOCRATES_", "CLAUDESIM_", "CLAUDE_"}

type Invocation struct {
	Arguments      []string          `json:"arguments"`
	Prompt         string            `json:"prompt"`
	Directory      string            `json:"directory"`
	Environment    map[string]string `json:"environment"`
	Model          string            `json:"model,omitempty"`
	Effort         string            `json:"effort,omitempty"`
	Schema         string            `json:"schema,omitempty"`
	SettingSources string            `json:"setting_sources,omitempty"`
	Settings       string            `json:"settings,omitempty"`
	Tools          string            `json:"tools,omitempty"`
	WritesRefused  []string          `json:"writes_refused,omitempty"`
	Child          int               `json:"child,omitempty"`
	Process        int               `json:"process"`
}

func (i Invocation) Has(argument string) bool {
	return slices.Contains(i.Arguments, argument)
}

func (i Invocation) ValueOf(flag string) (string, bool) {
	for position, argument := range i.Arguments {
		if argument == flag && position+1 < len(i.Arguments) {
			return i.Arguments[position+1], true
		}
		if value, found := strings.CutPrefix(argument, flag+"="); found {
			return value, true
		}
	}
	return "", false
}

func Append(path string, invocation Invocation) error {
	if path == "" {
		return nil
	}

	encoded, err := json.Marshal(invocation)
	if err != nil {
		return fmt.Errorf("encoding the invocation: %w", err)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("opening %s: %w", path, err)
	}
	if _, err := file.Write(append(encoded, '\n')); err != nil {
		_ = file.Close()
		return fmt.Errorf("writing to %s: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("closing %s: %w", path, err)
	}
	return nil
}

func ReadLog(path string) ([]Invocation, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	var invocations []Invocation
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var invocation Invocation
		if err := json.Unmarshal(line, &invocation); err != nil {
			return nil, fmt.Errorf("reading a line of %s: %w", path, err)
		}
		invocations = append(invocations, invocation)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	return invocations, nil
}

func recordedEnvironment(environment []string) map[string]string {
	named := map[string]bool{}
	for _, name := range strings.Split(os.Getenv(RecordEnvVariable), ",") {
		if trimmed := strings.TrimSpace(name); trimmed != "" {
			named[trimmed] = true
		}
	}

	recorded := map[string]string{}
	for _, entry := range environment {
		name, value, found := strings.Cut(entry, "=")
		if !found {
			continue
		}
		if named[name] || hasRecordedPrefix(name) {
			recorded[name] = value
		}
	}
	return recorded
}

func hasRecordedPrefix(name string) bool {
	for _, prefix := range recordedPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
