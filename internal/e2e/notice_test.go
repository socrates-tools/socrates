package e2e_test

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/assets"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func openFontLicence(copyright string) []string {
	return []string{
		copyright,
		"SIL Open Font License, Version 1.1",
		"SIL OPEN FONT LICENSE Version 1.1",
		`THE FONT SOFTWARE IS PROVIDED "AS IS"`,
		"PERMISSION & CONDITIONS",
	}
}

var licenceOfAFontFamily = map[string][]string{
	"Cinzel":         openFontLicence("Copyright 2020 The Cinzel Project Authors (https://github.com/NDISCOVER/Cinzel)"),
	"Manrope":        openFontLicence("Copyright 2019 The Manrope Project Authors (https://github.com/sharanda/manrope)"),
	"JetBrains Mono": openFontLicence("Copyright 2020 The JetBrains Mono Project Authors (https://github.com/JetBrains/JetBrainsMono)"),
}

var licenceOfABundledLibrary = map[string][]string{
	"highlight.js": {
		"BSD 3-Clause License",
		"Copyright (c) 2006, Ivan Sagalaev.",
		"Redistribution and use in source and binary forms, with or without",
		`THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"`,
	},
}

var declaredFamily = regexp.MustCompile(`@font-face\s*\{[^}]*font-family\s*:\s*['"]?([^;'"}]+)`)

func TestAnEmbeddedFontsLicenceTravelsWithTheWork(t *testing.T) {
	notice := rootDocument(t, convention.NoticeFile)

	faces, err := fs.Glob(assets.Files(), filepath.ToSlash(filepath.Join("fonts", "*.woff2")))
	if err != nil {
		t.Fatalf("reading what the bundle embeds: %v", err)
	}
	if len(faces) == 0 {
		t.Fatalf("the bundle embeds no font face, so this test proves nothing")
	}

	for _, face := range faces {
		if !strings.Contains(notice, filepath.Base(face)) {
			t.Errorf("%s names no %s, and that face ships inside the binary", convention.NoticeFile, filepath.Base(face))
		}
	}

	for _, family := range familiesTheStylesheetDeclares(t) {
		if !strings.Contains(notice, family) {
			t.Errorf("%s names no %q, and the embedded stylesheet asks for it", convention.NoticeFile, family)
		}

		clauses, stated := licenceOfAFontFamily[family]
		if !stated {
			t.Errorf("%q is embedded and this test states no licence text for it: state it here and carry it in %s", family, convention.NoticeFile)
			continue
		}
		for _, clause := range clauses {
			if !strings.Contains(notice, clause) {
				t.Errorf("%s carries no %q, so the licence of %q is referred to rather than carried", convention.NoticeFile, clause, family)
			}
		}
	}
}

func TestAnEmbeddedLibrarysLicenceTravelsWithTheWork(t *testing.T) {
	notice := rootDocument(t, convention.NoticeFile)

	for _, library := range librariesTheClientBundles(t) {
		if !strings.Contains(notice, library) {
			t.Errorf("%s names no %s, and it is bundled into the client that ships", convention.NoticeFile, library)
		}

		clauses, stated := licenceOfABundledLibrary[library]
		if !stated {
			t.Errorf("%s is bundled and this test states no licence text for it: state it here and carry it in %s", library, convention.NoticeFile)
			continue
		}
		for _, clause := range clauses {
			if !strings.Contains(notice, clause) {
				t.Errorf("%s carries no %q, so the licence of %s is referred to rather than carried", convention.NoticeFile, clause, library)
			}
		}
	}
}

func familiesTheStylesheetDeclares(t *testing.T) []string {
	t.Helper()

	held, err := fs.ReadFile(assets.Files(), "surface.css")
	if err != nil {
		t.Fatalf("reading the embedded stylesheet: %v", err)
	}

	var declared []string
	for _, face := range declaredFamily.FindAllStringSubmatch(string(held), -1) {
		family := strings.TrimSpace(face[1])
		if !slices.Contains(declared, family) {
			declared = append(declared, family)
		}
	}
	slices.Sort(declared)

	if len(declared) == 0 {
		t.Fatalf("the embedded stylesheet declares no font family, so this test proves nothing")
	}
	return declared
}

func librariesTheClientBundles(t *testing.T) []string {
	t.Helper()

	bundled, err := convention.LibrariesTheClientBundles(source.Root(t))
	if err != nil {
		t.Fatalf("reading what the client bundles: %v", err)
	}
	if len(bundled) == 0 {
		t.Fatalf("the client bundles nothing, so this test proves nothing")
	}
	return bundled
}
