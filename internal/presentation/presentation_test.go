package presentation_test

import (
	"strings"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/agent"
	"github.com/rui-ktei/socrates/internal/presentation"
	"github.com/rui-ktei/socrates/internal/reviewer"
	"github.com/rui-ktei/socrates/internal/runlog"
	"github.com/rui-ktei/socrates/internal/runstate"
)

var at = time.Date(2026, 8, 21, 9, 14, 2, 0, time.UTC)

func TestRenderingNeedsNoTerminal(t *testing.T) {
	lines := presentation.Lines([]runlog.Event{
		{At: at, Kind: reviewer.EventRunStarted, Subject: "20260821T091402Z-6616db11", Fields: map[string]any{
			"repository": "/home/a/project",
			"base":       "a9a3207995f002632584efb33312c14123de4560",
			"snapshot":   "4445282d62930b67e699165cda6a212f8c2e47ab",
		}},
		{At: at.Add(time.Second), Kind: agent.EventStarted, Subject: "reviewer", Fields: map[string]any{
			"model":   "opus",
			"effort":  "high",
			"process": float64(4242),
		}},
		{At: at.Add(90 * time.Second), Kind: agent.EventCompleted, Subject: "reviewer", Fields: map[string]any{
			"model":         "claude-opus-5",
			"input_tokens":  float64(18234),
			"output_tokens": float64(2101),
		}},
		{At: at.Add(91 * time.Second), Kind: reviewer.EventRunEnded, Subject: "20260821T091402Z-6616db11", Fields: map[string]any{
			"state":  string(runstate.Escalated),
			"reason": "three findings, one of them for the author",
		}},
	})

	if len(lines) != 4 {
		t.Fatalf("four events rendered as %d lines:\n%s", len(lines), strings.Join(lines, "\n"))
	}
	for position, needed := range [][]string{
		{"09:14:02", "/home/a/project", "4445282", "a9a3207"},
		{"09:14:03", "reviewer", "opus", "high"},
		{"09:15:32", "18234", "2101"},
		{"09:15:33", string(runstate.Escalated), "three findings, one of them for the author"},
	} {
		for _, wanted := range needed {
			if !strings.Contains(lines[position], wanted) {
				t.Errorf("line %d does not carry %q: %q", position+1, wanted, lines[position])
			}
		}
	}
	for _, line := range lines {
		if strings.Contains(line, "map[") || strings.Contains(line, "%!") {
			t.Errorf("a line was rendered from a field it could not read: %q", line)
		}
	}
}

func TestEveryKindTheCoreWritesIsRendered(t *testing.T) {
	for _, kind := range []string{
		reviewer.EventRunStarted,
		reviewer.EventRunEnded,
		reviewer.EventRoundStarted,
		reviewer.EventRoundParked,
		reviewer.EventResponse,
		agent.EventStarted,
		agent.EventCompleted,
		agent.EventFailed,
		agent.EventToolServersLoaded,
	} {
		t.Run(kind, func(t *testing.T) {
			lines := presentation.Lines([]runlog.Event{{At: at, Kind: kind, Subject: "a-run", Fields: fields()}})
			if len(lines) != 1 {
				t.Fatalf("one event rendered as %d lines: %v", len(lines), lines)
			}
			if strings.Contains(lines[0], kind) {
				t.Errorf("the kind %s is rendered by printing its own name: %q", kind, lines[0])
			}
			if strings.Contains(lines[0], "%!") || strings.Contains(lines[0], "map[") {
				t.Errorf("the kind %s is rendered from a field it could not read: %q", kind, lines[0])
			}
		})
	}
}

func TestAParkedRoundRendersWhatItWaitsOn(t *testing.T) {
	rendered := func(held map[string]any) string {
		t.Helper()

		lines := presentation.Lines([]runlog.Event{{At: at, Kind: reviewer.EventRoundParked, Subject: "a-run", Fields: held}})
		if len(lines) != 1 {
			t.Fatalf("a parked round rendered as %d lines: %v", len(lines), lines)
		}
		return lines[0]
	}

	t.Run("a round parked on a question alone", func(t *testing.T) {
		line := rendered(map[string]any{"round": float64(1), "patchable": float64(0), "needing_a_person": float64(2)})

		for _, needed := range []string{"0 findings can be patched", "2 need a person"} {
			if !strings.Contains(line, needed) {
				t.Errorf("a round waiting on two questions and no patches does not say %q: %q", needed, line)
			}
		}
	})

	t.Run("a round parked on patches alone", func(t *testing.T) {
		line := rendered(map[string]any{"round": float64(1), "patchable": float64(3), "needing_a_person": float64(0)})

		if !strings.Contains(line, "3 findings can be patched") {
			t.Errorf("a round waiting on three patches does not say so: %q", line)
		}
		if strings.Contains(line, "person") {
			t.Errorf("a round waiting on nobody names a person: %q", line)
		}
	})

	t.Run("a round parked on both", func(t *testing.T) {
		line := rendered(map[string]any{"round": float64(2), "patchable": float64(3), "needing_a_person": float64(1)})

		for _, needed := range []string{"3 findings can be patched", "1 need a person"} {
			if !strings.Contains(line, needed) {
				t.Errorf("a round waiting on both does not carry %q: %q", needed, line)
			}
		}
	})
}

func TestAnUnknownEventKindIsStillShown(t *testing.T) {
	lines := presentation.Lines([]runlog.Event{
		{At: at, Kind: "reviewer-thought-about-it", Subject: "r1", Fields: map[string]any{
			"depth": float64(3),
			"about": "the retry loop",
		}},
	})

	if len(lines) != 1 {
		t.Fatalf("an unknown kind rendered as %d lines: %v", len(lines), lines)
	}
	for _, needed := range []string{"09:14:02", "reviewer-thought-about-it", "r1", "depth", "3", "about", "the retry loop"} {
		if !strings.Contains(lines[0], needed) {
			t.Errorf("the rendering of an unknown kind drops %q: %q", needed, lines[0])
		}
	}
}

func TestARotatedLogIsMarked(t *testing.T) {
	lines := presentation.Lines([]runlog.Event{
		{At: at, Kind: agent.EventStarted, Subject: "reviewer", Fields: fields()},
		{At: at.Add(time.Second), Kind: runlog.EventRotated, Subject: runlog.FileName, Fields: map[string]any{
			"ceiling": float64(1024),
			"kept":    float64(512),
		}},
		{At: at.Add(2 * time.Second), Kind: agent.EventCompleted, Subject: "reviewer", Fields: fields()},
	})

	if len(lines) != 3 {
		t.Fatalf("a rotated log rendered as %d lines:\n%s", len(lines), strings.Join(lines, "\n"))
	}
	if !strings.Contains(lines[0], "discarded") {
		t.Errorf("the first line does not say that earlier events are gone: %q", lines[0])
	}
	for _, line := range lines[1:] {
		if strings.Contains(line, "discarded") {
			t.Errorf("the rotation is reported twice: %q", line)
		}
	}
}

func fields() map[string]any {
	return map[string]any{
		"repository":    "/home/a/project",
		"base":          "a9a3207995f002632584efb33312c14123de4560",
		"snapshot":      "4445282d62930b67e699165cda6a212f8c2e47ab",
		"round":         float64(1),
		"model":         "opus",
		"effort":        "high",
		"process":       float64(4242),
		"input_tokens":  float64(10),
		"output_tokens": float64(2),
		"reason":        "the model provider answered 529",
		"state":         string(runstate.Failed),
		"patchable":     float64(2),
		"finding":       "f1",
		"resolution":    "fixed",
		"unanswered":    float64(1),
	}
}
