package convention

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const (
	FrontDocumentFile      = "README.md"
	SecurityPolicyFile     = "SECURITY.md"
	ContributingFile       = "CONTRIBUTING.md"
	CodeOfConductFile      = "CODE_OF_CONDUCT.md"
	SupportFile            = "SUPPORT.md"
	GotchasFile            = "docs/gotchas.md"
	OwnershipFile          = ".github/CODEOWNERS"
	PullRequestTemplate    = ".github/PULL_REQUEST_TEMPLATE.md"
	IssueTemplateDirectory = ".github/ISSUE_TEMPLATE"
	EditorSettingsFile     = ".editorconfig"
	AttributeSettingsFile  = ".gitattributes"

	RequirementsHeading = "## What it needs"
	ToolchainHeading    = "## The toolchain"
	OutOfScopeHeading   = "## What socrates-review does not defend against"

	ReviewerCLIFloor   = "2.1.169"
	ProposalHome       = "GitHub issue"
	PrivateChannelPath = "security/advisories/new"

	SendsNothing     = "socrates-review sends nothing"
	SendsToTheVendor = "sends the change to"
	NotAffiliated    = "not affiliated"

	namedToBeExplained = 3
)

func Carries(root, named string) bool {
	_, err := os.Stat(filepath.Join(root, filepath.FromSlash(named)))
	return err == nil
}

func Document(root, named string) (string, error) {
	held, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(named)))
	if err != nil {
		return "", err
	}
	return string(held), nil
}

type FrontMatter struct {
	Sections []string
	Missing  []string
}

func RequirementsStatedFirst(root string, facts []string) (FrontMatter, error) {
	var front FrontMatter

	held, err := Document(root, FrontDocumentFile)
	if err != nil {
		return front, err
	}

	for line := range strings.SplitSeq(held, "\n") {
		if strings.HasPrefix(line, "## ") {
			front.Sections = append(front.Sections, strings.TrimSpace(line))
		}
	}

	section, found := sectionUnder(held, RequirementsHeading)
	if !found {
		front.Missing = append(front.Missing, fmt.Sprintf("%s has no %q section at all", FrontDocumentFile, RequirementsHeading))
		return front, nil
	}
	for _, fact := range facts {
		if !strings.Contains(section, fact) {
			front.Missing = append(front.Missing, fmt.Sprintf("%s states what socrates-review needs without naming %s", FrontDocumentFile, fact))
		}
	}
	return front, nil
}

type StatedFloor struct {
	Document string
	Floor    string
}

var moduleFloor = regexp.MustCompile(`(?m)^go\s+([0-9]+\.[0-9]+)`)

var documentedFloor = regexp.MustCompile(`Go ([0-9]+\.[0-9]+)`)

var documentsStatingTheFloor = []string{FrontDocumentFile, ContributingFile}

func ToolchainFloorRequired(root string) (string, error) {
	held, err := Document(root, "go.mod")
	if err != nil {
		return "", err
	}

	required := moduleFloor.FindStringSubmatch(held)
	if required == nil {
		return "", fmt.Errorf("go.mod requires no toolchain version")
	}
	return required[1], nil
}

func ToolchainFloorsStated(root string) ([]StatedFloor, error) {
	var stated []StatedFloor

	for _, named := range documentsStatingTheFloor {
		held, err := Document(root, named)
		if err != nil {
			return nil, err
		}
		for _, floor := range documentedFloor.FindAllStringSubmatch(held, -1) {
			stated = append(stated, StatedFloor{Document: named, Floor: floor[1]})
		}
	}
	return stated, nil
}

func FloorDisagreements(required string, stated []StatedFloor) []string {
	var disagreeing []string

	for _, one := range stated {
		if one.Floor != required {
			disagreeing = append(disagreeing, fmt.Sprintf("%s states Go %s and the module requires Go %s", one.Document, one.Floor, required))
		}
	}
	slices.Sort(disagreeing)
	return slices.Compact(disagreeing)
}

var installedBinary = regexp.MustCompile(`go install[^\n]*\./(?:[A-Za-z0-9_.-]+/)*cmd/(\S+)`)

func BinariesInstalled(root string) ([]string, error) {
	held, err := Document(root, "Makefile")
	if err != nil {
		return nil, err
	}

	var installed []string
	for _, named := range installedBinary.FindAllStringSubmatch(held, -1) {
		if !slices.Contains(installed, named[1]) {
			installed = append(installed, named[1])
		}
	}
	slices.Sort(installed)
	return installed, nil
}

