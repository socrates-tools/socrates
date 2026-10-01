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
	concurrencyKey = "concurrency"
	modelKey       = "roles.implementer.model"
	effortKey      = "roles.implementer.effort"
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
