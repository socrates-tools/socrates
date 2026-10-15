package skill

import (
	"fmt"
	"strings"

	"github.com/socrates-tools/socrates/internal/refusal"
)

const (
	ToFlag       = "to"
	AgentsToFlag = "agents-to"
)

func HostUsage() string {
	var flags []string
	for _, host := range Hosts() {
		flags = append(flags, host.Flag)
	}
	return fmt.Sprintf(hostUsage, strings.Join(flags, hostFlagsJoin))
}

type Directory struct {
	Flag  string
	Value string
}

func Selected(named string, directories ...Directory) ([]Host, error) {
	if named == "" {
		for _, one := range directories {
			if one.Value != "" {
				return nil, refusal.New(aDirectoryNeedsAHost, one.Flag, HostFlag, hostFlagsStated())
			}
		}
		return Hosts(), nil
	}

	host, found := HostFlagged(named)
	if !found {
		return nil, refusal.New(noSuchHost, HostFlag, named, hostFlagsStated())
	}
	return []Host{host}, nil
}

func hostFlagsStated() string {
	var stated []string
	for _, host := range Hosts() {
		stated = append(stated, fmt.Sprintf(hostFlagStated, HostFlag, host.Flag))
	}
	return strings.Join(stated, hostFlagsJoin)
}

func InstallFor(hosts []Host, target, invocation string, skillsFor func(Host) []Skill) ([]Installed, error) {
	var written []Installed

	for _, host := range hosts {
		directory := target
		if directory == "" {
			chosen, err := UserDirectory(host, invocation)
			if err != nil {
				return nil, err
			}
			directory = chosen
		}

		installed, err := Install(directory, skillsFor(host), invocation)
		if err != nil {
			return nil, err
		}
		written = append(written, installed...)
	}
	return written, nil
}
