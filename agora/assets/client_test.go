package assets_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/source"
)

const buildScript = "build.mjs"

func TestTheClientHoldsNoJavaScriptButItsBuildScript(t *testing.T) {
	// scenario: The client's source holds no JavaScript file other than the build script
	client := filepath.Join(source.Root(t), "agora", "web")

	for _, held := range within(t, client) {
		relative, err := filepath.Rel(client, held)
		if err != nil {
			t.Fatalf("naming %s: %v", held, err)
		}
		extension := filepath.Ext(held)
		if (extension == ".js" || extension == ".mjs" || extension == ".cjs" || extension == ".jsx") && relative != buildScript {
			t.Errorf("%s is JavaScript, and the client is TypeScript beside its one build script", filepath.ToSlash(relative))
		}
	}
}

func TestTheTypeCheckerRunsInStrictModeOverAllOfTheClient(t *testing.T) {
	// scenario: The type checker runs in strict mode over all of the client
	held, err := os.ReadFile(filepath.Join(source.Root(t), "agora", "web", "tsconfig.json"))
	if err != nil {
		t.Fatalf("reading the type checker's configuration: %v", err)
	}

	var configured struct {
		CompilerOptions struct {
			Strict bool `json:"strict"`
			NoEmit bool `json:"noEmit"`
		} `json:"compilerOptions"`
		Include []string `json:"include"`
	}
	if err := json.Unmarshal(held, &configured); err != nil {
		t.Fatalf("parsing the type checker's configuration: %v", err)
	}

	if !configured.CompilerOptions.Strict {
		t.Error("the type checker is not in strict mode")
	}
	if !configured.CompilerOptions.NoEmit {
		t.Error("the type checker emits files, and the bundle is what esbuild writes")
	}
	if !slices.Contains(configured.Include, "packages/**/*.ts") {
		t.Errorf("the type checker covers %v rather than every TypeScript source under packages", configured.Include)
	}
}
