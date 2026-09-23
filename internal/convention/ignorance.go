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

var SurfaceDirectories = []string{
	filepath.Join(ToolDirectory("agora"), "command"),
	filepath.Join(ToolDirectory("agora"), "config"),
	filepath.Join(ToolDirectory("agora"), "assets"),
	filepath.Join(ToolDirectory("agora"), "document"),
	filepath.Join(ToolDirectory("agora"), "render"),
	filepath.Join(ToolDirectory("agora"), "sanitise"),
	filepath.Join(ToolDirectory("agora"), "stylesheet"),
	filepath.Join(ToolDirectory("agora"), "surface"),
	filepath.Join(ToolDirectory("agora"), "web", "src"),
}

var surfaceExtensions = []string{".go", ".gohtml", ".css", ".js", ".json"}

var generatedDirectories = []string{"bundle", "node_modules"}

var wordsOfTheOtherTool = []string{
	`(?i)\breviews?\b`,
	`(?i)\breviewers?\b`,
	`(?i)\bfindings?\b`,
	`(?i)\bverdicts?\b`,
	`(?i)\brounds?\b`,
	`(?i)\b(a|the|each|every|one|that|this)\s+runs?\b`,
	`(?i)\brun's\b`,
}

type Ignorance struct {
	Naming []string
}

func WordsOfTheOtherTool(root string, directories []string) (Ignorance, error) {
	var read Ignorance

	where, err := ModulePath(root)
	if err != nil {
		return Ignorance{}, err
	}

	patterns := make([]*regexp.Regexp, 0, len(wordsOfTheOtherTool))
	for _, stated := range wordsOfTheOtherTool {
		patterns = append(patterns, regexp.MustCompile(stated))
	}

	var own *regexp.Regexp
	if where != "" {
		own = regexp.MustCompile(regexp.QuoteMeta(where) + `(?:/[A-Za-z0-9_.-]+)?`)
	}

	for _, directory := range directories {
		err := filepath.WalkDir(filepath.Join(root, directory), func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				generated := slices.Contains(generatedDirectories, entry.Name())
				if path != filepath.Join(root, directory) && (generated || skipped(path, entry.Name())) {
					return filepath.SkipDir
				}
				return nil
			}
			if !slices.Contains(surfaceExtensions, filepath.Ext(path)) {
				return nil
			}

			held, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}

			relative, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}

			for position, line := range strings.Split(string(held), "\n") {
				vocabulary := line
				if own != nil {
					vocabulary = own.ReplaceAllString(line, "")
				}
				for _, pattern := range patterns {
					if named := pattern.FindString(vocabulary); named != "" {
						read.Naming = append(read.Naming, fmt.Sprintf("%s:%d names %q", relative, position+1, named))
					}
				}
			}
			return nil
		})
		if err != nil {
			return Ignorance{}, err
		}
	}

	slices.Sort(read.Naming)
	return read, nil
}

func PackagesOutside(root, tool string, listed []string) ([]string, error) {
	var unread []string
	err := filepath.WalkDir(filepath.Join(root, tool), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != filepath.Join(root, tool) && (slices.Contains(generatedDirectories, entry.Name()) || skipped(path, entry.Name())) {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		directory, relErr := filepath.Rel(root, filepath.Dir(path))
		if relErr != nil {
			return relErr
		}
		if !slices.Contains(unread, directory) && !slices.ContainsFunc(listed, covering(directory)) {
			unread = append(unread, directory)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	slices.Sort(unread)
	return unread, nil
}

func covering(directory string) func(string) bool {
	return func(listed string) bool {
		return directory == listed || strings.HasPrefix(directory, listed+string(filepath.Separator))
	}
}
