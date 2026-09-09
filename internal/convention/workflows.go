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

type outward struct {
	named   string
	targets []string
	tags    []string
}

var theRealCLI = outward{
	named:   "the real CLI",
	targets: []string{"fixtures", "live-smoke"},
	tags:    []string{"record", "livesmoke"},
}

var aForge = outward{
	named:   "a forge",
	targets: []string{"forge-smoke"},
	tags:    []string{"forgesmoke"},
}

var buildTagsNamed = regexp.MustCompile(`-tags[= ]([A-Za-z0-9_,]+)`)

var credentialsTheRealCLIReads = []string{"ANTHROPIC_API_KEY", "CLAUDE_CODE_OAUTH_TOKEN", "OPENAI_API_KEY", "CODEX_API_KEY", "OPENCODE_API_KEY"}

const TheEntryPoint = "verify"

var gatesTheEntryPointRuns = []string{"lint", "test", "corpus-outcomes", "fixtures-diff"}

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

func reachesOutward(content, named string, outside outward) []string {
	var reaching []string

	for _, target := range outside.targets {
		if invocationOf(target).MatchString(content) {
			reaching = append(reaching, fmt.Sprintf("%s runs make %s, which reaches %s", named, target, outside.named))
		}
	}
	for _, tagged := range buildTagsNamed.FindAllStringSubmatch(content, -1) {
		for tag := range strings.SplitSeq(tagged[1], ",") {
			if slices.Contains(outside.tags, tag) {
				reaching = append(reaching, fmt.Sprintf("%s compiles the %s tag, which is behind the paths that reach %s", named, tag, outside.named))
			}
		}
	}
	return reaching
}

type WorkflowReach struct {
	Reaching  []string
	Workflows []string
}

func WorkflowsReachingTheRealCLI(root string) (WorkflowReach, error) {
	return workflowsReaching(root, theRealCLI, credentialsTheRealCLIReads)
}

func WorkflowsReachingAForge(root string) (WorkflowReach, error) {
	return workflowsReaching(root, aForge, nil)
}

func workflowsReaching(root string, outside outward, credentials []string) (WorkflowReach, error) {
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

		reach.Reaching = append(reach.Reaching, reachesOutward(string(content), name, outside)...)
		for _, credential := range credentials {
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
