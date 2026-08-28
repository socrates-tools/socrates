package convention_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rui-ktei/socrates/internal/convention"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

func TestNoUserFacingStringNamesAVariableTheBinariesDoNotRead(t *testing.T) {
	unread, err := convention.VariablesNamedToAPerson(source.Root(t))
	if err != nil {
		t.Fatalf("reading what the messages name: %v", err)
	}

	for _, named := range unread {
		t.Errorf("%s tells a person to set %s, and nothing in the source reads it", named.File, named.Name)
	}
}

func TestAMessageNamingAVariableNothingReadsIsReported(t *testing.T) {
	root := t.TempDir()
	writePackaged(t, filepath.Join(root, "reader", "reader.go"), "package reader\n\nconst Variable = \"SOCRATES_HOME\"\n")
	writePackaged(t, filepath.Join(root, "reader", "messages.go"), "package reader\n\nconst gone = \"Set SOCRATES_STATE_HOME and run the command again\"\n")

	unread, err := convention.VariablesNamedToAPerson(root)
	if err != nil {
		t.Fatalf("reading what the messages name: %v", err)
	}
	if len(unread) != 1 || unread[0].Name != "SOCRATES_STATE_HOME" {
		t.Fatalf("the check found %v, want the one variable nothing reads", unread)
	}
}

func TestAMessageNamingAVariableTheSourceReadsIsNotReported(t *testing.T) {
	root := t.TempDir()
	writePackaged(t, filepath.Join(root, "reader", "reader.go"), "package reader\n\nconst Variable = \"SOCRATES_HOME\"\n")
	writePackaged(t, filepath.Join(root, "reader", "messages.go"), "package reader\n\nconst held = \"Set SOCRATES_HOME and run the command again\"\n")

	unread, err := convention.VariablesNamedToAPerson(root)
	if err != nil {
		t.Fatalf("reading what the messages name: %v", err)
	}
	if len(unread) != 0 {
		t.Errorf("the check reported %v, and the source reads every variable named", unread)
	}
}

func writePackaged(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}
