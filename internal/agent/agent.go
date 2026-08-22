package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"slices"
)

type Role string

const Reviewer Role = "reviewer"

var roles = []Role{Reviewer}

var defaultRouting = map[Role]Routing{
	Reviewer: {Model: "opus", Effort: "high"},
}

type Routing struct {
	Model  string `yaml:"model"`
	Effort string `yaml:"effort"`
}

type Configuration map[Role]Routing

func Roles() []Role {
	return slices.Clone(roles)
}

func (r Role) Known() bool {
	return slices.Contains(roles, r)
}

func Default(role Role) (Routing, bool) {
	routing, found := defaultRouting[role]
	return routing, found
}

type Usage struct {
	InputTokens         int `json:"input_tokens"`
	OutputTokens        int `json:"output_tokens"`
	CacheReadTokens     int `json:"cache_read_tokens"`
	CacheCreationTokens int `json:"cache_creation_tokens"`
}

type Request struct {
	Role       Role
	Prompt     string
	Dir        string
	Schema     json.RawMessage
	Honours    func(json.RawMessage) error
	Transcript io.Writer
}

type Response struct {
	Structured  json.RawMessage
	Text        string
	Model       string
	Effort      string
	SessionID   string
	Usage       Usage
	Opened      bool
	ToolServers []string
	Tools       []string
}

type Agent interface {
	Run(context.Context, Request) (Response, error)
}

func resolveRouting(configured Configuration) (map[Role]Routing, error) {
	for role := range configured {
		if !role.Known() {
			return nil, fmt.Errorf(notARole, role, roles)
		}
	}

	resolved := make(map[Role]Routing, len(roles))
	for _, role := range roles {
		routing := defaultRouting[role]
		if stated, found := configured[role]; found {
			if stated.Model != "" {
				routing.Model = stated.Model
			}
			if stated.Effort != "" {
				routing.Effort = stated.Effort
			}
		}
		resolved[role] = routing
	}
	return resolved, nil
}
