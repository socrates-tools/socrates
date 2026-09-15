package simulation_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/simulator"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/socrates-review/agent"
)

func TestASupportedVendorWithNoSimulatorFailsTheSuite(t *testing.T) {
	supported := agent.Supported()
	if len(supported) == 0 {
		t.Fatal("no reviewer CLI is supported, so this check has nothing to look for")
	}

	root := source.Root(t)
	for _, vendor := range supported {
		named := simulator.SimulatorFor(vendor)

		for _, where := range []string{
			filepath.Join(root, "socrates-review", "cmd", named),
			filepath.Join(root, "socrates-review", named),
		} {
			if _, err := os.Stat(where); err != nil {
				t.Errorf("%s is supported and %s is not there: a vendor whose wire contract nothing simulates is a vendor nobody has tested: %v", vendor.Name, where, err)
			}
		}

		if !agent.Implemented(vendor.Name) {
			t.Errorf("%s is supported and socrates carries no implementation of it", vendor.Name)
		}
	}
}

func TestEveryScenarioIsProducibleByEverySimulator(t *testing.T) {
	// scenario: The simulator produces a failure on demand
	for _, vendor := range agent.Supported() {
		t.Run(vendor.Name, func(t *testing.T) {
			installed := simulator.InstallFor(t, vendor)
			if installed.Binary == "" {
				t.Fatalf("no simulator was installed for %s", vendor.Name)
			}
		})
	}
}
