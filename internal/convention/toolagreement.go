package convention

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const commandPackage = "command"

var suiteNamespaces = []string{"version", "skill", "help"}

func MountedWithoutATool(mounted []string, tools []Tool) []string {
	var unlisted []string
	for _, namespace := range mounted {
		if slices.Contains(suiteNamespaces, namespace) || listed(namespace, tools) {
			continue
		}
		unlisted = append(unlisted, namespace)
	}
	slices.Sort(unlisted)
	return unlisted
}

func ToolsNotMounted(mounted []string, tools []Tool) []string {
	var unmounted []string
	for _, tool := range tools {
		if !slices.Contains(mounted, tool.Namespace) {
			unmounted = append(unmounted, tool.Namespace)
		}
	}
	slices.Sort(unmounted)
	return unmounted
}

func ToolsWithoutACommandPackage(root string, tools []Tool) ([]string, error) {
	var lacking []string
	for _, tool := range tools {
		holds, err := holdsGoFiles(filepath.Join(root, tool.Directory, commandPackage))
		if err != nil {
			return nil, err
		}
		if !holds {
			lacking = append(lacking, tool.Namespace)
		}
	}
	slices.Sort(lacking)
	return lacking, nil
}

func CommandDirectoriesOutsideTheList(root string, tools []Tool) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var outside []string
	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())
		if !entry.IsDir() || skipped(path, entry.Name()) || listedDirectory(entry.Name(), tools) {
			continue
		}
		holds, err := holdsGoFiles(filepath.Join(path, commandPackage))
		if err != nil {
			return nil, err
		}
		if holds {
			outside = append(outside, entry.Name())
		}
	}
	return outside, nil
}

func listed(namespace string, tools []Tool) bool {
	return slices.ContainsFunc(tools, func(tool Tool) bool { return tool.Namespace == namespace })
}

func listedDirectory(directory string, tools []Tool) bool {
	return slices.ContainsFunc(tools, func(tool Tool) bool { return tool.Directory == directory })
}

func holdsGoFiles(directory string) (bool, error) {
	entries, err := os.ReadDir(directory)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			return true, nil
		}
	}
	return false, nil
}
