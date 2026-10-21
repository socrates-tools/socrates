package convention

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var SurfaceStyleDirectory = filepath.Join(ToolDirectory("agora"), "web", "packages", "surface")

var primitiveElementsStyled = regexp.MustCompile(`(?:^|[^-\w.#:\[(])(button|input|textarea)(?:[^-\w]|$)`)

var primitiveClassesStyled = regexp.MustCompile(`\.(copy|copy-status|expand|choice|decline|option-label|option-meaning|option-mark|words-field|words-label|message-field|message-label)(?:[^-\w]|$)`)

func PrimitivesStyledBySurface(root, directory string) ([]string, error) {
	var found []string

	err := filepath.WalkDir(filepath.Join(root, directory), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".css" {
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
		found = append(found, primitivesStyledIn(relative, string(held))...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return found, nil
}

func primitivesStyledIn(named, held string) []string {
	var found []string

	unannotated := commentWritten.ReplaceAllStringFunc(held, func(comment string) string {
		return notALineBreak.ReplaceAllString(comment, " ")
	})

	for position, line := range strings.Split(unannotated, "\n") {
		text := strings.TrimSpace(line)
		if !strings.HasSuffix(text, "{") && !strings.HasSuffix(text, ",") || strings.HasPrefix(text, "@") {
			continue
		}
		selector := strings.TrimSuffix(strings.TrimSuffix(text, "{"), ",")
		for _, written := range strings.Split(selector, ",") {
			written = strings.TrimSpace(written)
			if primitiveElementsStyled.MatchString(written) || primitiveClassesStyled.MatchString(written) {
				found = append(found, fmt.Sprintf("%s:%d styles %q, and the library's primitives style themselves", named, position+1, written))
			}
		}
	}
	return found
}
