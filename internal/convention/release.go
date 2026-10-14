package convention

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"
)

const ReleaseConfigurationFile = ".goreleaser.yaml"

type ReleaseBuild struct {
	ID               string
	Binary           string
	Compiles         string
	OperatingSystems []string
	Stamps           []string
}

type ReleaseArchive struct {
	ID     string
	Builds []string
	Files  []string
}

type releaseConfiguration struct {
	Builds []struct {
		ID     string   `yaml:"id"`
		Binary string   `yaml:"binary"`
		Main   string   `yaml:"main"`
		Goos   []string `yaml:"goos"`
		Ldflag []string `yaml:"ldflags"`
	} `yaml:"builds"`
	Archives []struct {
		ID    string   `yaml:"id"`
		IDs   []string `yaml:"ids"`
		Files []struct {
			Source string `yaml:"src"`
		} `yaml:"files"`
	} `yaml:"archives"`
}

func AReleaseBuilds(root string) ([]ReleaseBuild, error) {
	configured, err := releaseAsConfigured(root)
	if err != nil {
		return nil, err
	}

	var building []ReleaseBuild
	for _, one := range configured.Builds {
		building = append(building, ReleaseBuild{
			ID:               one.ID,
			Binary:           one.Binary,
			Compiles:         one.Main,
			OperatingSystems: slices.Clone(one.Goos),
			Stamps:           slices.Clone(one.Ldflag),
		})
	}
	return building, nil
}

func AReleaseArchives(root string) ([]ReleaseArchive, error) {
	configured, err := releaseAsConfigured(root)
	if err != nil {
		return nil, err
	}

	var archiving []ReleaseArchive
	for _, one := range configured.Archives {
		carried := ReleaseArchive{ID: one.ID, Builds: slices.Clone(one.IDs)}
		for _, file := range one.Files {
			carried.Files = append(carried.Files, file.Source)
		}
		archiving = append(archiving, carried)
	}
	return archiving, nil
}

func OperatingSystemsAReleaseBuilds(root string) ([]string, error) {
	building, err := AReleaseBuilds(root)
	if err != nil {
		return nil, err
	}

	var named []string
	for _, one := range building {
		for _, operating := range one.OperatingSystems {
			if !slices.Contains(named, operating) {
				named = append(named, operating)
			}
		}
	}
	slices.Sort(named)
	return named, nil
}

const ReleaseWorkflow = "release.yml"

type ReleaseGate struct {
	Verifying string
	Publishes []string
	Ungated   []string
}

type workflowJobs struct {
	Jobs map[string]struct {
		Needs yaml.Node `yaml:"needs"`
		Uses  string    `yaml:"uses"`
		Steps []struct {
			Uses string         `yaml:"uses"`
			With map[string]any `yaml:"with"`
		} `yaml:"steps"`
	} `yaml:"jobs"`
}

func releaseWorkflow(root string) (workflowJobs, error) {
	var document workflowJobs

	held, err := os.ReadFile(filepath.Join(root, WorkflowsDirectory, ReleaseWorkflow))
	if err != nil {
		return document, err
	}
	if err := yaml.Unmarshal(held, &document); err != nil {
		return document, fmt.Errorf("reading %s: %w", ReleaseWorkflow, err)
	}
	return document, nil
}

func HowTheReleaseIsGated(root string) (ReleaseGate, error) {
	var gate ReleaseGate

	document, err := releaseWorkflow(root)
	if err != nil {
		return gate, err
	}

	running, err := WorkflowsRunningTheGates(root)
	if err != nil {
		return gate, err
	}

	needs := map[string][]string{}
	for name, job := range document.Jobs {
		needs[name] = named(job.Needs)
		if job.Uses != "" && slices.Contains(running.Invoking, filepath.Base(job.Uses)) {
			gate.Verifying = name
		}
	}

	for name := range document.Jobs {
		if name == gate.Verifying {
			continue
		}
		gate.Publishes = append(gate.Publishes, name)
		if !reaches(needs, name, gate.Verifying) {
			gate.Ungated = append(gate.Ungated, name)
		}
	}

	slices.Sort(gate.Publishes)
	slices.Sort(gate.Ungated)
	return gate, nil
}

