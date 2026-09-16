package report

import (
	"fmt"

	"github.com/socrates-tools/socrates/review/config"
)

type Configuration struct {
	Path     string
	Settings []config.Setting
}

func ConfigurationAt(path string, loaded config.Config) Configuration {
	return Configuration{Path: path, Settings: loaded.Settings()}
}

func (c Configuration) Lines() []string {
	keyColumn, valueColumn := 0, 0
	for _, setting := range c.Settings {
		keyColumn = max(keyColumn, len(setting.Key))
		valueColumn = max(valueColumn, len(setting.Value))
	}

	lines := []string{c.Path, ""}
	for _, setting := range c.Settings {
		lines = append(lines, fmt.Sprintf(settingLine, keyColumn, setting.Key, valueColumn, setting.Value, setting.Source))
	}
	return lines
}
