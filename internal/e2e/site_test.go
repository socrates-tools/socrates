package e2e_test

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/source"
)

const fingerprintLength = 12

var assetsThePageFingerprints = []string{"site.css", "site.js"}

func TestTheAssembledPageAsksForEachAssetByWhatItHolds(t *testing.T) {
	site := assembledSite(t, requireMake(t))

	page, err := os.ReadFile(filepath.Join(site, "index.html"))
	if err != nil {
		t.Fatalf("reading the assembled page: %v", err)
	}

	for _, asset := range assetsThePageFingerprints {
		held, err := os.ReadFile(filepath.Join(site, asset))
		if err != nil {
			t.Fatalf("reading the assembled %s: %v", asset, err)
		}

		want := asset + "?v=" + blobName(held)[:fingerprintLength]
		if !bytes.Contains(page, []byte(`"`+want+`"`)) {
			t.Errorf("the assembled page does not ask for %s, so a browser holding the last release's %s keeps painting with it", want, asset)
		}
		if regexp.MustCompile(`"` + regexp.QuoteMeta(asset) + `"`).Match(page) {
			t.Errorf("the assembled page still asks for %s by its bare name", asset)
		}
	}
}

func assembledSite(t *testing.T, root string) string {
	t.Helper()

	site := filepath.Join(t.TempDir(), "site")
	command := exec.Command("make", "site", "SITE="+site)
	command.Dir = root
	command.Env = append(os.Environ(), source.Toolchain()...)

	var failures bytes.Buffer
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("make site: %v: %s", err, failures.String())
	}
	return site
}

func blobName(held []byte) string {
	named := sha1.New()
	fmt.Fprintf(named, "blob %d\x00", len(held))
	named.Write(held)
	return hex.EncodeToString(named.Sum(nil))
}
