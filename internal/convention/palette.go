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

var StyleDirectory = filepath.Join("agora", "web", "src")

var paletteSelector = ":root"

var colourWritten = regexp.MustCompile(`(?i)#[0-9a-f]{3,8}\b|\b(?:rgba?|hsla?|hwb|lab|lch|oklab|oklch|color)\(|[:\s,(](?:white|black|silver|gainsboro|whitesmoke)[;\s,)]`)

var declaration = regexp.MustCompile(`^\s*([-a-zA-Z]+)\s*:`)

type Palette struct {
	Loose []string
}

func ColoursWrittenWhereTheyAreUsed(root, directory string) (Palette, error) {
	var read Palette

	err := filepath.WalkDir(filepath.Join(root, directory), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != filepath.Join(root, directory) && skipped(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".css" {
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

		read.Loose = append(read.Loose, coloursLooseIn(relative, string(held))...)
		return nil
	})
	if err != nil {
		return Palette{}, err
	}

	slices.Sort(read.Loose)
	return read, nil
}

func coloursLooseIn(named, held string) []string {
	var loose []string
	var open []string

	for position, line := range strings.Split(held, "\n") {
		selecting := strings.TrimSpace(line)
		if strings.HasSuffix(selecting, "{") {
			open = append(open, strings.TrimSpace(strings.TrimSuffix(selecting, "{")))
		}

		if painted := strings.Trim(colourWritten.FindString(line), " :,;()"); painted != "" {
			property := ""
			if named := declaration.FindStringSubmatch(line); named != nil {
				property = named[1]
			}
			switch {
			case !strings.HasPrefix(property, "--"):
				loose = append(loose, fmt.Sprintf("%s:%d paints %q where it is used", named, position+1, painted))
			case !slices.ContainsFunc(open, func(selector string) bool { return strings.Contains(selector, paletteSelector) }):
				loose = append(loose, fmt.Sprintf("%s:%d states %q outside the block that states the palette", named, position+1, painted))
			}
		}

		if strings.HasPrefix(selecting, "}") && len(open) > 0 {
			open = open[:len(open)-1]
		}
	}

	return loose
}
