package settingsmenu

import (
	"context"
	"io"

	tea "charm.land/bubbletea/v2"
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/tty"
)

type Setting struct {
	Key     string
	Value   string
	Source  string
	Default string
	Choices []string
}

type Wording struct {
	CannotBeRead       string
	CannotBeWritten    string
	CannotBeUnderstood string
	ReadFromNowOn      string
	NothingIsATerminal string
}

type Definition struct {
	Path     string
	Parse    func(path string, content []byte) error
	Settings func(path string) ([]Setting, error)
	Wording  Wording
}

func Run(ctx context.Context, in io.Reader, out io.Writer, definition Definition) error {
	if !tty.Attached(in) || !tty.Attached(out) {
		return &refusal.Refusal{Reason: definition.Wording.NothingIsATerminal}
	}

	standing, err := opened(definition)
	if err != nil {
		return err
	}

	program := tea.NewProgram(standing, tea.WithContext(ctx), tea.WithInput(in), tea.WithOutput(out))
	if _, err := program.Run(); err != nil {
		return err
	}
	return standing.failed
}

type Repair struct {
	Key     string
	Line    int
	Written []byte
}

func RepairFirst(definition Definition, content []byte) (Repair, bool, error) {
	broken := definition.unreadableIn(definition.Path, content)
	if len(broken) == 0 {
		return Repair{}, false, nil
	}
	written, err := repaired(content, broken[0])
	if err != nil {
		return Repair{}, false, err
	}
	return Repair{Key: broken[0].key, Line: broken[0].line, Written: written}, true, nil
}
