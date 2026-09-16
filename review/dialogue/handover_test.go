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

func TestEveryFindingOnAChangeThatIsNotMineIsForwardedDroppedOrRejected(t *testing.T) {
	want := []dialogue.Resolution{dialogue.Forwarded, dialogue.Dropped, dialogue.Rejected}

	for _, one := range everyKind() {
		reported := finding("f1", one.severity, one.action)

		if diff := cmp.Diff(want, dialogue.Accepts(reported, dialogue.NotMine)); diff != "" {
			t.Errorf("a %s/%s finding on a change that is not mine accepts (-want +got):\n%s", one.severity, one.action, diff)
		}
		for _, offered := range want {
			answer := dialogue.Response{Finding: "f1", Resolution: offered, Reason: "what I saw"}
			if err := dialogue.Answering(reported, answer, dialogue.NotMine); err != nil {
				t.Errorf("%q against a %s/%s finding on a change that is not mine was refused: %v", offered, one.severity, one.action, err)
			}
		}
	}
}

func TestAResolutionThatAssumesIWillPatchIsRefusedOnAChangeThatIsNotMine(t *testing.T) {
	reported := finding("f1", findings.Error, findings.AutoFix)

	for _, offered := range []dialogue.Resolution{dialogue.Fixed, dialogue.Decided} {
		err := dialogue.Answering(reported, dialogue.Response{Finding: "f1", Resolution: offered, Reason: "it is done"}, dialogue.NotMine)
		if err == nil {
			t.Fatalf("%q was accepted against a change that is not mine to patch", offered)
		}
		for _, named := range []dialogue.Resolution{dialogue.Forwarded, dialogue.Dropped, dialogue.Rejected} {
			if !strings.Contains(err.Error(), string(named)) {
				t.Errorf("the refusal of %q does not name %q, which is one of the answers that apply: %v", offered, named, err)
			}
		}
	}
}

func TestForwardingAndDroppingAreRefusedOnMyOwnChange(t *testing.T) {
	for _, one := range everyKind() {
		reported := finding("f1", one.severity, one.action)

		for _, offered := range []dialogue.Resolution{dialogue.Forwarded, dialogue.Dropped} {
			if slices.Contains(dialogue.Accepts(reported, dialogue.Mine), offered) {
				t.Errorf("a %s/%s finding on my own change accepts %q, and there is nobody to hand it to", one.severity, one.action, offered)
			}

			err := dialogue.Answering(reported, dialogue.Response{Finding: "f1", Resolution: offered, Reason: "not worth raising"}, dialogue.Mine)
			if err == nil {
				t.Fatalf("%q was accepted against a %s/%s finding on my own change", offered, one.severity, one.action)
			}
			if !strings.Contains(err.Error(), "--not-my-change") {
				t.Errorf("the refusal of %q does not name the flag that says the change is somebody else's: %v", offered, err)
			}
			for _, applies := range dialogue.Accepts(reported, dialogue.Mine) {
				if !strings.Contains(err.Error(), "--"+string(applies)) {
					t.Errorf("the refusal of %q against a %s/%s finding does not name --%s, which is how it is answered here: %v", offered, one.severity, one.action, applies, err)
				}
			}
			for _, elsewhere := range []dialogue.Resolution{dialogue.Forwarded, dialogue.Dropped} {
				if strings.Contains(err.Error(), "--"+string(elsewhere)+",") || strings.Contains(err.Error(), "or --"+string(elsewhere)) {
					t.Errorf("the refusal of %q offers --%s, which is what it just refused: %v", offered, elsewhere, err)
				}
			}
		}
	}
}

func TestAResolutionThatLeavesNoArtifactCarriesAReason(t *testing.T) {
	needsAReason := []dialogue.Resolution{dialogue.Rejected, dialogue.Decided, dialogue.Dropped}

	for _, resolution := range dialogue.Resolutions() {
		wanted := slices.Contains(needsAReason, resolution)

		for _, blank := range []string{"", "   ", "\n\t "} {
			err := dialogue.Validate(dialogue.Response{Finding: "f1", Resolution: resolution, Reason: blank})
			if wanted && err == nil {
				t.Errorf("%q leaves nothing behind and was accepted with the reason %q", resolution, blank)
			}
			if !wanted && err != nil {
				t.Errorf("%q leaves an artifact behind and was refused for having no reason: %v", resolution, err)
			}
		}
		if err := dialogue.Validate(dialogue.Response{Finding: "f1", Resolution: resolution, Reason: "because"}); err != nil {
			t.Errorf("%q with a reason was refused: %v", resolution, err)
		}
	}
}

