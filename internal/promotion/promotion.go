package promotion

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ktei/socrates/internal/coverage"
)

type Condition string

const (
	ScenariosDropped   Condition = "the promoted specification would stop holding scenarios it holds now"
	SubstitutesNothing Condition = "the requirement it modifies is in no promoted specification, so archiving would add it beside whatever is there"
	AddsWhatExists     Condition = "the promoted specification already holds a requirement under that heading"
	PromotedAhead      Condition = "the promoted specification already holds the new name of a scenario the delta renames, so there is nothing for the rename to be checked against"
)

type Drop struct {
	Change      string
	Spec        string
	Requirement string
	Condition   Condition
	Scenarios   []string
}

type Report struct {
	Changes     int
	Substituted int
	Added       int
	Drops       []Drop
}

const (
	specFile          = "spec.md"
	specsDirectory    = "specs"
	changesDirectory  = "changes"
	archiveDirectory  = "archive"
	openspecDirectory = "openspec"

	requirementHeading = "### Requirement:"
	scenarioHeading    = "#### Scenario:"
	sectionHeading     = "## "
	requirementsSuffix = " Requirements"

	added    = "ADDED"
	modified = "MODIFIED"
	renamed  = "RENAMED"

	renameFrom = "- FROM:"
	renameTo   = "- TO:"
)

type requirement struct {
	heading   string
	section   string
	scenarios []string
}

func Check(root string) (Report, error) {
	changes := filepath.Join(root, openspecDirectory, changesDirectory)
	_, err := os.Stat(changes)
	if errors.Is(err, fs.ErrNotExist) {
		return Report{}, nil
	}
	if err != nil {
		return Report{}, err
	}

	underWay, err := changesUnderWay(changes)
	if err != nil {
		return Report{}, err
	}

	report := Report{Changes: underWay}
	err = filepath.WalkDir(changes, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == archiveDirectory {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Base(path) != specFile || filepath.Base(filepath.Dir(filepath.Dir(path))) != specsDirectory {
			return nil
		}

		against, err := comparison(root, path)
		if err != nil {
			return err
		}
		report.Substituted += against.Substituted
		report.Added += against.Added
		report.Drops = append(report.Drops, against.Drops...)
		return nil
	})
	if err != nil {
		return Report{}, err
	}

	slices.SortStableFunc(report.Drops, func(one, other Drop) int {
		if one.Spec != other.Spec {
			return strings.Compare(one.Spec, other.Spec)
		}
		return strings.Compare(one.Requirement, other.Requirement)
	})
	return report, nil
}

func changesUnderWay(changes string) (int, error) {
	entries, err := os.ReadDir(changes)
	if err != nil {
		return 0, err
	}

	under := 0
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != archiveDirectory {
			under++
		}
	}
	return under, nil
}

func comparison(root, delta string) (Report, error) {
	capability := filepath.Base(filepath.Dir(delta))
	change := filepath.Base(filepath.Dir(filepath.Dir(filepath.Dir(delta))))

	proposed, renames, scenarioRenames, err := requirementsIn(delta)
	if err != nil {
		return Report{}, err
	}
	promoted, err := promotedRequirements(root, capability)
	if err != nil {
		return Report{}, err
	}

	relative, err := filepath.Rel(root, delta)
	if err != nil {
		relative = delta
	}

	ahead := renamesAlreadyPromoted(promoted, scenarioRenames)

	var report Report
	for _, one := range proposed {
		drop := Drop{Change: change, Spec: relative, Requirement: one.heading}

		switch one.section {
		case added:
			report.Added++
			if _, present := promoted[coverage.Normalise(one.heading)]; present {
				drop.Condition = AddsWhatExists
				report.Drops = append(report.Drops, drop)
			}
		case modified:
			report.Substituted++
			substituted := one.heading
			if from, wasRenamed := renames[coverage.Normalise(one.heading)]; wasRenamed {
				substituted = from
			}

			held, present := promoted[coverage.Normalise(substituted)]
			if !present {
				drop.Condition = SubstitutesNothing
				report.Drops = append(report.Drops, drop)
				continue
			}
			if carried := alsoCarried(ahead, one.scenarios); len(carried) > 0 {
				unreadable := drop
				unreadable.Condition = PromotedAhead
				unreadable.Scenarios = carried
				report.Drops = append(report.Drops, unreadable)
			}
			if missing := absent(held.scenarios, one.scenarios, scenarioRenames); len(missing) > 0 {
				dropped := drop
				dropped.Condition = ScenariosDropped
				dropped.Scenarios = missing
				report.Drops = append(report.Drops, dropped)
			}
		}
	}
	return report, nil
}

