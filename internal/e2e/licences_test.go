package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func theTermsTheRepositoryHolds(t *testing.T) string {
	t.Helper()

	var held strings.Builder
	for _, named := range []string{convention.LicenceFile, convention.NoticeFile} {
		content, err := os.ReadFile(filepath.Join(source.Root(t), named))
		if err != nil {
			t.Fatalf("reading %s: %v", named, err)
		}
		held.Write(content)
	}
	return held.String()
}

func TestTheBinaryPrintsTheLicenceAndTheNoticeItCarries(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	want := theTermsTheRepositoryHolds(t)

	for _, spelled := range []string{"licences", "licenses"} {
		printed, failures, code := session.Top(spelled)
		if code != 0 {
			t.Fatalf("socrates %s exited %d: %s", spelled, code, failures)
		}
		if printed != want {
			t.Errorf("socrates %s does not print %s and then %s as the repository holds them", spelled, convention.LicenceFile, convention.NoticeFile)
		}
	}
}

func TestPrintingTheTermsRefusesAnArgument(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	printed, failures, code := session.Top("licences", "NOTICE")
	if code == 0 {
		t.Fatalf("socrates licences NOTICE exited 0:\n%s", printed)
	}
	if !strings.Contains(failures, `"NOTICE"`) {
		t.Errorf("socrates licences NOTICE does not name the argument it refused:\n%s", failures)
	}
}
