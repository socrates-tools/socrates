package convention

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var ClientManifest = filepath.Join(ToolDirectory("agora"), "web", "package.json")

var noticeSeparator = regexp.MustCompile(`(?m)^={20,}$`)

var copyrightStated = regexp.MustCompile(`(?m)^Copyright .+$`)

type MissingNotice struct {
	Library string
	Line    string
}

func NoticesTheClientDoesNotCarry(client, notice string, bundled []string) []MissingNotice {
	var missing []MissingNotice

	for _, library := range bundled {
		stated := copyrightStatedFor(notice, library)
		if len(stated) == 0 {
			missing = append(missing, MissingNotice{Library: library})
			continue
		}

		for _, line := range stated {
			if !strings.Contains(client, line) {
				missing = append(missing, MissingNotice{Library: library, Line: line})
			}
		}
	}
	return missing
}

func copyrightStatedFor(notice, library string) []string {
	sections := noticeSeparator.Split(notice, -1)

	for position, section := range sections {
		if !strings.HasPrefix(strings.TrimSpace(section), library) || position+1 >= len(sections) {
			continue
		}
		return copyrightStated.FindAllString(sections[position+1], -1)
	}
	return nil
}

func LibrariesTheClientBundles(root string) ([]string, error) {
	held, err := os.ReadFile(filepath.Join(root, ClientManifest))
	if err != nil {
		return nil, err
	}

	var declared struct {
		Dependencies map[string]string `json:"dependencies"`
	}
	if err := json.Unmarshal(held, &declared); err != nil {
		return nil, fmt.Errorf("reading %s: %w", ClientManifest, err)
	}

	bundled := make([]string, 0, len(declared.Dependencies))
	for library := range declared.Dependencies {
		bundled = append(bundled, library)
	}
	slices.Sort(bundled)
	return bundled, nil
}

func EmbeddedFilesWithoutANotice(embedded fs.FS, notice string) ([]string, error) {
	var unaccounted []string

	err := fs.WalkDir(embedded, ".", func(path string, entry fs.DirEntry, walked error) error {
		if walked != nil {
			return walked
		}
		if entry.IsDir() {
			return nil
		}
		if !strings.Contains(notice, filepath.ToSlash(path)) {
			unaccounted = append(unaccounted, filepath.ToSlash(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	slices.Sort(unaccounted)
	return unaccounted, nil
}

type noticeEntry struct {
	Module  string
	Version string
	Text    string
}

func noticeEntries(notice string) []noticeEntry {
	var entries []noticeEntry

	sections := noticeSeparator.Split(notice, -1)
	for position, section := range sections {
		lines := strings.Split(strings.TrimSpace(section), "\n")
		if len(lines) < 2 || !strings.HasPrefix(lines[1], "Linked into ") || position+1 >= len(sections) {
			continue
		}
		if fields := strings.Fields(lines[0]); len(fields) == 2 {
			entries = append(entries, noticeEntry{Module: fields[0], Version: fields[1], Text: sections[position+1]})
		}
	}
	return entries
}

func modulesNoticed(notice string) []string {
	var named []string
	for _, entry := range noticeEntries(notice) {
		named = append(named, entry.Module)
	}
	return named
}

func ModulesWithoutANotice(linked []string, notice string) []string {
	noticed := modulesNoticed(notice)

	var missing []string
	for _, module := range linked {
		if !slices.Contains(noticed, module) {
			missing = append(missing, module)
		}
	}
	slices.Sort(missing)
	return missing
}

func ModulesNoticedButNotLinked(linked []string, notice string) []string {
	var stale []string
	for _, module := range modulesNoticed(notice) {
		if !slices.Contains(linked, module) {
			stale = append(stale, module)
		}
	}
	slices.Sort(stale)
	return stale
}

type VersionDisagreement struct {
	Module  string
	Linked  string
	Noticed string
}

func ModulesNoticedAtAnotherVersion(linked map[string]string, notice string) []VersionDisagreement {
	var disagreeing []VersionDisagreement
	for _, entry := range noticeEntries(notice) {
		version, isLinked := linked[entry.Module]
		if isLinked && version != entry.Version {
			disagreeing = append(disagreeing, VersionDisagreement{Module: entry.Module, Linked: version, Noticed: entry.Version})
		}
	}
	slices.SortFunc(disagreeing, func(a, b VersionDisagreement) int { return strings.Compare(a.Module, b.Module) })
	return disagreeing
}

type LicenceDisagreement struct {
	Module string
	File   string
}

var licenceFileName = regexp.MustCompile(`(?i)^(licen[sc]e|patents|notice|copying)(\..*)?$`)

func ModulesWhoseLicenceTheNoticeDoesNotCarry(directories map[string]string, notice string) ([]LicenceDisagreement, error) {
	var disagreeing []LicenceDisagreement

	for _, entry := range noticeEntries(notice) {
		directory, isLinked := directories[entry.Module]
		if !isLinked {
			continue
		}

		files, err := licenceFilesIn(directory)
		if err != nil {
			return nil, fmt.Errorf("reading the licence files of %s: %w", entry.Module, err)
		}
		if len(files) == 0 {
			disagreeing = append(disagreeing, LicenceDisagreement{Module: entry.Module})
			continue
		}

		carried := withUnixLineEndings(entry.Text)
		for _, file := range files {
			held, err := os.ReadFile(filepath.Join(directory, file))
			if err != nil {
				return nil, fmt.Errorf("reading %s of %s: %w", file, entry.Module, err)
			}
			if !strings.Contains(carried, strings.TrimSpace(withUnixLineEndings(string(held)))) {
				disagreeing = append(disagreeing, LicenceDisagreement{Module: entry.Module, File: file})
			}
		}
	}

	slices.SortFunc(disagreeing, func(a, b LicenceDisagreement) int {
		return strings.Compare(a.Module+"\x00"+a.File, b.Module+"\x00"+b.File)
	})
	return disagreeing, nil
}

func licenceFilesIn(directory string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.Type().IsRegular() && licenceFileName.MatchString(entry.Name()) {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func withUnixLineEndings(text string) string {
	return strings.ReplaceAll(text, "\r\n", "\n")
}
