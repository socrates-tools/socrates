package configure

import (
	"context"
	"io"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/settingsmenu"
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

	return settingsmenu.Run(ctx, in, out, definitionAt(path))
}

func definitionAt(path string) settingsmenu.Definition {
	return settingsmenu.Definition{
		Path:     path,
		Parse:    config.Parse,
		Settings: settingsIn,
		Wording: settingsmenu.Wording{
			CannotBeRead:       cannotBeRead,
			CannotBeWritten:    cannotBeWritten,
			CannotBeUnderstood: cannotBeUnderstood,
			ReadFromNowOn:      readFromNowOn,
		},
	}
}

func settingsIn(path string) ([]settingsmenu.Setting, error) {
	loaded, err := config.LoadFrom(path)
	if err != nil {
		return nil, err
	}

	var offered []settingsmenu.Setting
	for _, setting := range loaded.Settings() {
		offered = append(offered, settingsmenu.Setting{
			Key:     setting.Key,
			Value:   setting.Value,
			Source:  string(setting.Source),
			Default: setting.Default,
			Choices: setting.Choices,
		})
	}
	return offered, nil
}
