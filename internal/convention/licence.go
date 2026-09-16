package convention

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"
)

const (
	LicenceFile = "LICENSE"
	NoticeFile  = "NOTICE"
)

var GeneratedBundle = filepath.Join(ToolDirectory("agora"), "assets", "bundle")

var sourceExtensions = []string{".go", ".js", ".mjs", ".css", ".gohtml", ".yaml", ".yml", ".sh"}

var sourceNames = []string{"Makefile"}

var shebang = "#!"

var commentOpenings = []string{"//", "#", "/*", "<!--"}

var blockClosings = map[string]string{"/*": "*/", "<!--": "-->"}

var termsWording = []string{
	"spdx-license-identifier",
	"copyright (c)",
	"copyright ©",
	"all rights reserved",
	"licensed under",
	"permission is hereby granted",
	"this file is part of",
	"gnu general public license",
}

type LicenceHeader struct {
	File string
	Line int
	Text string
}

func LicenceHeaders(root string) ([]LicenceHeader, error) {
	var carried []LicenceHeader

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walked error) error {
		if walked != nil {
			return walked
		}

		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			relative = path
		}

		if entry.IsDir() {
			if path != root && (skipped(entry.Name()) || relative == GeneratedBundle) {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.Type().IsRegular() {
			return nil
		}

		held, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if !utf8.Valid(held) || !authoredSource(relative, string(held)) {
			return nil
		}

		opening, line := openingComment(string(held))
		if !statesTerms(opening) {
			return nil
		}

		carried = append(carried, LicenceHeader{File: relative, Line: line, Text: firstLine(opening)})
		return nil
	})
	if err != nil {
		return nil, err
	}

	slices.SortFunc(carried, func(one, other LicenceHeader) int {
		return strings.Compare(one.File, other.File)
	})
	return carried, nil
}

func authoredSource(relative, held string) bool {
	if slices.Contains(sourceExtensions, filepath.Ext(relative)) {
		return true
	}
	if slices.Contains(sourceNames, filepath.Base(relative)) {
		return true
	}
	return strings.HasPrefix(held, shebang)
}

func openingComment(held string) (string, int) {
	lines := strings.Split(held, "\n")

	position := blank(lines, 0)
	if position < len(lines) && strings.HasPrefix(lines[position], shebang) {
		position = blank(lines, position+1)
	}
	if position >= len(lines) {
		return "", 0
	}

	marker := commentOpening(lines[position])
	if marker == "" {
		return "", 0
	}

	var block []string
	if closing, blocked := blockClosings[marker]; blocked {
		for at := position; at < len(lines); at++ {
			block = append(block, lines[at])
			if strings.Contains(lines[at], closing) {
				break
			}
		}
		return strings.Join(block, "\n"), position + 1
	}

	for at := position; at < len(lines) && commentOpening(lines[at]) == marker; at++ {
		block = append(block, lines[at])
	}
	return strings.Join(block, "\n"), position + 1
}

func blank(lines []string, from int) int {
	for from < len(lines) && strings.TrimSpace(lines[from]) == "" {
		from++
	}
	return from
}

func commentOpening(line string) string {
	trimmed := strings.TrimSpace(line)
	for _, marker := range commentOpenings {
		if strings.HasPrefix(trimmed, marker) {
			return marker
		}
	}
	return ""
}

func statesTerms(opening string) bool {
	if strings.Contains(opening, NoticeFile) {
		return false
	}

	lowered := strings.ToLower(opening)
	for _, wording := range termsWording {
		if strings.Contains(lowered, wording) {
			return true
		}
	}
	return false
}
