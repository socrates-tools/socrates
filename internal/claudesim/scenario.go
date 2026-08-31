package claudesim

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"

	"github.com/socrates-devtool/socrates/internal/corpus"
)

type Behaviour string

const (
	ReportFindings  Behaviour = "findings"
	TruncateStream  Behaviour = "truncated-stream"
	OmitResult      Behaviour = "no-result"
	OmitStructured  Behaviour = "no-structured-output"
	ViolateSchema   Behaviour = "schema-violation"
	ExitNonZero     Behaviour = "non-zero-exit"
	ReportAnError   Behaviour = "error-result"
	SayNothingAtAll Behaviour = "silence"
	RejectAtStartup Behaviour = "rejected-at-startup"
	OmitToolList    Behaviour = "no-tool-list"
)

const (
	defaultText             = "I read the change against the base."
	defaultModel            = "claude-sim-1"
	defaultSession          = "00000000-0000-4000-8000-000000000000"
	defaultHold             = 10 * time.Minute
	defaultStderr           = "error: the model provider answered 529"
	defaultAccount          = "API Error: 529 Overloaded. This is a server-side issue, usually temporary."
	defaultRejection        = "error: unknown option '--settings'"
	defaultViolation        = `{"findings":"none that matter","risk_level":"low","risk_rationale":"r"}`
	defaultStructuredOutput = `{"findings":[],"risk_level":"low","risk_rationale":"nothing in this change needs attention"}`
	connectedSource         = "connected"
)

var behaviours = []Behaviour{
	ReportFindings,
	RejectAtStartup,
	TruncateStream,
	OmitResult,
	OmitStructured,
	ViolateSchema,
	ExitNonZero,
	ReportAnError,
	SayNothingAtAll,
	OmitToolList,
}

type Write struct {
	Path           string `yaml:"path"`
	Content        string `yaml:"content"`
	DefeatsTheSeal bool   `yaml:"defeats_the_seal"`
}

type Response struct {
	When         string    `yaml:"when"`
	Behaviour    Behaviour `yaml:"behaviour"`
	Recorded     string    `yaml:"recorded"`
	Structured   string    `yaml:"structured"`
	Text         string    `yaml:"text"`
	Model        string    `yaml:"model"`
	SessionID    string    `yaml:"session_id"`
	Exit         *int      `yaml:"exit"`
	Stderr       string    `yaml:"stderr"`
	InputTokens  int       `yaml:"input_tokens"`
	OutputTokens int       `yaml:"output_tokens"`
	Events       int       `yaml:"events"`
	Interval     string    `yaml:"interval"`
	Hold         string    `yaml:"hold"`
	Writes       []Write   `yaml:"writes"`
	ChildCommand []string  `yaml:"child_command"`
	ExtraTools   []string  `yaml:"extra_tools"`
	ToolSources  []string  `yaml:"tool_sources"`
}

type Scenario struct {
	Default   *Response  `yaml:"default"`
	Responses []Response `yaml:"responses"`
}

func LoadScenario(path string) (*Scenario, error) {
	if path == "" {
		return &Scenario{}, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading the scenario %s: %w", path, err)
	}

	var scenario Scenario
	decoder := yaml.NewDecoder(strings.NewReader(string(content)))
	decoder.KnownFields(true)
	if err := decoder.Decode(&scenario); err != nil {
		return nil, fmt.Errorf("reading the scenario %s: %w", path, err)
	}

	if scenario.Default != nil {
		if err := scenario.Default.validate(); err != nil {
			return nil, fmt.Errorf("the scenario %s states a default that %w", path, err)
		}
		if err := replay(scenario.Default); err != nil {
			return nil, fmt.Errorf("the scenario %s states a default that %w", path, err)
		}
	}
	for position, response := range scenario.Responses {
		if err := response.validate(); err != nil {
			return nil, fmt.Errorf("the scenario %s states a response for %q that %w", path, response.When, err)
		}
		if err := replay(&scenario.Responses[position]); err != nil {
			return nil, fmt.Errorf("the scenario %s states a response for %q that %w", path, response.When, err)
		}
	}

	return &scenario, nil
}

func replay(response *Response) error {
	if response.Recorded == "" {
		return nil
	}
	if response.Structured != "" {
		return fmt.Errorf("names the recorded answer %s and carries a structured answer of its own, and a recorded answer is emitted as it was recorded", response.Recorded)
	}

	entry, err := corpus.Named(response.Recorded)
	if err != nil {
		return fmt.Errorf("names a recorded answer that %w", err)
	}
	response.Structured = string(entry.Document)
	return nil
}

func (s *Scenario) Match(prompt string) Response {
	for _, response := range s.Responses {
		if response.When != "" && strings.Contains(prompt, response.When) {
			return response.withDefaults()
		}
	}
	if s.Default != nil {
		return s.Default.withDefaults()
	}
	return Response{}.withDefaults()
}

func (r Response) validate() error {
	if r.Behaviour != "" && !isKnownBehaviour(r.Behaviour) {
		return fmt.Errorf("asks for the behaviour %q, which is not one of %v", r.Behaviour, behaviours)
	}
	if _, err := parseDuration(r.Interval); err != nil {
		return fmt.Errorf("states the interval %q, which is not a duration", r.Interval)
	}
	if _, err := parseDuration(r.Hold); err != nil {
		return fmt.Errorf("states the hold %q, which is not a duration", r.Hold)
	}
	return nil
}

func (r Response) withDefaults() Response {
	if r.Behaviour == "" {
		r.Behaviour = ReportFindings
	}
	if r.Text == "" {
		if r.Behaviour == ReportAnError {
			r.Text = defaultAccount
		} else {
			r.Text = defaultText
		}
	}
	if r.Model == "" {
		r.Model = defaultModel
	}
	if r.SessionID == "" {
		r.SessionID = defaultSession
	}
	if r.Events < 1 {
		r.Events = 1
	}
	if r.InputTokens == 0 {
		r.InputTokens = 100
	}
	if r.OutputTokens == 0 {
		r.OutputTokens = 50
	}
	if r.Structured == "" {
		switch r.Behaviour {
		case ViolateSchema:
			r.Structured = defaultViolation
		case OmitStructured, ReportAnError:
			r.Structured = ""
		default:
			r.Structured = defaultStructuredOutput
		}
	}
	if r.Behaviour == ExitNonZero || r.Behaviour == ReportAnError {
		if r.Exit == nil {
			failing := exitUsage
			r.Exit = &failing
		}
	}
	if r.Behaviour == ExitNonZero && r.Stderr == "" {
		r.Stderr = defaultStderr
	}
	return r
}

func (r Response) exitCode() int {
	if r.Exit == nil {
		return exitSuccess
	}
	return *r.Exit
}

func (r Response) reportsAnError() bool {
	return r.Behaviour == ReportAnError
}

func (r Response) interval() time.Duration {
	parsed, _ := parseDuration(r.Interval)
	return parsed
}

func (r Response) hold() time.Duration {
	parsed, _ := parseDuration(r.Hold)
	if parsed == 0 {
		return defaultHold
	}
	return parsed
}

func parseDuration(stated string) (time.Duration, error) {
	if stated == "" {
		return 0, nil
	}
	return time.ParseDuration(stated)
}

func isKnownBehaviour(behaviour Behaviour) bool {
	return slices.Contains(behaviours, behaviour)
}
