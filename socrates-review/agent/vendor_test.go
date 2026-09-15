package agent_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/socrates-review/agent"
)

func TestTheSupportedSetIsStatedOnce(t *testing.T) {
	if len(agent.Supported()) == 0 {
		t.Fatal("the supported set names no reviewer CLI, so nothing that reads it is being checked against anything")
	}

	for _, one := range agent.Supported() {
		if one.Binary == "" || one.Tier == "" {
			t.Errorf("%s is supported while %+v, and a supported entry states the binary it runs and the tier its document arrives on", one.Name, one)
		}
	}

	copies, err := convention.VendorsNamedOutsideTheStatement(source.Root(t), agent.Names())
	if err != nil {
		t.Fatalf("reading what names a reviewer CLI: %v", err)
	}
	for _, copied := range copies {
		t.Errorf("%s, which is a second copy of a set stated in %s", copied, convention.VendorStatementFile)
	}
}

func TestAnUnsupportedCLIIsStatedWithTheReasonItIsUnsupported(t *testing.T) {
	unsupported := agent.Unsupported()
	if len(unsupported) == 0 {
		t.Fatal("the statement names no unsupported reviewer CLI, so nothing says why a CLI a person runs is not driven")
	}

	for _, one := range unsupported {
		if len(one.Lacking()) == 0 {
			t.Errorf("%s does not qualify and names no criterion it fails", one.Name)
		}
		if strings.TrimSpace(one.Lacks()) == "" {
			t.Errorf("%s reads as lacking nothing, so a refusal naming it says nothing a person can act on", one.Name)
		}
	}
}

func TestQualifyingAndLackingAgree(t *testing.T) {
	for _, one := range agent.Vendors() {
		if one.Qualifies() != (len(one.Lacking()) == 0) {
			t.Errorf("%s qualifies: %t while lacking %v, so the conjunction and what it reports have drifted apart", one.Name, one.Qualifies(), one.Lacking())
		}
	}
}

func TestACLIThatConstrainsNoSchemaStillQualifies(t *testing.T) {
	extracting := agent.Vendor{
		Name:                      "constrains-nothing",
		Binary:                    "constrains-nothing",
		Tier:                      agent.Parsed,
		IsolatedFromTheRepository: true,
		ReportsWhyItFailed:        true,
		ModelSelectablePerRole:    true,
	}

	if !extracting.Qualifies() {
		t.Error("a CLI meeting every criterion does not qualify because its documents are extracted rather than constrained, so the tier is a gate after all")
	}

	onTheParsedTier := false
	for _, one := range agent.Supported() {
		if one.Tier == agent.Parsed {
			onTheParsedTier = true
		}
	}
	if !onTheParsedTier {
		t.Error("no supported CLI is on the tier that extracts its document, so nothing exercises it")
	}
}

func TestSupportFollowsFromTheCriteriaRatherThanFromAList(t *testing.T) {
	one, found := agent.Named(agent.DefaultVendor)
	if !found {
		t.Fatalf("%s is not in the statement", agent.DefaultVendor)
	}
	if !slices.Contains(agent.SupportedNames(), one.Name) {
		t.Fatalf("%s is not supported, so flipping a criterion on it establishes nothing", one.Name)
	}

	one.IsolatedFromTheRepository = false

	if one.Qualifies() {
		t.Error("a CLI that cannot keep the reviewed repository's instructions out of its session still qualifies, so support comes from a list rather than from the criteria")
	}
	if len(one.Lacking()) == 0 {
		t.Error("a CLI with the isolation criterion flipped reports lacking nothing")
	}
}

func TestTheDefaultIsSupported(t *testing.T) {
	if !agent.Supports(agent.DefaultVendor) {
		t.Errorf("the default reviewer CLI %s is not supported, so an absent configuration field is refused", agent.DefaultVendor)
	}
	if agent.TheDefault().Name != agent.DefaultVendor {
		t.Errorf("the default reads as %q rather than %q", agent.TheDefault().Name, agent.DefaultVendor)
	}
}
