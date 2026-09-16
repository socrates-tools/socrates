package configure

import (
	"context"
	"io"

	tea "charm.land/bubbletea/v2"
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/review/config"
	"github.com/socrates-tools/socrates/review/terminal"
)

func Run(ctx context.Context, in io.Reader, out io.Writer) error {
	if !terminal.Attached(in) || !terminal.Attached(out) {
		return refusal.New(noPersonAtAKeyboard, readingCommand)
	}

	path, err := config.Path()
	if err != nil {
		return err
	}

	standing, err := opened(path)
	if err != nil {
		return err
	}

	program := tea.NewProgram(standing, tea.WithContext(ctx), tea.WithInput(in), tea.WithOutput(out))
	if _, err := program.Run(); err != nil {
		return err
	}
	return standing.failed
}
