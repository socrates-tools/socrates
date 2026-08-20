package claudesim

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

type Behaviour string

const (
	ReportFindings   Behaviour = "findings"
	TruncateStream   Behaviour = "truncated-stream"
	OmitResult       Behaviour = "no-result"
	OmitStructured   Behaviour = "no-structured-output"
	ViolateSchema    Behaviour = "schema-violation"
	ExitNonZero      Behaviour = "non-zero-exit"
	SayNothingAtAll  Behaviour = "silence"
	RefuseStrictness Behaviour = "strict-tool-servers-refused"
)

const (
	strictnessRefused       = "You cannot use --strict-mcp-config when an enterprise MCP config is present"
	defaultText             = "I read the change against the base."
	defaultModel            = "claude-sim-1"
	defaultSession          = "00000000-0000-4000-8000-000000000000"
	defaultHold             = 10 * time.Minute
	defaultStderr           = "error: the model provider answered 529"
	defaultViolation        = `{"findings":"none that matter","risk_level":"low","risk_rationale":"r"}`
	defaultStructuredOutput = `{"findings":[],"risk_level":"low","risk_rationale":"nothing in this change needs attention"}`
)

var behaviours = []Behaviour{
	ReportFindings,
	RefuseStrictness,
	TruncateStream,
	OmitResult,
	OmitStructured,
	ViolateSchema,
	ExitNonZero,
	SayNothingAtAll,
}

type Write struct {
	Path           string `yaml:"path"`
	Content        string `yaml:"content"`
	DefeatsTheSeal bool   `yaml:"defeats_the_seal"`
}

type Response struct {
	When         string    `yaml:"when"`
	Behaviour    Behaviour `yaml:"behaviour"`
	Structured   string    `yaml:"structured"`
	Text         string    `yaml:"text"`
	Model        string    `yaml:"model"`
	SessionID    string    `yaml:"session_id"`
	Exit         int       `yaml:"exit"`
	Stderr       string    `yaml:"stderr"`
	InputTokens  int       `yaml:"input_tokens"`
	OutputTokens int       `yaml:"output_tokens"`
	Events       int       `yaml:"events"`
	Interval     string    `yaml:"interval"`
	Hold         string    `yaml:"hold"`
	Writes       []Write   `yaml:"writes"`
	ChildCommand []string  `yaml:"child_command"`
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
	}
	for _, response := range scenario.Responses {
		if err := response.validate(); err != nil {
			return nil, fmt.Errorf("the scenario %s states a response for %q that %w", path, response.When, err)
		}
	}

	return &scenario, nil
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
		r.Text = defaultText
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
		case OmitStructured:
			r.Structured = ""
		default:
			r.Structured = defaultStructuredOutput
		}
	}
	if r.Behaviour == ExitNonZero {
		if r.Exit == 0 {
			r.Exit = 1
		}
		if r.Stderr == "" {
			r.Stderr = defaultStderr
		}
	}
	return r
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
