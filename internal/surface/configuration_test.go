package surface_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rui-ktei/socrates/internal/refusal"
	"github.com/rui-ktei/socrates/internal/surface"
)

func TestAnAbsentConfigurationIsTheZeroValue(t *testing.T) {
	directory := t.TempDir()

	loaded, err := surface.LoadConfigurationFrom(filepath.Join(directory, "agora", "config.yaml"))
	if err != nil {
		t.Fatalf("loading an absent configuration file: %v", err)
	}
	if diff := cmp.Diff(surface.Configuration{}, loaded); diff != "" {
		t.Errorf("an absent configuration file is not the zero value (-want +got):\n%s", diff)
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("reading %s: %v", directory, err)
	}
	if len(entries) != 0 {
		t.Errorf("loading an absent configuration file created %v", entries)
	}
}

func TestBothAddressesAreRead(t *testing.T) {
	path := writeConfiguration(t, `bind: 127.0.0.1:7788
reached_at: https://agora.example:8443/
`)

	loaded, err := surface.LoadConfigurationFrom(path)
	if err != nil {
		t.Fatalf("loading %s: %v", path, err)
	}

	want := surface.Configuration{Bind: "127.0.0.1:7788", ReachedAt: "https://agora.example:8443/"}
	if diff := cmp.Diff(want, loaded); diff != "" {
		t.Errorf("what an operator stated was not read back (-want +got):\n%s", diff)
	}
}

func TestAConfigurationThatWillNotParseIsRefused(t *testing.T) {
	path := writeConfiguration(t, "bind: [127.0.0.1:0\n")

	if _, err := surface.LoadConfigurationFrom(path); !refusal.Is(err) {
		t.Fatalf("a configuration file that will not parse gave %v, want a refusal", err)
	}
}

func TestAConfigurationNamingAFieldAgoraDoesNotKnowIsRefused(t *testing.T) {
	path := writeConfiguration(t, "proxy: something.example\n")

	_, err := surface.LoadConfigurationFrom(path)
	if !refusal.Is(err) {
		t.Fatalf("a configuration file naming an unknown field gave %v, want a refusal", err)
	}
	if !strings.Contains(err.Error(), "proxy") {
		t.Errorf("the refusal does not name the field: %s", err)
	}
}

func TestTheConfigurationSitsAtTheTopOfAgorasOwnHome(t *testing.T) {
	stated := t.TempDir()
	t.Setenv(surface.HomeVariable, stated)

	path, err := surface.ConfigurationPath()
	if err != nil {
		t.Fatalf("resolving where agora reads its configuration: %v", err)
	}
	if want := filepath.Join(stated, "config.yaml"); path != want {
		t.Errorf("agora reads %s, and its home is %s", path, want)
	}

	session, err := surface.SessionPath()
	if err != nil {
		t.Fatalf("resolving where agora writes its session: %v", err)
	}
	if filepath.Dir(session) != stated {
		t.Errorf("agora writes its session at %s, outside the home its configuration was read from", session)
	}
}

func TestARelativeHomeIsRefusedByAgora(t *testing.T) {
	t.Setenv(surface.HomeVariable, filepath.Join("somewhere", "below"))

	for named, resolve := range map[string]func() (string, error){
		"its configuration": surface.ConfigurationPath,
		"its session":       surface.SessionPath,
		"its document":      surface.DocumentPath,
	} {
		path, err := resolve()
		if !refusal.Is(err) {
			t.Errorf("agora resolved %s to %q with a relative home, want a refusal: %v", named, path, err)
			continue
		}
		if !strings.Contains(err.Error(), surface.HomeVariable) || !strings.Contains(err.Error(), "somewhere") {
			t.Errorf("the refusal reads %q, without naming the variable and the path", err)
		}
	}
}

func writeConfiguration(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
	return path
}
