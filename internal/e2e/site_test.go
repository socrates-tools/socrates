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

var assetsThePageFingerprints = []string{"tokens.css", "site.generated.css", "site.js"}

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

func TestTheSiteIsAssembledWithoutNodeOrPnpmOnThePath(t *testing.T) {
	root := requireMake(t)

	refusing := t.TempDir()
	for _, tool := range []string{"node", "pnpm", "npm", "npx", "corepack"} {
		script := "#!/bin/sh\necho \"make site ran " + tool + "\" >&2\nexit 1\n"
		if err := os.WriteFile(filepath.Join(refusing, tool), []byte(script), 0o755); err != nil {
			t.Fatalf("writing a %s that refuses to run: %v", tool, err)
		}
	}

	site := assembledSiteWith(t, root, []string{"PATH=" + refusing + string(os.PathListSeparator) + os.Getenv("PATH")})

	for _, asset := range []string{"index.html", "site.generated.css", "site.js"} {
		if _, err := os.Stat(filepath.Join(site, asset)); err != nil {
			t.Errorf("the site assembled without Node holds no %s: %v", asset, err)
		}
	}
}

func assembledSite(t *testing.T, root string) string {
	t.Helper()

	return assembledSiteWith(t, root, nil)
}

func assembledSiteWith(t *testing.T, root string, environment []string) string {
	t.Helper()

	site := filepath.Join(t.TempDir(), "site")
	command := exec.Command("make", "site", "SITE="+site)
	command.Dir = root
	command.Env = append(append(os.Environ(), source.Toolchain()...), environment...)

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
