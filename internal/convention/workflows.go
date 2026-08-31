package convention

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const WorkflowsDirectory = ".github/workflows"

var targetsReachingTheRealCLI = []string{"fixtures", "live-smoke"}

var tagsReachingTheRealCLI = []string{"record", "livesmoke"}

var buildTagsNamed = regexp.MustCompile(`-tags[= ]([A-Za-z0-9_,]+)`)

var credentialsTheRealCLIReads = []string{"ANTHROPIC_API_KEY", "CLAUDE_CODE_OAUTH_TOKEN"}

const TheEntryPoint = "verify"

var gatesTheEntryPointRuns = []string{"lint", "test", "scenario-coverage", "scenario-drop", "corpus-outcomes", "fixtures-diff"}

type WorkflowGates struct {
	Invoking []string
	Listing  []string
}

func WorkflowsRunningTheGates(root string) (WorkflowGates, error) {
	var gates WorkflowGates

	paths, err := workflowPaths(filepath.Join(root, WorkflowsDirectory))
	if err != nil {
		return gates, err
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return gates, err
		}

		name := filepath.Base(path)
		if invocationOf(TheEntryPoint).Match(content) {
			gates.Invoking = append(gates.Invoking, name)
		}
		for _, gate := range gatesTheEntryPointRuns {
			if invocationOf(gate).Match(content) {
				gates.Listing = append(gates.Listing, fmt.Sprintf("%s runs make %s of its own", name, gate))
			}
		}
	}

	slices.Sort(gates.Invoking)
	slices.Sort(gates.Listing)
	return gates, nil
}

type WorkflowReach struct {
	Reaching  []string
	Workflows []string
}

func WorkflowsReachingTheRealCLI(root string) (WorkflowReach, error) {
	var reach WorkflowReach

	paths, err := workflowPaths(filepath.Join(root, WorkflowsDirectory))
	if err != nil {
		return reach, err
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return reach, err
		}

		name := filepath.Base(path)
		reach.Workflows = append(reach.Workflows, name)

		for _, target := range targetsReachingTheRealCLI {
			if invocationOf(target).Match(content) {
				reach.Reaching = append(reach.Reaching, fmt.Sprintf("%s runs make %s, which reaches the real CLI", name, target))
			}
		}
		for _, tagged := range buildTagsNamed.FindAllStringSubmatch(string(content), -1) {
			for _, tag := range strings.Split(tagged[1], ",") {
				if slices.Contains(tagsReachingTheRealCLI, tag) {
					reach.Reaching = append(reach.Reaching, fmt.Sprintf("%s compiles the %s tag, which is behind the paths that reach the real CLI", name, tag))
				}
			}
		}
		for _, credential := range credentialsTheRealCLIReads {
			if regexp.MustCompile(regexp.QuoteMeta(credential)).Match(content) {
				reach.Reaching = append(reach.Reaching, fmt.Sprintf("%s names %s, which is a credential the real CLI reads", name, credential))
			}
		}
	}
	return reach, nil
}

var runnerNamed = regexp.MustCompile(`(?m)^\s*runs-on:\s*(\S+)\s*$`)

var derived = regexp.MustCompile(`\$\{\{`)

type NamedRunner struct {
	Workflow string
	Runner   string
}

func RunnersNamedInWorkflows(root string) ([]NamedRunner, error) {
	var named []NamedRunner

	paths, err := workflowPaths(filepath.Join(root, WorkflowsDirectory))
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		for _, running := range runnerNamed.FindAllStringSubmatch(string(content), -1) {
			if derived.MatchString(running[1]) {
				continue
			}
			named = append(named, NamedRunner{Workflow: filepath.Base(path), Runner: strings.Trim(running[1], `"'`)})
		}
	}

	slices.SortFunc(named, func(one, other NamedRunner) int {
		if one.Workflow != other.Workflow {
			return strings.Compare(one.Workflow, other.Workflow)
		}
		return strings.Compare(one.Runner, other.Runner)
	})
	return named, nil
}

func workflowPaths(directory string) ([]string, error) {
	var paths []string
	for _, extension := range []string{"*.yml", "*.yaml"} {
		matched, err := filepath.Glob(filepath.Join(directory, extension))
		if err != nil {
			return nil, err
		}
		paths = append(paths, matched...)
	}
	slices.Sort(paths)
	return paths, nil
}

func invocationOf(target string) *regexp.Regexp {
	return regexp.MustCompile(`\bmake\s+` + regexp.QuoteMeta(target) + `(?:[^\w-]|$)`)
}
