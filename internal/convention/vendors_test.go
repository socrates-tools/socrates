package convention_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/review/agent"
)

func TestADocumentThatNamesFewerCLIsThanSocratesSupportsFails(t *testing.T) {
	root := source.Root(t)

	documented, err := convention.VendorsDocumented(root)
	if err != nil {
		t.Fatalf("reading what a person is told about the reviewer CLIs: %v", err)
	}
	if len(documented) == 0 {
		t.Fatalf("%q names no reviewer CLI, so a reader learns nothing about which are driven", convention.VendorsHeading)
	}

	named := make([]string, 0, len(documented))
	for _, row := range documented {
		named = append(named, strings.Trim(strings.TrimSpace(row), "`"))
	}

	supported := agent.SupportedNames()
	for _, disagreement := range convention.VendorDisagreements(supported, onlyThoseSupported(named, supported), convention.VendorsHeading) {
		t.Errorf("%s", disagreement)
	}
}

func TestAReaderLearnsWhyTheirCLIIsNotSupported(t *testing.T) {
	root := source.Root(t)

	section, err := convention.VendorSection(root)
	if err != nil {
		t.Fatalf("reading what a person is told about the reviewer CLIs: %v", err)
	}
	if strings.TrimSpace(section) == "" {
		t.Fatalf("%q is not there", convention.VendorsHeading)
	}

	unsupported := agent.Unsupported()
	if len(unsupported) == 0 {
		t.Fatal("every stated reviewer CLI is supported, so this check reads nothing")
	}

	for _, one := range unsupported {
		if !strings.Contains(section, one.Name) {
			t.Errorf("%s does not name %s, so a person who runs it learns nothing", convention.VendorsHeading, one.Name)
		}
		for _, lacking := range one.Lacking() {
			if !strings.Contains(section, lacking) {
				t.Errorf("%s names %s without saying it lacks %q", convention.VendorsHeading, one.Name, lacking)
			}
		}
	}

	if !strings.Contains(section, agent.UnsupportedAreUnplanned) {
		t.Errorf("%s names CLIs that are not supported without saying whether that is planned: %q", convention.VendorsHeading, agent.UnsupportedAreUnplanned)
	}
}

func TestEveryTierIsStatedForAReaderChoosingBetweenThem(t *testing.T) {
	root := source.Root(t)

	section, err := convention.VendorSection(root)
	if err != nil {
		t.Fatalf("reading what a person is told about the reviewer CLIs: %v", err)
	}

	for _, one := range agent.Supported() {
		row := rowFor(section, one.Name)
		if row == "" {
			t.Errorf("%s is supported and %s carries no row for it", one.Name, convention.VendorsHeading)
			continue
		}

		constrains := strings.Contains(row, "constrains")
		extracts := strings.Contains(row, "extracts")

		if one.Tier == agent.Enforced && !constrains {
			t.Errorf("%s is on the enforced tier and its row reads %q", one.Name, row)
		}
		if one.Tier == agent.Parsed && !extracts {
			t.Errorf("%s is on the extracting tier and its row reads %q", one.Name, row)
		}
	}
}

func onlyThoseSupported(named, supported []string) []string {
	var kept []string
	for _, one := range named {
		for _, supports := range supported {
			if one == supports {
				kept = append(kept, one)
			}
		}
	}
	return kept
}

func rowFor(section, name string) string {
	for _, line := range strings.Split(section, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "| `"+name+"`") {
			return line
		}
	}
	return ""
}

func TestTheSecurityPolicyNamesEveryReviewerCLIThatCarriesTheChangeToItsVendor(t *testing.T) {
	policy, err := convention.Document(source.Root(t), convention.SecurityPolicyFile)
	if err != nil {
		t.Fatalf("reading %s: %v", convention.SecurityPolicyFile, err)
	}
	_, section, found := strings.Cut(policy, "\n## Where your code goes\n")
	if !found {
		t.Fatalf("%s has no section saying where your code goes", convention.SecurityPolicyFile)
	}
	section, _, _ = strings.Cut(section, "\n## ")

	for _, name := range agent.SupportedNames() {
		if !strings.Contains(section, "`"+name+"`") {
			t.Errorf("%s says where your code goes and never names %s, which review drives and which sends the change to its vendor", convention.SecurityPolicyFile, name)
		}
	}
}
