package convention

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"
)

const (
	ReadOnlyToken         = "contents: read"
	CancelSuperseded      = "cancel-in-progress: true"
	DependabotFile        = ".github/dependabot.yml"
	TheVulnerabilityCheck = "govulncheck"
	TheCodeScan           = "github/codeql-action/analyze"
)

var ecosystemsADependabotWatches = []string{"gomod", "npm", "github-actions"}

var actionNamed = regexp.MustCompile(`(?m)^\s*(?:-\s+)?uses:\s*(\S+)`)

var pinnedToACommit = regexp.MustCompile(`^[0-9a-f]{40}$`)

type NamedAction struct {
	Workflow string
	Uses     string
	Named    string
}

func ActionsNamedByATag(root string) ([]NamedAction, error) {
	var tagged []NamedAction

	paths, err := workflowPaths(filepath.Join(root, WorkflowsDirectory))
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		for _, used := range actionNamed.FindAllStringSubmatch(string(content), -1) {
			action := strings.Trim(used[1], `"'`)
			if strings.HasPrefix(action, "./") {
				continue
			}
			named, at, versioned := strings.Cut(action, "@")
			if !versioned {
				tagged = append(tagged, NamedAction{Workflow: filepath.Base(path), Uses: named})
				continue
			}
			if !pinnedToACommit.MatchString(at) {
				tagged = append(tagged, NamedAction{Workflow: filepath.Base(path), Uses: named, Named: at})
			}
		}
	}

	slices.SortFunc(tagged, func(one, other NamedAction) int {
		if one.Workflow != other.Workflow {
			return strings.Compare(one.Workflow, other.Workflow)
		}
		return strings.Compare(one.Uses, other.Uses)
	})
	return tagged, nil
}

type WorkflowClaim struct {
	Workflow string
	Stated   string
}

func TokensWorkflowsAskFor(root string) ([]WorkflowClaim, error) {
	return claimed(root, "permissions")
}

func WhatWorkflowsSupersede(root string) ([]WorkflowClaim, error) {
	return claimed(root, "concurrency")
}

func claimed(root, key string) ([]WorkflowClaim, error) {
	var claims []WorkflowClaim

	paths, err := workflowPaths(filepath.Join(root, WorkflowsDirectory))
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		claims = append(claims, WorkflowClaim{Workflow: filepath.Base(path), Stated: topLevelBlock(string(content), key)})
	}
	return claims, nil
}

func topLevelBlock(document, key string) string {
	var held []string
	inside := false

	for _, line := range strings.Split(document, "\n") {
		if strings.HasPrefix(line, key+":") {
			inside = true
			held = append(held, strings.TrimSpace(strings.TrimPrefix(line, key+":")))
			continue
		}
		if !inside {
			continue
		}
		if line == "" || strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			held = append(held, strings.TrimSpace(line))
			continue
		}
		break
	}

	if !inside {
		return ""
	}
	return strings.TrimSpace(strings.Join(held, "\n"))
}

const TheBrowserSeam = "AGORA_BROWSER"

var browserVersionNamed = regexp.MustCompile(`chrome-version:\s*'?([0-9][0-9.]*)'?`)

var browserResolvedFromTheRunner = regexp.MustCompile(`command -v\s+(?:google-chrome|chromium)`)

type BrowserInAWorkflow struct {
	Workflow      string
	Version       string
	FromTheRunner bool
	NamesTheSeam  bool
}

func BrowsersWorkflowsDrive(root string) ([]BrowserInAWorkflow, error) {
	var driving []BrowserInAWorkflow

	paths, err := workflowPaths(filepath.Join(root, WorkflowsDirectory))
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if !invocationOf(TheEntryPoint).Match(content) {
			continue
		}

		found := BrowserInAWorkflow{
			Workflow:      filepath.Base(path),
			FromTheRunner: browserResolvedFromTheRunner.Match(content),
			NamesTheSeam:  strings.Contains(string(content), TheBrowserSeam),
		}
		if named := browserVersionNamed.FindSubmatch(content); named != nil {
			found.Version = string(named[1])
		}
		driving = append(driving, found)
	}
	return driving, nil
}

func WhenWorkflowsRun(root string) ([]WorkflowClaim, error) {
	return claimed(root, "on")
}

var triggersNobodyAsksFor = []string{"push", "pull_request", "schedule"}

const OnRequest = "workflow_dispatch"

type WorkflowStart struct {
	Workflow  string
	Unasked   []string
	OnRequest bool
}

func WhatStartsAWorkflow(root string, named string) (WorkflowStart, error) {
	start := WorkflowStart{Workflow: named}

	held, err := os.ReadFile(filepath.Join(root, WorkflowsDirectory, named))
	if err != nil {
		return start, err
	}

	stated := topLevelBlock(string(held), "on")
	start.OnRequest = strings.Contains(stated, OnRequest)
	for _, trigger := range triggersNobodyAsksFor {
		if regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(trigger) + `:?$`).MatchString(stated) {
			start.Unasked = append(start.Unasked, trigger)
		}
	}
	return start, nil
}

func WorkflowsChecking(root, checked string) ([]string, error) {
	var checking []string

	paths, err := workflowPaths(filepath.Join(root, WorkflowsDirectory))
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if strings.Contains(string(content), checked) {
			checking = append(checking, filepath.Base(path))
		}
	}
	return checking, nil
}

func WorkflowsReadingTheStatement(root, statement string) (WorkflowGates, error) {
	var reading WorkflowGates

	paths, err := workflowPaths(filepath.Join(root, WorkflowsDirectory))
	if err != nil {
		return reading, err
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return reading, err
		}
		if !invocationOf(TheEntryPoint).Match(content) {
			continue
		}

		name := filepath.Base(path)
		if strings.Contains(string(content), statement) {
			reading.Invoking = append(reading.Invoking, name)
			continue
		}
		reading.Listing = append(reading.Listing, fmt.Sprintf("%s runs the suite without reading %s", name, statement))
	}
	return reading, nil
}

type dependabotConfiguration struct {
	Updates []struct {
		Ecosystem string `yaml:"package-ecosystem"`
		Directory string `yaml:"directory"`
	} `yaml:"updates"`
}

func EcosystemsUnwatched(root string) ([]string, error) {
	held, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(DependabotFile)))
	if err != nil {
		return nil, err
	}

	var configured dependabotConfiguration
	if err := yaml.Unmarshal(held, &configured); err != nil {
		return nil, fmt.Errorf("reading %s: %w", DependabotFile, err)
	}

	watched := map[string]string{}
	for _, update := range configured.Updates {
		watched[update.Ecosystem] = update.Directory
	}

	var unwatched []string
	for _, ecosystem := range ecosystemsADependabotWatches {
		directory, proposing := watched[ecosystem]
		if !proposing {
			unwatched = append(unwatched, fmt.Sprintf("%s proposes no update for %s", DependabotFile, ecosystem))
			continue
		}
		entry, err := os.Stat(filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(directory, "/"))))
		if err != nil || !entry.IsDir() {
			unwatched = append(unwatched, fmt.Sprintf("%s proposes an update for %s in %s, and there is no such directory", DependabotFile, ecosystem, directory))
		}
	}
	return unwatched, nil
}
