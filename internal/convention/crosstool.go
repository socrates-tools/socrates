package convention

import (
	"path/filepath"
	"slices"
	"strings"
)

type CrossToolImport struct {
	File     string
	Tool     string
	Imported string
	Other    string
}

func CrossToolImports(root string, module string, tools []Tool) ([]CrossToolImport, error) {
	files, err := readTree(root)
	if err != nil {
		return nil, err
	}

	var crossing []CrossToolImport
	for _, file := range files {
		if file.test {
			continue
		}
		from := toolOwning(file.directory, tools)
		if from == "" {
			continue
		}

		for _, imported := range file.imported {
			beneath, inside := withinModule(imported, module)
			if !inside {
				continue
			}
			to := toolOwning(beneath, tools)
			if to == "" || to == from {
				continue
			}
			crossing = append(crossing, CrossToolImport{File: file.relative, Tool: from, Imported: imported, Other: to})
		}
	}

	slices.SortFunc(crossing, func(one, other CrossToolImport) int {
		return strings.Compare(one.File, other.File)
	})
	return crossing, nil
}

func toolOwning(directory string, tools []Tool) string {
	for _, tool := range tools {
		if directory == tool.Directory || strings.HasPrefix(directory, tool.Directory+string(filepath.Separator)) {
			return tool.Namespace
		}
	}
	return ""
}

func withinModule(imported, module string) (string, bool) {
	if module == "" {
		return "", false
	}
	if imported == module {
		return ".", true
	}
	beneath, inside := strings.CutPrefix(imported, module+"/")
	if !inside {
		return "", false
	}
	return filepath.FromSlash(beneath), true
}
