package dialogue_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/review/dialogue"
	"github.com/socrates-tools/socrates/review/findings"
)

func TestAResponseIsAFixARejectionOrADecision(t *testing.T) {
	t.Run("Recording a fix", func(t *testing.T) {
		recorded := held(dialogue.Response{Finding: "f1", Resolution: dialogue.Fixed, Reason: "I moved the guard above the loop."})

		answer, answered := recorded.For("f1")
		if !answered || answer.Resolution != dialogue.Fixed {
			t.Fatalf("f1 reads as %+v", answer)
		}
		if answer.Reason != "I moved the guard above the loop." {
			t.Errorf("the note reads %q", answer.Reason)
		}
	})

	t.Run("Recording a rejection", func(t *testing.T) {
		const reason = "  The loop runs at most twice, and the caller's deadline bounds it.\nNothing retries forever.  "

		recorded := held(dialogue.Response{Finding: "f2", Resolution: dialogue.Rejected, Reason: reason})

		answer, _ := recorded.For("f2")
		if answer.Reason != reason {
			t.Errorf("the reason reads %q, want it verbatim", answer.Reason)
		}

		encoded, err := dialogue.Line(answer)
		if err != nil {
			t.Fatalf("encoding the response: %v", err)
		}
		read, err := dialogue.Parse(encoded)
		if err != nil {
			t.Fatalf("reading the responses back: %v", err)
		}
		if diff := cmp.Diff(recorded, read); diff != "" {
			t.Errorf("what was stored differs from what was read back (-stored +read):\n%s", diff)
		}
	})

	t.Run("A rejection with no reason", func(t *testing.T) {
		for _, blank := range []string{"", "   ", "\n\t "} {
			err := dialogue.Validate(dialogue.Response{Finding: "f1", Resolution: dialogue.Rejected, Reason: blank})
			if err == nil {
				t.Fatalf("a rejection reasoned %q was accepted", blank)
			}
			if !strings.Contains(err.Error(), "claim") {
				t.Errorf("the refusal does not say that a rejection is a claim: %v", err)
			}
		}
		if err := dialogue.Validate(dialogue.Response{Finding: "f1", Resolution: dialogue.Fixed}); err != nil {
			t.Errorf("a fix with no note was refused: %v", err)
		}
	})

	t.Run("Recording a decision", func(t *testing.T) {
		const reason = "  Vic decided the endpoint stays, because the migration reads it.\n"

		question := finding("f1", findings.Error, findings.AskUser)
		response := dialogue.Response{Finding: "f1", Resolution: dialogue.Decided, Reason: reason}
		if err := dialogue.Answering(question, response, dialogue.Mine); err != nil {
			t.Fatalf("a decision against a finding that asks a person was refused: %v", err)
		}

		recorded := held(response)
		answer, answered := recorded.For("f1")
		if !answered || answer.Resolution != dialogue.Decided {
			t.Fatalf("f1 reads as %+v, want a decision", answer)
		}
		if answer.Reason != reason {
			t.Errorf("the reason reads %q, want it verbatim", answer.Reason)
		}

		encoded, err := dialogue.Line(answer)
		if err != nil {
			t.Fatalf("encoding the response: %v", err)
		}
		read, err := dialogue.Parse(encoded)
		if err != nil {
			t.Fatalf("reading the responses back: %v", err)
		}
		if diff := cmp.Diff(recorded, read); diff != "" {
			t.Errorf("what was stored differs from what was read back (-stored +read):\n%s", diff)
		}
	})

	t.Run("A decision with no reason", func(t *testing.T) {
		for _, blank := range []string{"", "   ", "\n\t "} {
			err := dialogue.Validate(dialogue.Response{Finding: "f1", Resolution: dialogue.Decided, Reason: blank})
			if err == nil {
				t.Fatalf("a decision reasoned %q was accepted", blank)
			}
			if !strings.Contains(err.Error(), "claim") {
				t.Errorf("the refusal does not say that a decision is a claim: %v", err)
			}
		}
	})

	t.Run("A patch offered where a decision is required", func(t *testing.T) {
		question := finding("f1", findings.Error, findings.AskUser)

		for _, offered := range []dialogue.Response{
			{Finding: "f1", Resolution: dialogue.Fixed, Reason: "I took the endpoint out"},
			{Finding: "f1", Resolution: dialogue.Rejected, Reason: "the endpoint is needed"},
		} {
			err := dialogue.Answering(question, offered, dialogue.Mine)
			if err == nil {
				t.Fatalf("%q was accepted against a finding that asks a person", offered.Resolution)
			}
			if !strings.Contains(err.Error(), string(dialogue.Decided)) {
				t.Errorf("the refusal of %q does not name the resolution that applies: %v", offered.Resolution, err)
			}
		}
	})

	t.Run("A decision offered where none was asked for", func(t *testing.T) {
		decided := dialogue.Response{Finding: "f1", Resolution: dialogue.Decided, Reason: "Vic said to leave it"}

		for _, asking := range []findings.Action{findings.AutoFix, findings.NoOp} {
			err := dialogue.Answering(finding("f1", findings.Error, asking), decided, dialogue.Mine)
			if err == nil {
				t.Fatalf("a decision was accepted against a finding actioned %q", asking)
			}
			for _, named := range []dialogue.Resolution{dialogue.Fixed, dialogue.Rejected} {
				if !strings.Contains(err.Error(), string(named)) {
					t.Errorf("the refusal against a finding actioned %q does not name %q: %v", asking, named, err)
				}
			}
		}
	})

	t.Run("Changing a response", func(t *testing.T) {
		document := document(finding("f1", findings.Error, findings.AutoFix))

		recorded := held(
			dialogue.Response{Finding: "f1", Resolution: dialogue.Rejected, Reason: "the reviewer misread the guard"},
			dialogue.Response{Finding: "f1", Resolution: dialogue.Fixed, Reason: "on reflection it was right"},
		)

		if diff := cmp.Diff([]dialogue.Response{{Finding: "f1", Resolution: dialogue.Fixed, Reason: "on reflection it was right"}}, recorded.Latest()); diff != "" {
			t.Errorf("the next reviewer would receive (-want +got):\n%s", diff)
		}
		if len(recorded.Recorded) != 2 {
			t.Errorf("the run holds %d responses, so the replacement is not visible", len(recorded.Recorded))
		}
		if waiting := recorded.Unanswered(document, dialogue.Mine); waiting != nil {
			t.Errorf("f1 still reads as unanswered: %v", waiting)
		}
	})
}

