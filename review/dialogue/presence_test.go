package dialogue_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/review/dialogue"
	"github.com/socrates-tools/socrates/review/findings"
	"github.com/socrates-tools/socrates/review/runstate"
)

func TestAFindingThatAsksAPersonTakesASettlementWhereNobodyIsToAsk(t *testing.T) {
	question := finding("f1", findings.Error, findings.AskUser)

	if diff := cmp.Diff([]dialogue.Resolution{dialogue.Settled}, dialogue.Accepts(question, dialogue.Mine, dialogue.NobodyToAsk)); diff != "" {
		t.Errorf("a finding that asks a person on a run with nobody to ask accepts (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]dialogue.Resolution{dialogue.Decided}, dialogue.Accepts(question, dialogue.Mine, dialogue.SomebodyToAsk)); diff != "" {
		t.Errorf("a finding that asks a person on a run with somebody to ask accepts (-want +got):\n%s", diff)
	}

	settled := dialogue.Response{Finding: "f1", Resolution: dialogue.Settled, Reason: "the function stays, because the importer reads it"}
	if err := dialogue.Answering(question, settled, dialogue.Mine, dialogue.NobodyToAsk); err != nil {
		t.Errorf("a settlement of a finding that asks a person on a run with nobody to ask was refused: %v", err)
	}
}

func TestADecisionIsRefusedWhereNobodyIsToAskNamingTheSettlement(t *testing.T) {
	question := finding("f1", findings.Error, findings.AskUser)

	for _, offered := range []dialogue.Resolution{dialogue.Decided, dialogue.Fixed, dialogue.Rejected} {
		err := dialogue.Answering(question, dialogue.Response{Finding: "f1", Resolution: offered, Reason: "because"}, dialogue.Mine, dialogue.NobodyToAsk)
		if err == nil {
			t.Fatalf("%q was accepted against a finding that asks a person on a run with nobody to ask", offered)
		}
		if !strings.Contains(err.Error(), "--"+string(dialogue.Settled)) {
			t.Errorf("the refusal of %q does not name --%s, which is how the finding is answered: %v", offered, dialogue.Settled, err)
		}
	}
}

func TestASettlementIsRefusedWhereSomebodyIsToAskNamingTheDecision(t *testing.T) {
	settled := dialogue.Response{Finding: "f1", Resolution: dialogue.Settled, Reason: "the function stays"}

	err := dialogue.Answering(finding("f1", findings.Error, findings.AskUser), settled, dialogue.Mine, dialogue.SomebodyToAsk)
	if err == nil {
		t.Fatal("a settlement was accepted on a run with somebody to ask")
	}
	if !strings.Contains(err.Error(), "--"+string(dialogue.Decided)) {
		t.Errorf("the refusal does not name --%s, which is how the finding is answered: %v", dialogue.Decided, err)
	}
}

func TestASettlementAnswersNothingButAFindingThatAsksAPerson(t *testing.T) {
	settled := dialogue.Response{Finding: "f1", Resolution: dialogue.Settled, Reason: "the function stays"}

	for _, asking := range []findings.Action{findings.AutoFix, findings.NoOp} {
		reported := finding("f1", findings.Error, asking)
		if slices.Contains(dialogue.Accepts(reported, dialogue.Mine, dialogue.NobodyToAsk), dialogue.Settled) {
			t.Errorf("a finding actioned %q accepts a settlement", asking)
		}
		if err := dialogue.Answering(reported, settled, dialogue.Mine, dialogue.NobodyToAsk); err == nil {
			t.Errorf("a settlement was accepted against a finding actioned %q", asking)
		}
	}
}

func TestASettlementNeedsAReasonAndRoundTrips(t *testing.T) {
	for _, blank := range []string{"", "   ", "\n\t "} {
		err := dialogue.Validate(dialogue.Response{Finding: "f1", Resolution: dialogue.Settled, Reason: blank})
		if err == nil {
			t.Fatalf("a settlement reasoned %q was accepted", blank)
		}
		if !strings.Contains(err.Error(), "settlement") || !strings.Contains(err.Error(), "needs a reason") {
			t.Errorf("the refusal does not say a settlement needs a reason: %v", err)
		}
	}

	recorded := held(dialogue.Response{Finding: "f1", Resolution: dialogue.Settled, Reason: "the function stays"})
	answer, answered := recorded.For("f1")
	if !answered || answer.Resolution != dialogue.Settled {
		t.Fatalf("f1 reads as %+v, want a settlement", answer)
	}
}

func TestARoundOfQuestionsDoesNotAwaitAPersonWhereNobodyIsToAsk(t *testing.T) {
	asking := document(finding("f1", findings.Error, findings.AskUser))
	classified := dialogue.Classified(asking, dialogue.Mine)

	if !dialogue.AwaitsAPerson(classified, runstate.AwaitingResponse, dialogue.SomebodyToAsk) {
		t.Error("a parked round holding a question does not await a person where somebody is to ask")
	}
	if dialogue.AwaitsAPerson(classified, runstate.AwaitingResponse, dialogue.NobodyToAsk) {
		t.Error("a parked round holding a question awaits a person where nobody is to ask")
	}
	if !dialogue.AwaitsAPerson(classified, runstate.Escalated, dialogue.NobodyToAsk) {
		t.Error("an escalated round does not await a person where nobody is to ask, and nobody settled it")
	}

	waiting := (dialogue.Responses{}).Unanswered(asking, dialogue.Mine)
	if diff := cmp.Diff([]string{"f1"}, waiting); diff != "" {
		t.Errorf("the question is still unanswered where nobody is to ask (-want +got):\n%s", diff)
	}
	if asked := (dialogue.Responses{}).AskingAPerson(asking, dialogue.Mine); !slices.Equal(asked, []string{"f1"}) {
		t.Errorf("%v ask a person, want f1", asked)
	}
}

func TestTheResponseRulesSaySettledBelongsToARunWithNobodyToAsk(t *testing.T) {
	stated := strings.Join(dialogue.Rules(), "\n")

	for _, said := range []string{"`settled` belongs only to a run begun with `--nobody-to-ask`", "which you never begin"} {
		if !strings.Contains(stated, said) {
			t.Errorf("the response rules do not say %q: %s", said, stated)
		}
	}
	if strings.Contains(stated, "socrates review begin --nobody-to-ask") {
		t.Errorf("the response rules tell the session to begin a run with nobody to ask: %s", stated)
	}
}
