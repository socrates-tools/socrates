package convention

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ktei/socrates/internal/testing/source"
)

func TestNoSourceFileCarriesALicenceHeader(t *testing.T) {
	carried, err := LicenceHeaders(source.Root(t))
	if err != nil {
		t.Fatalf("reading the source: %v", err)
	}

	for _, header := range carried {
		t.Errorf("%s:%d opens with a licence header, and the terms are stated once at the root: %s", header.File, header.Line, header.Text)
	}
}

func TestALicenceHeaderIsReported(t *testing.T) {
	cases := []struct {
		name     string
		file     string
		content  string
		reported []LicenceHeader
	}{
		{
			name:     "an identifier above the package clause is reported",
			file:     filepath.Join("internal", "thing", "thing.go"),
			content:  "// SPDX-License-Identifier: MIT\n\npackage thing\n",
			reported: []LicenceHeader{{Line: 1, Text: "// SPDX-License-Identifier: MIT"}},
		},
		{
			name:     "a copyright block above the package clause is reported",
			file:     filepath.Join("internal", "thing", "thing.go"),
			content:  "/*\nCopyright (c) 2026 Somebody Else\nLicensed under the terms below.\n*/\n\npackage thing\n",
			reported: []LicenceHeader{{Line: 1, Text: "/*"}},
		},
		{
			name:     "a header under a shebang is reported",
			file:     filepath.Join("scripts", "run-it"),
			content:  "#!/bin/sh\n# Copyright (c) 2026 Somebody Else\nset -e\n",
			reported: []LicenceHeader{{Line: 2, Text: "# Copyright (c) 2026 Somebody Else"}},
		},
		{
			name:     "a header in a stylesheet is reported",
			file:     filepath.Join("web", "src", "thing.css"),
			content:  "/*! highlight.js, all rights reserved */\n.thing { color: red; }\n",
			reported: []LicenceHeader{{Line: 1, Text: "/*! highlight.js, all rights reserved */"}},
		},
		{
			// scenario: A borrowed work's notice is not the header this forbids
			name:    "a borrowed work's notice is not the header this forbids",
			file:    filepath.Join("web", "src", "surface.js"),
			content: "/*! highlight.js is bundled here.\n    BSD 3-Clause License.\n    Copyright (c) 2006, Ivan Sagalaev.\n    The full text is in " + NoticeFile + " at the root. */\nimport hljs from 'highlight.js/lib/core'\n",
		},
		{
			name:    "a build directive is not reported",
			file:    filepath.Join("internal", "thing", "thing.go"),
			content: "//go:build linux\n\npackage thing\n",
		},
		{
			name:    "terms further down a file are not reported",
			file:    filepath.Join("internal", "thing", "thing.go"),
			content: "package thing\n\n// Copyright (c) 2026 Somebody Else\nfunc Do() {}\n",
		},
		{
			name:    "a document naming the terms is not reported",
			file:    "README.md",
			content: "# thing\n\nCopyright (c) 2026 Somebody Else, licensed under MIT.\n",
		},
		{
			name:    "the generated bundle is not read",
			file:    filepath.Join(GeneratedBundle, "surface.js"),
			content: "/*! Copyright (c) 2006 Ivan Sagalaev, all rights reserved */\n(()=>{})();\n",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			root := t.TempDir()
			write(t, root, testCase.file, testCase.content)

			carried, err := LicenceHeaders(root)
			if err != nil {
				t.Fatalf("reading the fabricated source: %v", err)
			}

			want := testCase.reported
			for position := range want {
				want[position].File = testCase.file
			}
			if diff := cmp.Diff(want, carried); diff != "" {
				t.Errorf("the reported licence headers differ (-want +got):\n%s", diff)
			}
		})
	}
}
