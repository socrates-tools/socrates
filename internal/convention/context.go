package convention

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var ContextFile = filepath.Join("openspec", "config.yaml")

const commandDirectory = "cmd"

var buildProducts = []string{"BINARY", "SURFACE", "SHORT_NAME"}

var commandNamed = regexp.MustCompile(`\b` + commandDirectory + `/([A-Za-z0-9._-]+)`)

var opaqueIdentifier = regexp.MustCompile(`\b[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\b`)

var makefileAssignment = regexp.MustCompile(`^([A-Z_]+)\s*:=\s*(\S+)`)

type ContextNaming struct {
	Absent  []string
	Unnamed []string
}

type ContextCitations struct {
	Unreachable              []string
	NamesOpenSpecAsNormative bool
}

func WhatTheNormativeContextNames(root string) (ContextNaming, error) {
	var naming ContextNaming

	context, err := normativeContext(root)
	if err != nil {
		return naming, err
	}

	for _, named := range commandNamed.FindAllStringSubmatch(context, -1) {
		if slices.Contains(naming.Absent, named[1]) {
			continue
		}
		if entry, err := os.Stat(filepath.Join(root, commandDirectory, named[1])); err != nil || !entry.IsDir() {
			naming.Absent = append(naming.Absent, named[1])
		}
	}
	slices.Sort(naming.Absent)

	built, err := binariesTheBuildProduces(root)
	if err != nil {
		return naming, err
	}
	for _, binary := range built {
		if !regexp.MustCompile(`\b` + regexp.QuoteMeta(binary) + `\b`).MatchString(context) {
			naming.Unnamed = append(naming.Unnamed, binary)
		}
	}
	slices.Sort(naming.Unnamed)

	return naming, nil
}

func WhatTheNormativeContextCites(root string) (ContextCitations, error) {
	var cited ContextCitations

	context, err := normativeContext(root)
	if err != nil {
		return cited, err
	}

	for position, line := range strings.Split(context, "\n") {
		for _, opaque := range opaqueIdentifier.FindAllString(line, -1) {
			cited.Unreachable = append(cited.Unreachable, fmt.Sprintf("%s:%d cites %s, which nobody outside the machine that made it can open", ContextFile, position+1, opaque))
		}
	}

	cited.NamesOpenSpecAsNormative = strings.Contains(context, "openspec/") && strings.Contains(context, "normative")
	return cited, nil
}

func normativeContext(root string) (string, error) {
	held, err := os.ReadFile(filepath.Join(root, ContextFile))
	if err != nil {
		return "", err
	}
	return string(held), nil
}

func binariesTheBuildProduces(root string) ([]string, error) {
	held, err := os.ReadFile(filepath.Join(root, "Makefile"))
	if err != nil {
		return nil, err
	}

	var built []string
	for _, line := range strings.Split(string(held), "\n") {
		assigned := makefileAssignment.FindStringSubmatch(line)
		if assigned == nil || !slices.Contains(buildProducts, assigned[1]) {
			continue
		}
		built = append(built, filepath.Base(assigned[2]))
	}
	slices.Sort(built)
	return slices.Compact(built), nil
}
