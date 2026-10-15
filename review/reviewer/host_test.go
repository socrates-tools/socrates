package reviewer

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/review/agent"
)

func TestTheReviewerCLIIsTheStatedOneWhateverTheHost(t *testing.T) {
	for _, host := range []string{"", skill.ClaudeCode.Flag, skill.Codex.Flag} {
		got, err := ResolveReviewerCLI("codex", host)
		if err != nil || got != "codex" {
			t.Errorf("a stated codex on host %q resolves to %q, %v", host, got, err)
		}
	}
}

func TestAnUnsetReviewerCLIFollowsTheHostAndKeepsTheDefaultWithoutOne(t *testing.T) {
	for host, want := range map[string]string{"": agent.DefaultVendor, skill.ClaudeCode.Flag: agent.DefaultVendor, skill.Codex.Flag: "codex"} {
		got, err := ResolveReviewerCLI("", host)
		if err != nil || got != want {
			t.Errorf("an unset reviewer_cli on host %q resolves to %q, %v, want %q", host, got, err, want)
		}
	}
}

func TestAHostNamingNeitherHostIsRefusedNamingBoth(t *testing.T) {
	_, err := ResolveReviewerCLI("", "a-host-nobody-has")
	if err == nil {
		t.Fatal("a host nobody has was accepted")
	}
	for _, named := range []string{"a-host-nobody-has", skill.ClaudeCode.Flag, skill.Codex.Flag} {
		if !strings.Contains(err.Error(), named) {
			t.Errorf("the refusal reads %q, without naming %s", err, named)
		}
	}
}

func TestAHostNamingNeitherHostIsRefusedEvenWhereAReviewerCLIIsStated(t *testing.T) {
	if _, err := ResolveReviewerCLI("codex", "a-host-nobody-has"); err == nil {
		t.Fatal("a host nobody has was accepted beside a stated reviewer CLI")
	}
}

func TestAStatedReviewerCLINobodyDrivesIsRefusedByName(t *testing.T) {
	unsupported := agent.Unsupported()[0].Name
	if _, err := ResolveReviewerCLI(unsupported, skill.Codex.Flag); err == nil || !strings.Contains(err.Error(), unsupported) {
		t.Errorf("a stated %s on Codex was not refused by name: %v", unsupported, err)
	}
}
