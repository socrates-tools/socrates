package config

import "strconv"

type Source string

const (
	Set     Source = "set"
	Default Source = "default"
)

type Setting struct {
	Key     string
	Value   string
	Source  Source
	Default string
	Choices []string
}

const (
	concurrencyKey   = "concurrency"
	modelKey         = "roles.implementer.model"
	effortKey        = "roles.implementer.effort"
	defaultBranchKey = "default_branch"
)

func (c Config) Settings() []Setting {
	stated := ""
	if c.Concurrency != nil {
		stated = strconv.Itoa(int(*c.Concurrency))
	}

	effort := inForce(effortKey, c.Roles.Implementer.Effort, DefaultEffort)
	effort.Choices = Efforts()

	return []Setting{
		inForce(concurrencyKey, stated, strconv.Itoa(DefaultConcurrency)),
		inForce(modelKey, c.Roles.Implementer.Model, DefaultModel),
		effort,
	}
}

func (c Config) SettingsIn(repository string) []Setting {
	settings := c.Settings()
	if repository == "" {
		return settings
	}
	remembered := c.DefaultBranchOf(repository)
	if remembered == "" {
		return append(settings, Setting{Key: defaultBranchKey, Source: Default})
	}
	return append(settings, Setting{Key: defaultBranchKey, Value: remembered, Source: Set})
}

func inForce(key, stated, owned string) Setting {
	if stated != "" {
		return Setting{Key: key, Value: stated, Source: Set, Default: owned}
	}
	return Setting{Key: key, Value: owned, Source: Default, Default: owned}
}
