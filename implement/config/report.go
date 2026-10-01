package config

import "fmt"

type Report struct {
	Path     string
	Settings []Setting
}

func ReportAt(path string, loaded Config) Report {
	return Report{Path: path, Settings: loaded.Settings()}
}

func (r Report) Lines() []string {
	keyColumn, valueColumn := 0, 0
	for _, setting := range r.Settings {
		keyColumn = max(keyColumn, len(setting.Key))
		valueColumn = max(valueColumn, len(setting.Value))
	}

	lines := []string{r.Path, ""}
	for _, setting := range r.Settings {
		lines = append(lines, fmt.Sprintf(settingLine, keyColumn, setting.Key, valueColumn, setting.Value, setting.Source))
	}
	return lines
}