func TestOnlyBlockingFindingsWaitForAResponse(t *testing.T) {
	// scenario: Non-blocking findings need no response
	document := document(
		finding("blocking", findings.Warning, findings.AutoFix),
		finding("informational", findings.Info, findings.NoOp),
	)

	if diff := cmp.Diff([]string{"blocking"}, dialogue.Responses{}.Unanswered(document, dialogue.Mine)); diff != "" {
		t.Errorf("the findings waiting for a response differ (-want +got):\n%s", diff)
	}

	recorded := held(dialogue.Response{Finding: "blocking", Resolution: dialogue.Fixed})
	if waiting := recorded.Unanswered(document, dialogue.Mine); waiting != nil {
		t.Errorf("an informational finding with no response is still waited on: %v", waiting)
	}
}

func TestAResponseToAFindingThatDoesNotExist(t *testing.T) {
	document := document(finding("f1", findings.Error, findings.AutoFix))

	if _, reported := dialogue.Reported(document, "f1"); !reported {
		t.Error("the round's own finding reads as one it did not produce")
	}
	if _, reported := dialogue.Reported(document, "f9"); reported {
		t.Error("a finding the round did not produce reads as one it did")
	}
}

func document(held ...findings.Finding) findings.Document {
	return findings.Document{
		Findings:      held,
		RiskLevel:     findings.LowRisk,
		RiskRationale: "nothing much",
	}
}

func held(responses ...dialogue.Response) dialogue.Responses {
	return dialogue.Responses{Recorded: responses}
}