func promotedRequirements(root, capability string) (map[string]requirement, error) {
	path := filepath.Join(root, openspecDirectory, "specs", capability, specFile)
	held, _, _, err := requirementsIn(path)
	if err != nil {
		return nil, err
	}

	promoted := map[string]requirement{}
	for _, one := range held {
		promoted[coverage.Normalise(one.heading)] = one
	}
	return promoted, nil
}

func requirementsIn(path string) ([]requirement, map[string]string, map[string]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil, nil
		}
		return nil, nil, nil, fmt.Errorf(readingASpecification, path, err)
	}

	var held []requirement
	renames := map[string]string{}
	scenarioRenames := map[string]string{}
	section := ""
	from := ""
	fromScenario := ""

	for line := range strings.SplitSeq(string(content), "\n") {
		trimmed := strings.TrimSpace(line)

		if named, headed := strings.CutPrefix(trimmed, requirementHeading); headed {
			held = append(held, requirement{heading: strings.TrimSpace(named), section: section})
			continue
		}
		if named, headed := strings.CutPrefix(trimmed, scenarioHeading); headed {
			if len(held) > 0 {
				held[len(held)-1].scenarios = append(held[len(held)-1].scenarios, strings.TrimSpace(named))
			}
			continue
		}
		if stated, headed := strings.CutPrefix(trimmed, sectionHeading); headed {
			section = strings.ToUpper(strings.TrimSpace(strings.TrimSuffix(stated, requirementsSuffix)))
			continue
		}

		if section != renamed {
			continue
		}
		if named, marked := strings.CutPrefix(trimmed, renameFrom); marked {
			if scenario, isScenario := renamedScenario(named); isScenario {
				fromScenario = scenario
			} else {
				from = renamedRequirement(named)
			}
		}
		if named, marked := strings.CutPrefix(trimmed, renameTo); marked {
			if scenario, isScenario := renamedScenario(named); isScenario && fromScenario != "" {
				scenarioRenames[coverage.Normalise(fromScenario)] = scenario
				fromScenario = ""
			} else if !isScenario && from != "" {
				renames[coverage.Normalise(renamedRequirement(named))] = from
				from = ""
			}
		}
	}
	return held, renames, scenarioRenames, nil
}

func renamedScenario(line string) (string, bool) {
	quoted := strings.Trim(strings.TrimSpace(line), "`")
	named, headed := strings.CutPrefix(strings.TrimSpace(quoted), scenarioHeading)
	if !headed {
		return "", false
	}
	return strings.TrimSpace(named), true
}

func renamedRequirement(line string) string {
	quoted := strings.Trim(strings.TrimSpace(line), "`")
	named, headed := strings.CutPrefix(strings.TrimSpace(quoted), requirementHeading)
	if !headed {
		return strings.TrimSpace(quoted)
	}
	return strings.TrimSpace(named)
}

func absent(held, carried []string, renames map[string]string) []string {
	var missing []string
	for _, scenario := range held {
		wanted := scenario
		if renamed, was := renames[coverage.Normalise(scenario)]; was {
			wanted = renamed
		}
		if !slices.ContainsFunc(carried, func(one string) bool {
			return coverage.Normalise(one) == coverage.Normalise(wanted)
		}) {
			missing = append(missing, scenario)
		}
	}
	return missing
}

func renamesAlreadyPromoted(promoted map[string]requirement, renames map[string]string) []string {
	holds := func(name string) bool {
		for _, one := range promoted {
			if slices.ContainsFunc(one.scenarios, func(scenario string) bool {
				return coverage.Normalise(scenario) == coverage.Normalise(name)
			}) {
				return true
			}
		}
		return false
	}

	var ahead []string
	for from, to := range renames {
		if coverage.Normalise(from) == coverage.Normalise(to) {
			continue
		}
		if holds(to) && !holds(from) {
			ahead = append(ahead, to)
		}
	}
	slices.Sort(ahead)
	return ahead
}

func alsoCarried(ahead, carried []string) []string {
	var named []string
	for _, one := range ahead {
		if slices.ContainsFunc(carried, func(scenario string) bool {
			return coverage.Normalise(scenario) == coverage.Normalise(one)
		}) {
			named = append(named, one)
		}
	}
	return named
}
