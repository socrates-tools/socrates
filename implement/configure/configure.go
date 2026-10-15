package configure

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/implement/skill"
	"github.com/socrates-tools/socrates/internal/settingsmenu"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
)

const (
	defaultBranchKey = "default_branch"
	rolesPrefix      = "roles"
)

func Run(ctx context.Context, in io.Reader, out io.Writer, repository string) error {
	path, err := config.Path()
	if err != nil {
		return err
	}

	return settingsmenu.Run(ctx, in, out, definitionAt(path, repository))
}

func definitionAt(path, repository string) settingsmenu.Definition {
	return settingsmenu.Definition{
		Path:           path,
		Parse:          config.Parse,
		Settings:       settingsIn(repository),
		Tool:           toolName,
		MenuCommand:    menuCommand,
		ReadingCommand: readingCommand,
		ReadFromNowOn:  readFromNowOn,
		Validate:       validatedIn(repository),
		Followed:       followedBy(path),
	}
}

func settingsIn(repository string) func(path string) ([]settingsmenu.Setting, error) {
	return func(path string) ([]settingsmenu.Setting, error) {
		loaded, err := config.LoadFrom(path)
		if err != nil {
			return nil, err
		}

		var offered []settingsmenu.Setting
		for _, setting := range loaded.SettingsIn(repository) {
			offered = append(offered, settingsmenu.Setting{
				Key:      setting.Key,
				Path:     setting.Path,
				Value:    setting.Value,
				Source:   string(setting.Source),
				Default:  setting.Default,
				Choices:  setting.Choices,
				Typeable: setting.Typeable,
			})
		}
		return offered, nil
	}
}

func validatedIn(repository string) func(setting settingsmenu.Setting, value string) error {
	return func(setting settingsmenu.Setting, value string) error {
		if setting.Key != defaultBranchKey {
			return nil
		}
		return config.RequireBranch(repository, value)
	}
}

func followedBy(path string) func(setting settingsmenu.Setting) (string, error) {
	return func(setting settingsmenu.Setting) (string, error) {
		if !shapesTheDefinition(setting.Key) {
			return readFromNowOn, nil
		}

		loaded, err := config.LoadFrom(path)
		if err != nil {
			return "", err
		}
		installed, wasInstalled, err := skill.Reinstall(skillkit.ClaudeCode, loaded)
		if err != nil {
			return "", err
		}
		if !wasInstalled {
			return definitionNotInstalled, nil
		}
		role := loaded.Implementer()
		return fmt.Sprintf(definitionRewritten, installed, role.Model, role.Effort), nil
	}
}

func shapesTheDefinition(key string) bool {
	return key == rolesPrefix || strings.HasPrefix(key, rolesPrefix+".")
}
