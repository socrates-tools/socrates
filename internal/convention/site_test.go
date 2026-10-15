package convention_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestTheLandingPageInstallsWithTheCommandTheReadmeStates(t *testing.T) {
	unstated, err := convention.InstallCommandsTheLandingPageOmits(source.Root(t))
	if err != nil {
		t.Fatalf("reading the install commands: %v", err)
	}

	for _, command := range unstated {
		t.Errorf("%s installs with %q and %s does not, so a person reads two ways to install", convention.FrontDocumentFile, command, convention.LandingPage)
	}
}

func TestALandingPageInstallingOtherwiseThanTheReadmeIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.FrontDocumentFile), "# x\n\n```sh\ncurl -fsSL https://github.com/socrates-tools/socrates/releases/latest/download/install.sh | sh\n```\n")
	writeNested(t, filepath.Join(root, filepath.FromSlash(convention.LandingPage)), "<pre><span>curl</span> -fsSL https://github.com/socrates-tools/socrates/releases/latest/download/setup.sh <span>|</span> sh</pre>\n")

	unstated, err := convention.InstallCommandsTheLandingPageOmits(root)
	if err != nil {
		t.Fatalf("reading fabricated install commands: %v", err)
	}

	want := []string{"curl -fsSL https://github.com/socrates-tools/socrates/releases/latest/download/install.sh | sh"}
	if diff := cmp.Diff(want, unstated); diff != "" {
		t.Errorf("what a landing page installing otherwise reports differs (-want +got):\n%s", diff)
	}
}

func TestTheReadmesInstallCommandDrawnAcrossMarkupIsStated(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.FrontDocumentFile), "# x\n\n```sh\ncurl -fsSL https://github.com/socrates-tools/socrates/releases/latest/download/install.sh | sh\n```\n")
	writeNested(t, filepath.Join(root, filepath.FromSlash(convention.LandingPage)), "<pre><span class=\"call\">curl</span> <span>-fsSL</span> https://github.com/socrates-tools/socrates/releases/latest/download/install.sh <span>|</span> <span>sh</span></pre>\n")

	unstated, err := convention.InstallCommandsTheLandingPageOmits(root)
	if err != nil {
		t.Fatalf("reading fabricated install commands: %v", err)
	}
	if len(unstated) != 0 {
		t.Errorf("a landing page drawing the command across highlighting was read as omitting %v", unstated)
	}
}
