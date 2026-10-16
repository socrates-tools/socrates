package platform

import (
	"encoding/json"
	"runtime"
	"slices"
	"strings"

	"github.com/socrates-tools/socrates/internal/process"
)

const ReleaseTargetsFile = "platforms.json"

type Platform struct {
	OperatingSystem           string
	Name                      string
	Runner                    string
	Architectures             []string
	ReapsByProcessGroup       bool
	FilesystemEnforcesTheSeal bool
	TheSuiteRunsThere         bool
}

type Build struct {
	OperatingSystem string
	Architecture    string
}

type Target struct {
	OS     string `json:"os"`
	Runner string `json:"runner"`
}

func (p Platform) Qualifies() bool {
	return p.ReapsByProcessGroup && p.FilesystemEnforcesTheSeal && p.TheSuiteRunsThere
}

var stated = []Platform{
	{
		OperatingSystem:           "linux",
		Name:                      "Linux",
		Runner:                    "ubuntu-latest",
		Architectures:             []string{"amd64", "arm64"},
		ReapsByProcessGroup:       true,
		FilesystemEnforcesTheSeal: true,
		TheSuiteRunsThere:         true,
	},
	{
		OperatingSystem:           "darwin",
		Name:                      "macOS",
		Runner:                    "macos-latest",
		Architectures:             []string{"amd64", "arm64"},
		ReapsByProcessGroup:       true,
		FilesystemEnforcesTheSeal: true,
		TheSuiteRunsThere:         true,
	},
}

func Supported() []Platform {
	supported := make([]Platform, 0, len(stated))
	for _, one := range stated {
		if one.Qualifies() {
			supported = append(supported, one)
		}
	}
	return supported
}

func OperatingSystems() []string {
	var named []string
	for _, one := range Supported() {
		named = append(named, one.OperatingSystem)
	}
	return named
}

func Architectures() []string {
	var named []string
	for _, one := range Supported() {
		for _, architecture := range one.Architectures {
			if !slices.Contains(named, architecture) {
				named = append(named, architecture)
			}
		}
	}
	slices.Sort(named)
	return named
}

func Builds() []Build {
	var builds []Build
	for _, one := range Supported() {
		for _, architecture := range one.Architectures {
			builds = append(builds, Build{OperatingSystem: one.OperatingSystem, Architecture: architecture})
		}
	}
	return builds
}

func Names() []string {
	var named []string
	for _, one := range Supported() {
		named = append(named, one.Name)
	}
	return named
}

func Runners() []string {
	var running []string
	for _, one := range Supported() {
		running = append(running, one.Runner)
	}
	return running
}

func Targets() []Target {
	var targets []Target
	for _, one := range Supported() {
		targets = append(targets, Target{OS: one.OperatingSystem, Runner: one.Runner})
	}
	return targets
}

func Rendered() ([]byte, error) {
	rendered, err := json.MarshalIndent(Targets(), "", "  ")
	if err != nil {
		return nil, err
	}
	return append(rendered, '\n'), nil
}

func Stated() string {
	return strings.Join(Names(), between)
}

func Supports(operatingSystem string) bool {
	return slices.Contains(OperatingSystems(), operatingSystem)
}

func TheOneRunning() string {
	return runtime.GOOS
}

func ReviewIsSupported() bool {
	return process.ReapsByGroup && Supports(TheOneRunning())
}
