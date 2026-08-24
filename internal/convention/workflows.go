package convention

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
)

const WorkflowsDirectory = ".github/workflows"

var targetsReachingTheRealCLI = []string{"fixtures", "live-smoke"}

var credentialsTheRealCLIReads = []string{"ANTHROPIC_API_KEY", "CLAUDE_CODE_OAUTH_TOKEN"}

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
		for _, credential := range credentialsTheRealCLIReads {
			if regexp.MustCompile(regexp.QuoteMeta(credential)).Match(content) {
				reach.Reaching = append(reach.Reaching, fmt.Sprintf("%s names %s, which is a credential the real CLI reads", name, credential))
			}
		}
	}
	return reach, nil
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
