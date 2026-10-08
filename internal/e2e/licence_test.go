package e2e_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

var mitClauses = []string{
	"MIT License",
	"Permission is hereby granted, free of charge, to any person obtaining a copy",
	"to use, copy, modify, merge, publish, distribute, sublicense, and/or sell",
	"The above copyright notice and this permission notice shall be included in all",
	`THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND`,
	"IN NO EVENT SHALL THE",
}

var copyrightLine = regexp.MustCompile(`(?m)^Copyright \(c\) (\d{4}) (\S.*\S)$`)

func TestTheRootStatesTheTerms(t *testing.T) {
	t.Parallel()
	stated := rootDocument(t, convention.LicenceFile)

	for _, clause := range mitClauses {
		if !strings.Contains(stated, clause) {
			t.Errorf("%s states no %q, so it does not state the MIT licence in full", convention.LicenceFile, clause)
		}
	}

	held := copyrightLine.FindStringSubmatch(stated)
	if held == nil {
		t.Fatalf("%s carries no copyright line, so nobody is named as holding the terms", convention.LicenceFile)
	}
	if held[2] == "" {
		t.Errorf("%s names the year %s and no holder", convention.LicenceFile, held[1])
	}
}

func rootDocument(t *testing.T, named string) string {
	t.Helper()

	held, err := os.ReadFile(filepath.Join(source.Root(t), named))
	if err != nil {
		t.Fatalf("the root carries no %s: %v", named, err)
	}
	return string(held)
}
