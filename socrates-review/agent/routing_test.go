package agent_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/simulator"
	"github.com/socrates-tools/socrates/socrates-review/agent"
	"github.com/socrates-tools/socrates/socrates-review/findings"
)

func TestTheReviewersModelIsConfigurableOnItsOwn(t *testing.T) {
	simulated := simulator.Install(t)

	reviewing, defaulted := agent.Default(agent.Reviewer)
	if !defaulted {
		t.Fatal("the reviewer has no default routing")
	}

	claude := build(t, agent.Options{
		Routing: agent.Configuration{
			agent.Reviewer: {Model: "sonnet"},
		},
	})

	if _, err := claude.Run(context.Background(), request()); err != nil {
		t.Fatalf("running the reviewer: %v", err)
	}

	invocation := simulated.OnlyInvocation()
	if invocation.Model != "sonnet" {
		t.Errorf("the reviewer was invoked with the model %q, want sonnet", invocation.Model)
	}
	if invocation.Effort != reviewing.Effort {
		t.Errorf("naming a model changed the effort to %q, want the default %q", invocation.Effort, reviewing.Effort)
	}

	for _, role := range agent.Roles() {
		if _, hasDefault := agent.Default(role); !hasDefault {
			t.Errorf("the role %q resolves to no default", role)
		}
	}
}

func TestTheReviewersEffortIsConfigurableOnItsOwn(t *testing.T) {
	simulated := simulator.Install(t)

	reviewing, _ := agent.Default(agent.Reviewer)
	claude := build(t, agent.Options{
		Routing: agent.Configuration{
			agent.Reviewer: {Effort: "max"},
		},
	})

	if _, err := claude.Run(context.Background(), request()); err != nil {
		t.Fatalf("running the reviewer: %v", err)
	}

	invocation := simulated.OnlyInvocation()
	if invocation.Effort != "max" {
		t.Errorf("the reviewer was invoked with the effort %q, want max", invocation.Effort)
	}
	if invocation.Model != reviewing.Model {
		t.Errorf("naming an effort changed the model to %q, want the default %q", invocation.Model, reviewing.Model)
	}
}

func TestAnEmptyConfigurationStillRuns(t *testing.T) {
	simulated := simulator.Install(t)

	claude := build(t, agent.Options{})

	if _, err := claude.Run(context.Background(), request()); err != nil {
		t.Fatalf("running the reviewer with no configuration: %v", err)
	}

	reviewing, _ := agent.Default(agent.Reviewer)
	invocation := simulated.OnlyInvocation()
	if invocation.Model != reviewing.Model || invocation.Effort != reviewing.Effort {
		t.Errorf("the reviewer was invoked with %q and %q, want the defaults %q and %q",
			invocation.Model, invocation.Effort, reviewing.Model, reviewing.Effort)
	}
}

func TestARoleSocratesDoesNotHaveIsRefusedAtConstruction(t *testing.T) {
	_, err := agent.NewClaude(agent.Options{
		Routing: agent.Configuration{
			agent.Role("librarian"): {Model: "sonnet"},
		},
	})
	if err == nil {
		t.Fatal("a role socrates does not have was accepted, and it would fall back at the call site")
	}
	if !strings.Contains(err.Error(), "librarian") {
		t.Errorf("the refusal reads %q, without naming the role", err)
	}
}

func TestARequestForAnUnknownRoleFails(t *testing.T) {
	simulator.Install(t)
	claude := build(t, agent.Options{})

	asked := request()
	asked.Role = agent.Role("librarian")

	if _, err := claude.Run(context.Background(), asked); err == nil {
		t.Fatal("a request for a role socrates does not have was answered")
	}
}

func build(t *testing.T, options agent.Options) agent.Agent {
	t.Helper()

	claude, err := agent.NewClaude(options)
	if err != nil {
		t.Fatalf("constructing the agent: %v", err)
	}
	return claude
}

func request() agent.Request {
	return agent.Request{
		Role:    agent.Reviewer,
		Prompt:  "You are reviewing a change you did not write.",
		Schema:  []byte(findings.Schema),
		Honours: honoursFindings,
	}
}

func honoursFindings(structured json.RawMessage) error {
	_, err := findings.Parse(structured)
	return err
}
