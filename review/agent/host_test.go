package agent

import (
	"testing"

	"github.com/socrates-tools/socrates/internal/skill"
)

func TestEveryHostHasASupportedVendorOfItsOwn(t *testing.T) {
	seen := map[string]string{}
	for _, host := range skill.Hosts() {
		vendor, found := ForHost(host)
		if !found {
			t.Fatalf("%s maps to no vendor", host.Name)
		}
		if !vendor.Qualifies() || !Implemented(vendor.Name) {
			t.Errorf("%s maps to %s, which socrates review does not drive", host.Name, vendor.Name)
		}
		if other, taken := seen[vendor.Name]; taken {
			t.Errorf("%s and %s map to the same vendor %s", host.Name, other, vendor.Name)
		}
		seen[vendor.Name] = host.Name
	}
}

func TestClaudeCodeMapsToTheDefaultVendorAndCodexToCodex(t *testing.T) {
	if vendor, _ := ForHost(skill.ClaudeCode); vendor.Name != DefaultVendor {
		t.Errorf("Claude Code maps to %q, want the default %q", vendor.Name, DefaultVendor)
	}
	codex, _ := Named("codex")
	if vendor, _ := ForHost(skill.Codex); vendor.Name != codex.Name {
		t.Errorf("Codex maps to %q, want %q", vendor.Name, codex.Name)
	}
}

func TestAHostThatIsNotStatedMapsToNoVendor(t *testing.T) {
	if vendor, found := ForHost(skill.Host{Flag: "nobody"}); found {
		t.Errorf("an unstated host maps to %q", vendor.Name)
	}
}
