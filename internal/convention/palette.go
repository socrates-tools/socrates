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

var StyleDirectory = filepath.Join(ToolDirectory("agora"), "web")

var LandingPageStyleDirectory = "site"

var TokensDirectory = filepath.Join(ToolDirectory("agora"), "web", "packages", "ui", "tokens")

const GeneratedStylesheetSuffix = ".generated.css"

var paletteSelector = ":root"

var colourWritten = regexp.MustCompile(`(?i)#[0-9a-f]{3,8}\b|\b(?:rgba?|hsla?|hwb|lab|lch|oklab|oklch|color)\(|[:\s,(](?:white|black|silver|gainsboro|whitesmoke)[;\s,)]`)

var declaration = regexp.MustCompile(`^\s*([-a-zA-Z0-9_]+)\s*:`)

var fontDeclared = regexp.MustCompile(`(?:^|[\s;{])(font-family|font)\s*:\s*([^;{}]*)`)

var commentWritten = regexp.MustCompile(`(?s)/\*.*?\*/`)

var notALineBreak = regexp.MustCompile(`[^\n]`)

var fontLeftToTheTokens = regexp.MustCompile(`^(?:inherit|initial|unset|revert|revert-layer|var\(--font-[a-z-]+\))$`)

var fontTakenFromTheTokens = regexp.MustCompile(`var\(--font-[a-z-]+\)`)

type Painting struct {
	Colours  []string
	Families []string
}

func PaintWrittenWhereItIsUsed(root, directory, tokens string) (Painting, error) {
	var read Painting

	err := filepath.WalkDir(filepath.Join(root, directory), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != filepath.Join(root, directory) && skipped(path, entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".css" || strings.HasSuffix(path, GeneratedStylesheetSuffix) {
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

		if isWithin(relative, tokens) {
			read.Colours = append(read.Colours, coloursLooseIn(relative, string(held), true)...)
		} else {
			read.Colours = append(read.Colours, coloursLooseIn(relative, string(held), false)...)
			read.Families = append(read.Families, familiesNamedIn(relative, string(held))...)
		}
		return nil
	})
	if err != nil {
		return Painting{}, err
	}

	slices.Sort(read.Colours)
	slices.Sort(read.Families)
	return read, nil
}

func isWithin(relative, directory string) bool {
	within, err := filepath.Rel(directory, relative)
	return err == nil && within != ".." && !strings.HasPrefix(within, ".."+string(filepath.Separator))
}

func familiesNamedIn(named, held string) []string {
	var loose []string

	unannotated := commentWritten.ReplaceAllStringFunc(held, func(comment string) string {
		return notALineBreak.ReplaceAllString(comment, " ")
	})

	for _, declared := range fontDeclared.FindAllStringSubmatchIndex(unannotated, -1) {
		property := unannotated[declared[2]:declared[3]]
		value := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(unannotated[declared[4]:declared[5]]), "!important"))
		value = strings.Join(strings.Fields(value), " ")
		if fontLeftToTheTokens.MatchString(value) {
			continue
		}
		if property == "font" && fontTakenFromTheTokens.MatchString(value) {
			continue
		}
		position := strings.Count(unannotated[:declared[2]], "\n") + 1
		loose = append(loose, fmt.Sprintf("%s:%d names the font family %s outside the library's tokens", named, position, value))
	}

	return loose
}

func coloursLooseIn(named, held string, stated bool) []string {
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
			case !stated:
				loose = append(loose, fmt.Sprintf("%s:%d states %q outside the library's tokens", named, position+1, painted))
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
