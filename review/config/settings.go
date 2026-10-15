package config

import (
	"strconv"
	"time"

	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/guidance"
	"github.com/socrates-tools/socrates/review/prompt"
	"github.com/socrates-tools/socrates/review/runlog"
	"github.com/socrates-tools/socrates/review/runstore"
)

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
	reviewerCLIKey     = "reviewer_cli"
	reviewerModelKey   = "roles.reviewer.model"
	reviewerEffortKey  = "roles.reviewer.effort"
	diffCeilingKey     = "diff_ceiling"
	guidanceCeilingKey = "guidance_ceiling"
	agentDeadlineKey   = "agent_deadline"
	logCeilingKey      = "log_ceiling"
	runRetentionKey    = "run_retention"
)

func (c Config) Settings() []Setting {
	vendor, _ := agent.Named(c.Vendor())
	routing := vendor.RoutingFor(agent.Reviewer, c.Roles)
	unrouted := vendor.Routes(agent.Reviewer)
	stated := c.Roles[agent.Reviewer]

	return []Setting{
		chosenFrom(inForce(reviewerCLIKey, c.ReviewerCLI, followsTheHost), agent.SupportedNames()),
		routed(reviewerModelKey, routing.Model, unrouted.Model, stated.Model != ""),
		routed(reviewerEffortKey, routing.Effort, unrouted.Effort, stated.Effort != ""),
		inForce(diffCeilingKey, counted(int64(c.DiffCeiling)), counted(prompt.DefaultCeiling)),
		inForce(guidanceCeilingKey, counted(int64(c.GuidanceCeiling)), counted(guidance.DefaultCeiling)),
		inForce(agentDeadlineKey, lasting(c.AgentDeadline), lasting(agent.DefaultDeadline)),
		inForce(logCeilingKey, counted(c.LogCeiling), counted(runlog.DefaultCeiling)),
		inForce(runRetentionKey, lasting(c.RunRetention), lasting(runstore.DefaultRetention)),
	}
}

func inForce(key, stated, owned string) Setting {
	if stated != "" {
		return Setting{Key: key, Value: stated, Source: Set, Default: owned}
	}
	return Setting{Key: key, Value: owned, Source: Default, Default: owned}
}

func routed(key, resolved, owned string, stated bool) Setting {
	source := Default
	if stated {
		source = Set
	}
	return Setting{Key: key, Value: resolved, Source: source, Default: owned}
}

func chosenFrom(setting Setting, choices []string) Setting {
	setting.Choices = choices
	return setting
}

func counted(measure int64) string {
	if measure <= 0 {
		return ""
	}
	return strconv.FormatInt(measure, 10)
}

func lasting(span time.Duration) string {
	switch {
	case span <= 0:
		return ""
	case span%time.Hour == 0:
		return strconv.FormatInt(int64(span/time.Hour), 10) + "h"
	case span%time.Minute == 0:
		return strconv.FormatInt(int64(span/time.Minute), 10) + "m"
	case span%time.Second == 0:
		return strconv.FormatInt(int64(span/time.Second), 10) + "s"
	}
	return span.String()
}