func BinariesUnexplained(root string) ([]string, error) {
	installed, err := BinariesInstalled(root)
	if err != nil {
		return nil, err
	}

	held, err := Document(root, FrontDocumentFile)
	if err != nil {
		return nil, err
	}

	var unexplained []string
	for _, binary := range installed {
		explained := false
		for _, section := range strings.Split(held, "\n## ")[1:] {
			if strings.Count(section, binary) >= namedToBeExplained {
				explained = true
				break
			}
		}
		if !explained {
			unexplained = append(unexplained, fmt.Sprintf("%s puts %s on a path and %s explains it nowhere", "make install", binary, FrontDocumentFile))
		}
	}
	return unexplained, nil
}

var targetsNamed = []*regexp.Regexp{
	regexp.MustCompile("`make ([a-z][a-z-]*)`"),
	regexp.MustCompile(`(?m)^make ([a-z][a-z-]*)\b`),
}

var documentsTellingAContributorWhatToRun = []string{FrontDocumentFile, ContributingFile}

func TargetsDocumentsTell(root string) ([]string, error) {
	var told []string

	for _, named := range documentsTellingAContributorWhatToRun {
		held, err := Document(root, named)
		if err != nil {
			return nil, err
		}
		for _, naming := range targetsNamed {
			for _, target := range naming.FindAllStringSubmatch(held, -1) {
				if !slices.Contains(told, target[1]) {
					told = append(told, target[1])
				}
			}
		}
	}

	slices.Sort(told)
	return told, nil
}

var documentsStatingTheReach = []string{FrontDocumentFile, SecurityPolicyFile}

func ReachUnstated(root string) ([]string, error) {
	var unstated []string

	for _, named := range documentsStatingTheReach {
		held, err := Document(root, named)
		if err != nil {
			return nil, err
		}

		itself := strings.Contains(held, SendsNothing)
		vendor := strings.Contains(held, SendsToTheVendor)

		switch {
		case itself && !vendor:
			unstated = append(unstated, fmt.Sprintf("%s says %q and never says the reviewer CLI it drives %s its vendor, which reads as complete and is not", named, SendsNothing, SendsToTheVendor))
		case vendor && !itself:
			unstated = append(unstated, fmt.Sprintf("%s says the reviewer CLI %s its vendor and never says %q", named, SendsToTheVendor, SendsNothing))
		case !itself && !vendor:
			unstated = append(unstated, fmt.Sprintf("%s states neither half of where what socrates-review reads goes, and silence there is read as unknown", named))
		}
	}
	return unstated, nil
}

func NonAffiliationStated(root string) (bool, error) {
	held, err := Document(root, FrontDocumentFile)
	if err != nil {
		return false, err
	}
	return strings.Contains(held, NotAffiliated), nil
}

var respondsWithin = regexp.MustCompile(`\b[0-9]+ days\b`)

func SecurityPolicyGaps(root string) ([]string, error) {
	held, err := Document(root, SecurityPolicyFile)
	if err != nil {
		return nil, err
	}

	var gaps []string
	if !strings.Contains(held, PrivateChannelPath) {
		gaps = append(gaps, fmt.Sprintf("%s names no private channel, so a reporter's only way in is a public issue", SecurityPolicyFile))
	}
	if !respondsWithin.MatchString(held) {
		gaps = append(gaps, fmt.Sprintf("%s states no window in which a report is answered", SecurityPolicyFile))
	}

	scope, found := sectionUnder(held, OutOfScopeHeading)
	if !found {
		gaps = append(gaps, fmt.Sprintf("%s has no %q section, so a reporter cannot tell whether what they found is worth sending", SecurityPolicyFile, OutOfScopeHeading))
		return gaps, nil
	}
	if !strings.Contains(scope, "hook") {
		gaps = append(gaps, fmt.Sprintf("%s says what it does not defend against without naming the hook an administrator's policy configures, which is the residual %s already states", SecurityPolicyFile, FrontDocumentFile))
	}
	return gaps, nil
}

type GuideReading struct {
	ProposalAt  int
	ToolchainAt int
	NamesEntry  bool
	Listing     []string
}

func HowTheGuideReads(root string) (GuideReading, error) {
	var reading GuideReading

	held, err := Document(root, ContributingFile)
	if err != nil {
		return reading, err
	}

	reading.ProposalAt = strings.Index(held, ProposalHome)
	reading.ToolchainAt = strings.Index(held, ToolchainHeading)
	reading.NamesEntry = invocationOf(TheEntryPoint).MatchString(held)

	for _, gate := range gatesTheEntryPointRuns {
		if invocationOf(gate).MatchString(held) {
			reading.Listing = append(reading.Listing, fmt.Sprintf("%s tells a contributor to run make %s of its own", ContributingFile, gate))
		}
	}
	return reading, nil
}

