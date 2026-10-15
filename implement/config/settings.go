package config

import (
	"slices"
	"strconv"
)

type Source string

const (
	Set     Source = "set"
	Default Source = "default"
)

type Setting struct {
	Key      string
	Path     []string
	Value    string
	Source   Source
	Default  string
	Choices  []string
	Typeable bool
}

var modelAliases = []string{"sonnet", "opus", "haiku", "fable", "inherit"}

func ModelAliases() []string {
	return slices.Clone(modelAliases)
}

const (
	concurrencyKey   = "concurrency"
	modelKey         = "roles.implementer.model"
	effortKey        = "roles.implementer.effort"
	codexModelKey    = "roles.implementer.codex.model"
	codexEffortKey   = "roles.implementer.codex.effort"
	defaultBranchKey = "default_branch"
)

func (c Config) Settings() []Setting {
	stated := ""
	if c.Concurrency != nil {
		stated = strconv.Itoa(int(*c.Concurrency))
	}

	effort := inForce(effortKey, c.Roles.Implementer.Effort, DefaultEffort)
	effort.Choices = Efforts()

	model := inForce(modelKey, c.Roles.Implementer.Model, DefaultModel)
	model.Choices = ModelAliases()
	model.Typeable = true

	codexEffort := inForce(codexEffortKey, c.Roles.Implementer.Codex.Effort, DefaultCodexEffort)
	codexEffort.Choices = CodexEfforts()

	codexModel := inForce(codexModelKey, c.Roles.Implementer.Codex.Model, "")
	codexModel.Typeable = true

	return []Setting{
		inForce(concurrencyKey, stated, strconv.Itoa(DefaultConcurrency)),
		model,
		effort,
		codexModel,
		codexEffort,
	}
}

func (c Config) SettingsIn(repository string) []Setting {
	settings := c.Settings()
	if repository == "" {
		return settings
	}
	row := Setting{Key: defaultBranchKey, Path: []string{defaultBranchesKey, repository}, Source: Default}
	if remembered := c.DefaultBranchOf(repository); remembered != "" {
		row.Value, row.Source = remembered, Set
	}
	return append(settings, row)
}

func inForce(key, stated, owned string) Setting {
	if stated != "" {
		return Setting{Key: key, Value: stated, Source: Set, Default: owned}
	}
	return Setting{Key: key, Value: owned, Source: Default, Default: owned}
}

func (c Config) Implementer() Role {
	role := c.Roles.Implementer
	if role.Model == "" {
		role.Model = DefaultModel
	}
	if role.Effort == "" {
		role.Effort = DefaultEffort
	}
	return role
}

func (c Config) CodexImplementer() CodexRole {
	role := c.Roles.Implementer.Codex
	if role.Effort == "" {
		role.Effort = DefaultCodexEffort
	}
	return role
}