func TestEveryFindingOfAChangeThatIsNotMineIsAPersons(t *testing.T) {
	kinds := everyKind()

	for selection := range 1 << len(kinds) {
		held := selected(kinds, selection)
		classified := dialogue.Classified(document(held...), dialogue.NotMine)

		if len(classified.Patchable) != 0 {
			t.Fatalf("%v on a change that is not mine reads %d findings as patchable, and nothing here can be patched", describe(held), len(classified.Patchable))
		}
		if len(classified.NoOps) != 0 {
			t.Fatalf("%v on a change that is not mine reads %d findings as asking for nothing, and forwarding one is something", describe(held), len(classified.NoOps))
		}
		if len(classified.NeedingAPerson) != len(held) {
			t.Fatalf("%v on a change that is not mine reads %d of %d findings as a person's", describe(held), len(classified.NeedingAPerson), len(held))
		}
	}
}

func TestEveryRoundOfAChangeThatIsNotMineWaitsOnAPerson(t *testing.T) {
	kinds := everyKind()

	for selection := range 1 << len(kinds) {
		held := selected(kinds, selection)
		classified := dialogue.Classified(document(held...), dialogue.NotMine)

		for round := range dialogue.Cap + 2 {
			asked := round + 1
			reached := dialogue.After(classified, asked)

			if blocks(held) && asked < dialogue.Cap && reached != runstate.AwaitingResponse {
				t.Fatalf("round %d of %v on a change that is not mine is %q, and forwarding is what addresses every one of them", asked, describe(held), reached)
			}

			got := dialogue.AwaitsAPerson(classified, reached)
			want := reached != runstate.Clean
			if got != want {
				t.Fatalf("round %d of %v on a change that is not mine led to %q and waits on a person: %t, want %t", asked, describe(held), reached, got, want)
			}
		}
	}
}

func TestARoundOfAChangeThatIsNotMineWaitsOnEveryFindingThatHoldsIt(t *testing.T) {
	kinds := everyKind()

	for selection := range 1 << len(kinds) {
		held := selected(kinds, selection)

		waiting := (dialogue.Responses{}).Unanswered(document(held...), dialogue.NotMine)
		if diff := cmp.Diff(blockingIdentifiers(held), waiting); diff != "" {
			t.Fatalf("%v on a change that is not mine waits on (-want +got):\n%s", describe(held), diff)
		}
		for _, unanswered := range waiting {
			asked := (dialogue.Responses{}).AwaitingAPerson(document(held...), dialogue.NotMine)
			if !slices.Contains(asked, unanswered) {
				t.Fatalf("%s of %v is unanswered and reads as nobody's to answer, and the whole round is a person's", unanswered, describe(held))
			}
		}
	}
}

func TestOneRoundIsDerivedRatherThanCapped(t *testing.T) {
	beyondAPatch := document(finding("f1", findings.Error, findings.NoOp))

	if got := dialogue.After(dialogue.Classified(beyondAPatch, dialogue.Mine), 1); got != runstate.Escalated {
		t.Errorf("a finding that holds my own change while asking for nothing is %q, want %q", got, runstate.Escalated)
	}
	if got := dialogue.After(dialogue.Classified(beyondAPatch, dialogue.NotMine), 1); got != runstate.AwaitingResponse {
		t.Errorf("the same finding on a change that is not mine is %q, want %q: forwarding it is what addresses it, and a round that ended could not be forwarded from", got, runstate.AwaitingResponse)
	}
	if dialogue.Cap != 2 {
		t.Errorf("the round cap reads %d, and nothing about a change that is not mine moves it", dialogue.Cap)
	}
}

func TestTheHandoverRulesAreStatedOnce(t *testing.T) {
	rules := dialogue.HandoverRules()
	if len(rules) == 0 {
		t.Fatal("nothing states how a finding on somebody else's change is answered")
	}

	stated := strings.Join(rules, "\n")
	for _, named := range []dialogue.Resolution{dialogue.Forwarded, dialogue.Dropped, dialogue.Rejected} {
		if !strings.Contains(stated, string(named)) {
			t.Errorf("the rules do not name %q, which is one of the three answers: %s", named, stated)
		}
	}
	for _, absent := range []dialogue.Resolution{dialogue.Fixed, dialogue.Decided} {
		if strings.Contains(stated, "--"+string(absent)) {
			t.Errorf("the rules name --%s, and nothing here can be patched: %s", absent, stated)
		}
	}
	if strings.Join(dialogue.Rules(), "\n") == stated {
		t.Error("the rules for a change that is not mine are the rules for my own change")
	}
}
