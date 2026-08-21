package corpus

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/rui-ktei/socrates/internal/findings"
)

const ModulePath = "github.com/rui-ktei/socrates"

const Directory = "internal/corpus/answers"

const (
	answersDirectory = "answers"
	documentFile     = "answer.json"
	circumstances    = "entry.json"
	moduleDirective  = "module"
	moduleFile       = "go.mod"
)

//go:embed answers
var recorded embed.FS

type Entry struct {
	Name       string `json:"-"`
	Model      string `json:"model"`
	Effort     string `json:"effort"`
	Round      int    `json:"round"`
	Base       string `json:"base"`
	Snapshot   string `json:"snapshot"`
	Run        string `json:"run"`
	RecordedOn string `json:"recorded_on"`
	Subject    string `json:"subject"`
	Document   []byte `json:"-"`
}

func Entries() ([]Entry, error) {
	return read(recorded, answersDirectory)
}

func Named(name string) (Entry, error) {
	entries, err := Entries()
	if err != nil {
		return Entry{}, err
	}

	for _, entry := range entries {
		if entry.Name == name {
			return entry, nil
		}
	}

	held := make([]string, 0, len(entries))
	for _, entry := range entries {
		held = append(held, entry.Name)
	}
	return Entry{}, fmt.Errorf(noSuchAnswer, name, strings.Join(held, ", "))
}

func Read(directory string) ([]Entry, error) {
	return read(os.DirFS(directory), ".")
}

func Record(root, repository string, entry Entry) (string, error) {
	if err := entry.Validate(); err != nil {
		return "", err
	}
	if err := socratesOwnRepository(repository); err != nil {
		return "", err
	}

	at := filepath.Join(root, Directory, entry.Name)
	if _, err := os.Stat(at); err == nil {
		return "", fmt.Errorf(answerAlreadyRecorded, entry.Name, at)
	}
	if err := os.MkdirAll(at, 0o755); err != nil {
		return "", fmt.Errorf(writingAnAnswer, entry.Name, err)
	}

	if err := os.WriteFile(filepath.Join(at, documentFile), entry.Document, 0o644); err != nil {
		return "", fmt.Errorf(writingAnAnswer, entry.Name, err)
	}

	stated, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return "", fmt.Errorf(writingAnAnswer, entry.Name, err)
	}
	if err := os.WriteFile(filepath.Join(at, circumstances), append(stated, '\n'), 0o644); err != nil {
		return "", fmt.Errorf(writingAnAnswer, entry.Name, err)
	}
	return at, nil
}

func (e Entry) Validate() error {
	var absent []string
	for name, stated := range map[string]bool{
		"name":        e.Name == "",
		"model":       e.Model == "",
		"effort":      e.Effort == "",
		"round":       e.Round < 1,
		"base":        e.Base == "",
		"snapshot":    e.Snapshot == "",
		"run":         e.Run == "",
		"recorded_on": e.RecordedOn == "",
		"subject":     e.Subject == "",
		"document":    len(e.Document) == 0,
	} {
		if stated {
			absent = append(absent, name)
		}
	}
	slices.Sort(absent)

	if len(absent) > 0 {
		return fmt.Errorf(entryStatesNothingFor, e.Name, strings.Join(absent, ", "))
	}
	if e.Name != filepath.Base(e.Name) || strings.ContainsAny(e.Name, `/\`) {
		return fmt.Errorf(nameIsAPath, e.Name)
	}
	if strings.HasPrefix(e.Name, ".") || strings.HasPrefix(e.Name, "_") {
		return fmt.Errorf(nameTheCorpusWouldNotSee, e.Name)
	}
	if _, err := findings.Parse(e.Document); err != nil {
		return fmt.Errorf(documentSocratesCannotRoute, e.Name, err)
	}
	return nil
}

func read(fsys fs.FS, at string) ([]Entry, error) {
	held, err := fs.ReadDir(fsys, at)
	if err != nil {
		return nil, fmt.Errorf(readingTheCorpus, at, err)
	}

	var entries []Entry
	for _, held := range held {
		if !held.IsDir() {
			continue
		}

		entry, err := one(fsys, path(at, held.Name()), held.Name())
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	slices.SortFunc(entries, func(one, other Entry) int {
		return strings.Compare(one.Name, other.Name)
	})
	return entries, nil
}

func one(fsys fs.FS, at, name string) (Entry, error) {
	stated, err := fs.ReadFile(fsys, path(at, circumstances))
	if err != nil {
		return Entry{}, fmt.Errorf(readingAnAnswer, name, err)
	}

	var entry Entry
	decoder := json.NewDecoder(strings.NewReader(string(stated)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&entry); err != nil {
		return Entry{}, fmt.Errorf(readingAnAnswer, name, err)
	}

	document, err := fs.ReadFile(fsys, path(at, documentFile))
	if err != nil {
		return Entry{}, fmt.Errorf(readingAnAnswer, name, err)
	}

	entry.Name = name
	entry.Document = document
	return entry, nil
}

func path(at, name string) string {
	if at == "." {
		return name
	}
	return at + "/" + name
}

func socratesOwnRepository(repository string) error {
	stated, err := os.ReadFile(filepath.Join(repository, moduleFile))
	if err != nil {
		return fmt.Errorf(notSocrates, repository, err)
	}

	for line := range strings.SplitSeq(string(stated), "\n") {
		named, declared := strings.CutPrefix(strings.TrimSpace(line), moduleDirective+" ")
		if !declared {
			continue
		}
		if strings.TrimSpace(named) == ModulePath {
			return nil
		}
		return fmt.Errorf(anotherRepository, repository, strings.TrimSpace(named))
	}
	return fmt.Errorf(notSocrates, repository, fs.ErrNotExist)
}