type ReleaseScan struct {
	Scanning  string
	Unscanned []string
}

func HowTheReleaseIsScanned(root string) (ReleaseScan, error) {
	var scan ReleaseScan

	document, err := releaseWorkflow(root)
	if err != nil {
		return scan, err
	}

	scanning, err := WorkflowsChecking(root, TheCodeScan)
	if err != nil {
		return scan, err
	}

	needs := map[string][]string{}
	for name, job := range document.Jobs {
		needs[name] = named(job.Needs)
		if job.Uses != "" && slices.Contains(scanning, filepath.Base(job.Uses)) {
			scan.Scanning = name
		}
	}

	for name, job := range document.Jobs {
		if job.Uses != "" {
			continue
		}
		if !reaches(needs, name, scan.Scanning) {
			scan.Unscanned = append(scan.Unscanned, name)
		}
	}

	slices.Sort(scan.Unscanned)
	return scan, nil
}

const ProvenanceAction = "actions/attest-build-provenance"

type ReleaseSignature struct {
	Over    string
	Command string
	Keyless bool
}

type ReleaseVerification struct {
	Checksums  string
	BillsOver  []string
	Signatures []ReleaseSignature
	Provenance []string
	Archives   []string
}

type verificationConfiguration struct {
	Checksum struct {
		NameTemplate string `yaml:"name_template"`
	} `yaml:"checksum"`
	SBOMs []struct {
		Artifacts string `yaml:"artifacts"`
	} `yaml:"sboms"`
	Signs []struct {
		Cmd       string   `yaml:"cmd"`
		Args      []string `yaml:"args"`
		Artifacts string   `yaml:"artifacts"`
	} `yaml:"signs"`
	Archives []struct {
		NameTemplate string `yaml:"name_template"`
	} `yaml:"archives"`
}

func WhatAReleaseCarriesForVerification(root string) (ReleaseVerification, error) {
	var carried ReleaseVerification

	held, err := os.ReadFile(filepath.Join(root, ReleaseConfigurationFile))
	if err != nil {
		return carried, err
	}

	var configured verificationConfiguration
	if err := yaml.Unmarshal(held, &configured); err != nil {
		return carried, fmt.Errorf("reading %s: %w", ReleaseConfigurationFile, err)
	}

	carried.Checksums = configured.Checksum.NameTemplate
	for _, bill := range configured.SBOMs {
		carried.BillsOver = append(carried.BillsOver, bill.Artifacts)
	}
	for _, signature := range configured.Signs {
		carried.Signatures = append(carried.Signatures, ReleaseSignature{
			Over:    signature.Artifacts,
			Command: signature.Cmd,
			Keyless: !slices.ContainsFunc(signature.Args, func(argument string) bool { return strings.Contains(argument, "--key") }),
		})
	}
	for _, archive := range configured.Archives {
		carried.Archives = append(carried.Archives, archive.NameTemplate)
	}

	document, err := releaseWorkflow(root)
	if err != nil {
		return carried, err
	}
	for _, job := range document.Jobs {
		for _, step := range job.Steps {
			if action, _, _ := strings.Cut(step.Uses, "@"); action != ProvenanceAction {
				continue
			}
			if subject, stated := step.With["subject-path"]; stated {
				carried.Provenance = append(carried.Provenance, fmt.Sprint(subject))
			}
		}
	}

	slices.Sort(carried.Provenance)
	return carried, nil
}

func WhatTheReleaseReaches(root string) ([]string, error) {
	held, err := os.ReadFile(filepath.Join(root, ReleaseConfigurationFile))
	if err != nil {
		return nil, err
	}

	var reaching []string
	for _, outside := range []outward{theRealCLI, aForge} {
		reaching = append(reaching, reachesOutward(string(held), ReleaseConfigurationFile, outside)...)
	}
	for _, credential := range credentialsTheRealCLIReads {
		if strings.Contains(string(held), credential) {
			reaching = append(reaching, fmt.Sprintf("%s names %s, which is a credential the real CLI reads", ReleaseConfigurationFile, credential))
		}
	}

	slices.Sort(reaching)
	return reaching, nil
}

