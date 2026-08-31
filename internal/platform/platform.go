package platform

import (
	"runtime"
	"slices"
	"strings"

	"github.com/socrates-devtool/socrates/internal/process"
)

const ReleaseTargetsFile = "platforms.json"

type Platform struct {
	OperatingSystem           string
	Name                      string
	ReapsByProcessGroup       bool
	FilesystemEnforcesTheSeal bool
	TheSuiteRunsThere         bool
}

func (p Platform) Qualifies() bool {
	return p.ReapsByProcessGroup && p.FilesystemEnforcesTheSeal && p.TheSuiteRunsThere
}

var stated = []Platform{
	{
		OperatingSystem:           "linux",
		Name:                      "Linux",
		ReapsByProcessGroup:       true,
		FilesystemEnforcesTheSeal: true,
		TheSuiteRunsThere:         true,
	},
	{
		OperatingSystem:           "darwin",
		Name:                      "macOS",
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

func Names() []string {
	var named []string
	for _, one := range Supported() {
		named = append(named, one.Name)
	}
	return named
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
