package config

import (
	"strconv"
	"time"

	"github.com/socrates-tools/socrates-review/internal/agent"
	"github.com/socrates-tools/socrates-review/internal/guidance"
	"github.com/socrates-tools/socrates-review/internal/loopback"
	"github.com/socrates-tools/socrates-review/internal/prompt"
	"github.com/socrates-tools/socrates-review/internal/runlog"
	"github.com/socrates-tools/socrates-review/internal/runstore"
)

type Source string

const (
	Set     Source = "set"
	Default Source = "default"
)

type Setting struct {
	Key    string
	Value  string
	Source Source
}

const (
	reviewerCLIKey      = "reviewer_cli"
	reviewerModelKey    = "roles.reviewer.model"
	reviewerEffortKey   = "roles.reviewer.effort"
	diffCeilingKey      = "diff_ceiling"
	guidanceCeilingKey  = "guidance_ceiling"
	agentDeadlineKey    = "agent_deadline"
	logCeilingKey       = "log_ceiling"
	runRetentionKey     = "run_retention"
	surfaceBindKey      = "surface.bind"
	surfaceReachedAtKey = "surface.reached_at"
)

func (c Config) Settings() []Setting {
	vendor, _ := agent.Named(c.Vendor())
	routing := vendor.RoutingFor(agent.Reviewer, c.Roles)
	stated := c.Roles[agent.Reviewer]

	return []Setting{
		inForce(reviewerCLIKey, c.ReviewerCLI, agent.DefaultVendor),
		routed(reviewerModelKey, routing.Model, stated.Model != ""),
		routed(reviewerEffortKey, routing.Effort, stated.Effort != ""),
		inForce(diffCeilingKey, counted(int64(c.DiffCeiling)), counted(prompt.DefaultCeiling)),
		inForce(guidanceCeilingKey, counted(int64(c.GuidanceCeiling)), counted(guidance.DefaultCeiling)),
		inForce(agentDeadlineKey, lasting(c.AgentDeadline), lasting(agent.DefaultDeadline)),
		inForce(logCeilingKey, counted(c.LogCeiling), counted(runlog.DefaultCeiling)),
		inForce(runRetentionKey, lasting(c.RunRetention), lasting(runstore.DefaultRetention)),
		inForce(surfaceBindKey, c.Surface.Bind, loopback.DefaultAddress),
		inForce(surfaceReachedAtKey, c.Surface.ReachedAt, ""),
	}
}

func inForce(key, stated, owned string) Setting {
	if stated != "" {
		return Setting{Key: key, Value: stated, Source: Set}
	}
	return Setting{Key: key, Value: owned, Source: Default}
}

func routed(key, resolved string, stated bool) Setting {
	source := Default
	if stated {
		source = Set
	}
	return Setting{Key: key, Value: resolved, Source: source}
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
