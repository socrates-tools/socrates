package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/config"
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/state"
)

func write(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), state.ConfigurationFile)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
	return path
}

func TestTheFileIsAgorasOwnUnderTheBrandsRoot(t *testing.T) {
	stated := t.TempDir()
	t.Setenv(state.HomeVariable, stated)

	path, err := config.Path()
	if err != nil {
		t.Fatalf("resolving the configuration path: %v", err)
	}
	if want := filepath.Join(stated, "agora", state.ConfigurationFile); path != want {
		t.Errorf("agora reads %s, and its own file under the brand's root is %s", path, want)
	}
}

func TestNoFileIsNothingStated(t *testing.T) {
	directory := t.TempDir()

	loaded, err := config.LoadFrom(filepath.Join(directory, "agora", state.ConfigurationFile))
	if err != nil {
		t.Fatalf("loading an absent configuration file: %v", err)
	}
	if diff := cmp.Diff(config.Config{}, loaded); diff != "" {
		t.Errorf("an absent configuration file states something (-want +got):\n%s", diff)
	}
	if entries, _ := os.ReadDir(directory); len(entries) != 0 {
		t.Errorf("loading an absent configuration file created %v", entries)
	}
}

func TestBothFieldsOfTheSurfaceBlockAreRead(t *testing.T) {
	loaded, err := config.LoadFrom(write(t, "surface:\n  bind: 127.0.0.1:7788\n  reached_at: https://agora.example:8443/\n"))
	if err != nil {
		t.Fatalf("loading a surface block: %v", err)
	}

	want := config.Config{Surface: config.Surface{Bind: "127.0.0.1:7788", ReachedAt: "https://agora.example:8443/"}}
	if diff := cmp.Diff(want, loaded); diff != "" {
		t.Errorf("the surface block differs (-want +got):\n%s", diff)
	}
}

func TestAFileThatWillNotParseIsARefusalNamingIt(t *testing.T) {
	path := write(t, "surface: [not a mapping\n")

	_, err := config.LoadFrom(path)
	if !refusal.Is(err) {
		t.Fatalf("a file that does not parse gave %v, want a refusal", err)
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("the refusal does not name the file: %s", err)
	}
}

func TestAFieldTheSurfaceBlockDoesNotHaveIsRefused(t *testing.T) {
	_, err := config.LoadFrom(write(t, "surface:\n  proxy: something.example\n"))
	if !refusal.Is(err) {
		t.Fatalf("a surface block naming a field it does not have gave %v, want a refusal", err)
	}
	if !strings.Contains(err.Error(), "proxy") {
		t.Errorf("the refusal does not name the field: %s", err)
	}
}

func TestAFieldOfTheSurfaceBlockStatedOutsideItIsRefused(t *testing.T) {
	for _, stated := range []string{"bind: 127.0.0.1:0\n", "reached_at: https://agora.example:8443/\n"} {
		named, _, _ := strings.Cut(stated, ":")

		_, err := config.LoadFrom(write(t, stated))
		if !refusal.Is(err) {
			t.Fatalf("a file stating %s outside the surface block gave %v, want a refusal", named, err)
		}
		if !strings.Contains(err.Error(), named) {
			t.Errorf("the refusal does not name the field: %s", err)
		}
	}
}

func TestAnUnreadableFileIsARefusalNamingIt(t *testing.T) {
	path := filepath.Join(t.TempDir(), state.ConfigurationFile)
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatalf("making a directory where the file would be: %v", err)
	}

	_, err := config.LoadFrom(path)
	if !refusal.Is(err) {
		t.Fatalf("a file that cannot be read gave %v, want a refusal", err)
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("the refusal does not name the file: %s", err)
	}
}
