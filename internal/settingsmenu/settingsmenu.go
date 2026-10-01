package settingsmenu

import (
	"context"
	"io"

	tea "charm.land/bubbletea/v2"
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
}

type Definition struct {
	Path     string
	Parse    func(path string, content []byte) error
	Settings func(path string) ([]Setting, error)
	Wording  Wording
}

func Run(ctx context.Context, in io.Reader, out io.Writer, definition Definition) error {
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