var documentsAContributorNeeds = []string{
	CodeOfConductFile,
	ContributingFile,
	SecurityPolicyFile,
	SupportFile,
	OwnershipFile,
	PullRequestTemplate,
}

func DocumentsMissing(root string) ([]string, error) {
	var missing []string

	for _, named := range documentsAContributorNeeds {
		if !Carries(root, named) {
			missing = append(missing, fmt.Sprintf("the repository carries no %s", named))
			continue
		}
		held, err := Document(root, named)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(held) == "" {
			missing = append(missing, fmt.Sprintf("%s is empty", named))
		}
	}

	if Carries(root, OwnershipFile) {
		owners, err := Document(root, OwnershipFile)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(owners, "@") {
			missing = append(missing, fmt.Sprintf("%s names nobody, so ownership is stated by an empty file", OwnershipFile))
		}
	}

	if Carries(root, CodeOfConductFile) {
		conduct, err := Document(root, CodeOfConductFile)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(conduct, PrivateChannelPath) && !strings.Contains(conduct, "@") {
			missing = append(missing, fmt.Sprintf("%s names no contact, so a report has nowhere to go", CodeOfConductFile))
		}
	}

	if !Carries(root, IssueTemplateDirectory) {
		return append(missing, fmt.Sprintf("the repository carries no %s", IssueTemplateDirectory)), nil
	}

	templates, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(IssueTemplateDirectory)))
	if err != nil {
		return nil, err
	}

	forms := 0
	for _, entry := range templates {
		if !entry.IsDir() && entry.Name() != "config.yml" {
			forms++
		}
	}
	if forms == 0 {
		missing = append(missing, fmt.Sprintf("%s holds no template", IssueTemplateDirectory))
	}
	return missing, nil
}

var editorSettingsTheConventionsNeed = []string{
	"root = true",
	"charset = utf-8",
	"end_of_line = lf",
	"insert_final_newline = true",
	"indent_style = tab",
}

var attributeSettingsTheConventionsNeed = []string{"text=auto", "eol=lf"}

var textInTheBundle = []string{".css", ".js", ".json", ".md", ".gohtml"}

func EditorSettingsDisagreeing(root string) ([]string, error) {
	var disagreeing []string

	if !Carries(root, EditorSettingsFile) {
		return []string{fmt.Sprintf("the repository carries no %s, so a contributor's editor decides the indentation and the line endings", EditorSettingsFile)}, nil
	}

	settings, err := Document(root, EditorSettingsFile)
	if err != nil {
		return nil, err
	}
	for _, stated := range editorSettingsTheConventionsNeed {
		if !strings.Contains(settings, stated) {
			disagreeing = append(disagreeing, fmt.Sprintf("%s states no %q", EditorSettingsFile, stated))
		}
	}

	if !Carries(root, AttributeSettingsFile) {
		return append(disagreeing, fmt.Sprintf("the repository carries no %s, so a checkout on another platform can arrive with line endings of its own", AttributeSettingsFile)), nil
	}

	attributes, err := Document(root, AttributeSettingsFile)
	if err != nil {
		return nil, err
	}
	for _, stated := range attributeSettingsTheConventionsNeed {
		if !strings.Contains(attributes, stated) {
			disagreeing = append(disagreeing, fmt.Sprintf("%s states no %q", AttributeSettingsFile, stated))
		}
	}

	embedded, err := extensionsEmbedded(root)
	if err != nil {
		return nil, err
	}
	for _, extension := range embedded {
		if slices.Contains(textInTheBundle, extension) {
			continue
		}
		if !strings.Contains(attributes, "*"+extension+" binary") {
			disagreeing = append(disagreeing, fmt.Sprintf("%s embeds a %s and %s does not mark it binary, so a checkout can rewrite its bytes", GeneratedBundle, extension, AttributeSettingsFile))
		}
	}

	slices.Sort(disagreeing)
	return slices.Compact(disagreeing), nil
}

func extensionsEmbedded(root string) ([]string, error) {
	var found []string

	if !Carries(root, GeneratedBundle) {
		return nil, nil
	}

	err := filepath.WalkDir(filepath.Join(root, GeneratedBundle), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if extension := filepath.Ext(path); extension != "" && !slices.Contains(found, extension) {
			found = append(found, extension)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	slices.Sort(found)
	return found, nil
}
