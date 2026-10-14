package agent

import (
	"fmt"
	"slices"
	"strings"
)

type Tier string

const (
	Enforced Tier = "enforced"
	Parsed   Tier = "parsed"
)

type Vendor struct {
	Name                      string
	Binary                    string
	Tier                      Tier
	Routing                   Configuration
	IsolatedFromTheRepository bool
	ReportsWhyItFailed        bool
	ModelSelectablePerRole    bool
	ReportsItsToolInventory   bool
}

func (v Vendor) Qualifies() bool {
	return v.IsolatedFromTheRepository &&
		v.ReportsWhyItFailed &&
		v.ModelSelectablePerRole
}

func (v Vendor) Lacking() []string {
	var lacking []string
	if !v.IsolatedFromTheRepository {
		lacking = append(lacking, isolationFromTheRepository)
	}
	if !v.ReportsWhyItFailed {
		lacking = append(lacking, anAccountOfWhyItFailed)
	}
	if !v.ModelSelectablePerRole {
		lacking = append(lacking, aModelPerRole)
	}
	return lacking
}

func (v Vendor) Lacks() string {
	return strings.Join(v.Lacking(), and)
}

var stated = []Vendor{
	{
		Name:                      "claude",
		Binary:                    "claude",
		Tier:                      Enforced,
		Routing:                   Configuration{Reviewer: {Model: "opus", Effort: "high"}},
		IsolatedFromTheRepository: true,
		ReportsWhyItFailed:        true,
		ModelSelectablePerRole:    true,
		ReportsItsToolInventory:   true,
	},
	{
		Name:                      "codex",
		Binary:                    "codex",
		Tier:                      Enforced,
		Routing:                   Configuration{Reviewer: {Effort: "high"}},
		IsolatedFromTheRepository: true,
		ReportsWhyItFailed:        true,
		ModelSelectablePerRole:    true,
		ReportsItsToolInventory:   false,
	},
	{
		Name:                      "opencode",
		Binary:                    "opencode",
		Tier:                      Parsed,
		IsolatedFromTheRepository: true,
		ReportsWhyItFailed:        true,
		ModelSelectablePerRole:    true,
		ReportsItsToolInventory:   false,
	},
	{
		Name:                      "cursor-agent",
		Binary:                    "cursor-agent",
		IsolatedFromTheRepository: false,
		ReportsWhyItFailed:        true,
		ModelSelectablePerRole:    true,
		ReportsItsToolInventory:   false,
	},
	{
		Name:                      "copilot",
		Binary:                    "copilot",
		IsolatedFromTheRepository: false,
		ReportsWhyItFailed:        false,
		ModelSelectablePerRole:    true,
		ReportsItsToolInventory:   false,
	},
	{
		Name:                      "pi",
		Binary:                    "pi",
		IsolatedFromTheRepository: false,
		ReportsWhyItFailed:        false,
		ModelSelectablePerRole:    true,
		ReportsItsToolInventory:   false,
	},
}

const DefaultVendor = "claude"

func Vendors() []Vendor {
	return slices.Clone(stated)
}

func Supported() []Vendor {
	supported := make([]Vendor, 0, len(stated))
	for _, one := range stated {
		if one.Qualifies() {
			supported = append(supported, one)
		}
	}
	return supported
}

func Unsupported() []Vendor {
	var unsupported []Vendor
	for _, one := range stated {
		if !one.Qualifies() {
			unsupported = append(unsupported, one)
		}
	}
	return unsupported
}

func SupportedNames() []string {
	var named []string
	for _, one := range Supported() {
		named = append(named, one.Name)
	}
	return named
}

func Names() []string {
	var named []string
	for _, one := range stated {
		named = append(named, one.Name)
	}
	return named
}

func Stated() string {
	return strings.Join(SupportedNames(), between)
}

func Named(name string) (Vendor, bool) {
	for _, one := range stated {
		if one.Name == name {
			return one, true
		}
	}
	return Vendor{}, false
}

func Supports(name string) bool {
	one, found := Named(name)
	return found && one.Qualifies()
}

func TheDefault() Vendor {
	one, _ := Named(DefaultVendor)
	return one
}

var implementations = map[string]func(Vendor, Options) (Agent, error){
	"claude":   newClaude,
	"codex":    newCodex,
	"opencode": newOpencode,
}

func Implemented(name string) bool {
	_, built := implementations[name]
	return built
}

func New(named string, options Options) (Agent, error) {
	if named == "" {
		named = DefaultVendor
	}

	vendor, known := Named(named)
	if !known || !vendor.Qualifies() {
		return nil, fmt.Errorf(notASupportedReviewerCLI, named, Stated())
	}

	build, built := implementations[vendor.Name]
	if !built {
		return nil, fmt.Errorf(noImplementation, named)
	}
	return build(vendor, options)
}
