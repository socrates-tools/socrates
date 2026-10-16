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

func modulesNoticed(notice string) []string {
	var named []string

	for _, section := range noticeSeparator.Split(notice, -1) {
		lines := strings.Split(strings.TrimSpace(section), "\n")
		if len(lines) < 2 || !strings.HasPrefix(lines[1], "Linked into ") {
			continue
		}
		if fields := strings.Fields(lines[0]); len(fields) == 2 {
			named = append(named, fields[0])
		}
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