const (
	ChangelogFile = "CHANGELOG.md"
	Unreleased    = "Unreleased"

	entryHeading = "## ["
)

type RecordedRelease struct {
	Version string
	Dated   string
	Says    string
}

func ReleasesRecorded(root string) ([]RecordedRelease, error) {
	held, err := os.ReadFile(filepath.Join(root, ChangelogFile))
	if err != nil {
		return nil, err
	}

	var recorded []RecordedRelease
	for _, section := range strings.Split(string(held), entryHeading)[1:] {
		heading, body, _ := strings.Cut(section, "\n")
		version, dated, _ := strings.Cut(heading, "]")

		if version == Unreleased {
			continue
		}
		if ended := strings.Index(body, "\n## "); ended >= 0 {
			body = body[:ended]
		}
		recorded = append(recorded, RecordedRelease{
			Version: version,
			Dated:   strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(dated), "-")),
			Says:    strings.TrimSpace(body),
		})
	}
	return recorded, nil
}

func named(node yaml.Node) []string {
	var one string
	if err := node.Decode(&one); err == nil {
		return []string{one}
	}

	var many []string
	if err := node.Decode(&many); err == nil {
		return many
	}
	return nil
}

func reaches(needs map[string][]string, from, wanted string) bool {
	if wanted == "" {
		return false
	}

	seen := map[string]bool{}
	pending := []string{from}
	for len(pending) > 0 {
		name := pending[0]
		pending = pending[1:]
		if seen[name] {
			continue
		}
		seen[name] = true

		for _, needed := range needs[name] {
			if needed == wanted {
				return true
			}
			pending = append(pending, needed)
		}
	}
	return false
}

func releaseAsConfigured(root string) (releaseConfiguration, error) {
	var configured releaseConfiguration

	held, err := os.ReadFile(filepath.Join(root, ReleaseConfigurationFile))
	if err != nil {
		return configured, err
	}
	if err := yaml.Unmarshal(held, &configured); err != nil {
		return configured, fmt.Errorf("reading %s: %w", ReleaseConfigurationFile, err)
	}
	return configured, nil
}

func BuildsProducing(building []ReleaseBuild, binary string) []ReleaseBuild {
	var producing []ReleaseBuild
	for _, one := range building {
		if one.Binary == binary {
			producing = append(producing, one)
		}
	}
	return producing
}

var stampedSymbol = regexp.MustCompile(`-X\s+(\S+?)=`)

func SymbolsAReleaseStamps(stamps []string) []string {
	var named []string
	for _, stamp := range stamps {
		for _, written := range stampedSymbol.FindAllStringSubmatch(stamp, -1) {
			qualified := written[1]
			position := strings.LastIndex(qualified, ".")
			if position < 0 {
				continue
			}
			named = append(named, qualified[:position])
		}
	}
	return named
}

const InstallerFile = "install.sh"

type ReleaseAttachments struct {
	Attached    []string
	Checksummed []string
}

type attachmentConfiguration struct {
	Release struct {
		ExtraFiles []struct {
			Glob string `yaml:"glob"`
		} `yaml:"extra_files"`
	} `yaml:"release"`
	Checksum struct {
		ExtraFiles []struct {
			Glob string `yaml:"glob"`
		} `yaml:"extra_files"`
	} `yaml:"checksum"`
}

func WhatAReleaseAttaches(root string) (ReleaseAttachments, error) {
	var attaching ReleaseAttachments

	held, err := os.ReadFile(filepath.Join(root, ReleaseConfigurationFile))
	if err != nil {
		return attaching, err
	}

	var configured attachmentConfiguration
	if err := yaml.Unmarshal(held, &configured); err != nil {
		return attaching, fmt.Errorf("reading %s: %w", ReleaseConfigurationFile, err)
	}

	for _, extra := range configured.Release.ExtraFiles {
		attaching.Attached = append(attaching.Attached, path.Clean(extra.Glob))
	}
	for _, extra := range configured.Checksum.ExtraFiles {
		attaching.Checksummed = append(attaching.Checksummed, path.Clean(extra.Glob))
	}
	return attaching, nil
}
