package convention

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	ContextMapFile     = "CONTEXT-MAP.md"
	GlossaryFile       = "CONTEXT.md"
	InstructionsFile   = "AGENTS.md"
	ClaudeInstructions = "CLAUDE.md"
	decisionsDirectory = "docs/adr"
)

var (
	recordName      = regexp.MustCompile(`^(\d{4})-.*\.md$`)
	recordTitle     = regexp.MustCompile(`^# (\d+)\. `)
	decisionNamed   = regexp.MustCompile(`\bADR\s+(\d+)\b`)
	linkTarget      = regexp.MustCompile(`\]\(([^)\s]+)\)`)
	linkedElsewhere = regexp.MustCompile(`^(?:[a-zA-Z][a-zA-Z0-9+.-]*:|#|/)`)
)

func DecisionDirectories(tools []Tool) []string {
	directories := []string{filepath.FromSlash(decisionsDirectory)}
	for _, tool := range tools {
		directories = append(directories, filepath.Join(tool.Directory, filepath.FromSlash(decisionsDirectory)))
	}
	return directories
}

func DecisionRecords(root string, tools []Tool) (map[string][]string, error) {
	records := map[string][]string{}
	for _, directory := range DecisionDirectories(tools) {
		entries, err := os.ReadDir(filepath.Join(root, directory))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			matched := recordName.FindStringSubmatch(entry.Name())
			if entry.IsDir() || matched == nil {
				continue
			}
			records[matched[1]] = append(records[matched[1]], filepath.Join(directory, entry.Name()))
		}
	}
	return records, nil
}

func NumbersHeldByMoreThanOneRecord(records map[string][]string) []string {
	var numbers []string
	for number, held := range records {
		if len(held) > 1 {
			numbers = append(numbers, number)
		}
	}
	slices.Sort(numbers)
	return numbers
}

func RecordsTitledForAnotherNumber(root string, tools []Tool) ([]string, error) {
	records, err := DecisionRecords(root, tools)
	if err != nil {
		return nil, err
	}

	var disagreeing []string
	for number, held := range records {
		named, _ := strconv.Atoi(number)
		for _, record := range held {
			text, err := os.ReadFile(filepath.Join(root, record))
			if err != nil {
				return nil, err
			}
			titled := recordTitle.FindStringSubmatch(firstLine(string(text)))
			if titled == nil {
				disagreeing = append(disagreeing, record+" has no title carrying its number")
				continue
			}
			if carried, _ := strconv.Atoi(titled[1]); carried != named {
				disagreeing = append(disagreeing, fmt.Sprintf("%s is titled for %d", record, carried))
			}
		}
	}
	slices.Sort(disagreeing)
	return disagreeing, nil
}

func ReferencesToNoDecision(root string, tools []Tool) ([]string, error) {
	records, err := DecisionRecords(root, tools)
	if err != nil {
		return nil, err
	}

	var unresolved []string
	err = walkDocuments(root, func(relative, text string) {
		for position, line := range strings.Split(text, "\n") {
			for _, named := range decisionNamed.FindAllStringSubmatch(line, -1) {
				number, _ := strconv.Atoi(named[1])
				if _, held := records[fmt.Sprintf("%04d", number)]; !held {
					unresolved = append(unresolved, fmt.Sprintf("%s:%d names ADR %d, which no directory holds", relative, position+1, number))
				}
			}
		}
	})
	if err != nil {
		return nil, err
	}
	slices.Sort(unresolved)
	return unresolved, nil
}

func LinksReachingNothing(root string) ([]string, error) {
	var broken []string
	err := walkDocuments(root, func(relative, text string) {
		fenced := false
		for position, line := range strings.Split(text, "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "```") {
				fenced = !fenced
				continue
			}
			if fenced {
				continue
			}
			for _, found := range linkTarget.FindAllStringSubmatch(line, -1) {
				target, _, _ := strings.Cut(found[1], "#")
				if target == "" || linkedElsewhere.MatchString(found[1]) {
					continue
				}
				if _, err := os.Stat(filepath.Join(root, filepath.Dir(relative), filepath.FromSlash(target))); err != nil {
					broken = append(broken, fmt.Sprintf("%s:%d links to %s, which does not exist", relative, position+1, target))
				}
			}
		}
	})
	if err != nil {
		return nil, err
	}
	slices.Sort(broken)
	return broken, nil
}

func ContextGaps(root string, tools []Tool) ([]string, error) {
	var gaps []string

	for _, shared := range []string{GlossaryFile, ContextMapFile, InstructionsFile} {
		if !Carries(root, shared) {
			gaps = append(gaps, shared+" does not exist")
		}
	}
	if !importsTheInstructionsAlone(root, ClaudeInstructions) {
		gaps = append(gaps, ClaudeInstructions+" does not hold @"+InstructionsFile+" alone")
	}

	held, err := Document(root, ContextMapFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	for _, tool := range tools {
		glossary := tool.Directory + "/" + GlossaryFile
		instructions := tool.Directory + "/" + InstructionsFile
		claude := tool.Directory + "/" + ClaudeInstructions
		decisions := tool.Directory + "/" + decisionsDirectory

		if !strings.Contains(held, glossary) {
			gaps = append(gaps, fmt.Sprintf("%s does not name %s", ContextMapFile, glossary))
		}
		if !Carries(root, glossary) {
			gaps = append(gaps, glossary+" does not exist")
		}
		if !Carries(root, instructions) {
			gaps = append(gaps, instructions+" does not exist")
		}
		if !importsTheInstructionsAlone(root, claude) {
			gaps = append(gaps, claude+" does not hold @"+InstructionsFile+" alone")
		}
		if !holdsRecords(filepath.Join(root, filepath.FromSlash(decisions))) {
			gaps = append(gaps, decisions+" holds no decision record")
		}
	}
	return gaps, nil
}

func importsTheInstructionsAlone(root, relative string) bool {
	path := filepath.Join(root, filepath.FromSlash(relative))
	if standing, err := os.Lstat(path); err != nil || !standing.Mode().IsRegular() {
		return false
	}
	held, err := os.ReadFile(path)
	return err == nil && string(held) == "@"+InstructionsFile+"\n"
}

func GlossaryNamingTheOtherTool(root, directory string) ([]string, error) {
	relative := filepath.Join(directory, GlossaryFile)
	held, err := os.ReadFile(filepath.Join(root, relative))
	if err != nil {
		return nil, err
	}

	var naming []string
	for position, line := range strings.Split(string(held), "\n") {
		for _, pattern := range otherToolPatterns {
			if found := pattern.FindString(line); found != "" {
				naming = append(naming, fmt.Sprintf("%s:%d names %q", relative, position+1, found))
			}
		}
	}
	return naming, nil
}

func holdsRecords(directory string) bool {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return false
	}
	return slices.ContainsFunc(entries, func(entry fs.DirEntry) bool {
		return !entry.IsDir() && recordName.MatchString(entry.Name())
	})
}

func walkDocuments(root string, read func(relative, text string)) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root && skipped(path, entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".md" || entry.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		held, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		read(relative, string(held))
		return nil
	})
}
