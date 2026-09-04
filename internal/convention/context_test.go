package convention_test

import (
	"testing"

	"github.com/socrates-tools/socrates-review/internal/convention"
	"github.com/socrates-tools/socrates-review/internal/testing/source"
)

func TestTheContextNamesTheBinaryTheRepositoryBuilds(t *testing.T) {
	naming, err := convention.WhatTheNormativeContextNames(source.Root(t))
	if err != nil {
		t.Fatalf("reading the normative context: %v", err)
	}

	for _, absent := range naming.Absent {
		t.Errorf("%s names the command directory cmd/%s, and the repository has no such directory", convention.ContextFile, absent)
	}
	for _, unnamed := range naming.Unnamed {
		t.Errorf("%s names nothing called %s, and that is what the build produces", convention.ContextFile, unnamed)
	}
}

func TestTheContextCitesNothingAReaderCannotReach(t *testing.T) {
	cited, err := convention.WhatTheNormativeContextCites(source.Root(t))
	if err != nil {
		t.Fatalf("reading the normative context: %v", err)
	}

	for _, unreachable := range cited.Unreachable {
		t.Errorf("%s, and a citation that cannot be opened looks like evidence without being any", unreachable)
	}
	if !cited.NamesOpenSpecAsNormative {
		t.Errorf("%s does not state that openspec/ is what the project's design is normative in", convention.ContextFile)
	}
}
