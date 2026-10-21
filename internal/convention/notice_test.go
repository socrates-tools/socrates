package convention

import (
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
)

func TestAnEmbeddedAssetWithNoEntryFailsTheCheck(t *testing.T) {
	embedded := fstest.MapFS{
		"bundle/surface.js":                   &fstest.MapFile{Data: []byte("(()=>{})();")},
		"bundle/fonts/IBMPlexSerif-400.woff2": &fstest.MapFile{Data: []byte("a face somebody bundled")},
	}
	notice := "bundle/surface.js is the client, and it is socrates' own work.\n"

	unaccounted, err := EmbeddedFilesWithoutANotice(embedded, notice)
	if err != nil {
		t.Fatalf("reading what is embedded: %v", err)
	}

	want := []string{"bundle/fonts/IBMPlexSerif-400.woff2"}
	if diff := cmp.Diff(want, unaccounted); diff != "" {
		t.Errorf("the unaccounted files differ (-want +got):\n%s", diff)
	}
}

const aNoticeNamingALibraryOfEachShippedFile = "==========================\n" +
	"lit 3.3.3\n" +
	"Bundled into bundle/surface.js.\n" +
	"==========================\n\n" +
	"Copyright (c) 2017 Google LLC. All rights reserved.\n\n" +
	"==========================\n" +
	"tailwindcss 4.3.3\n" +
	"Bundled into bundle/surface.css.\n" +
	"==========================\n\n" +
	"Copyright (c) Tailwind Labs, Inc.\n"

func theBundle(js, css string) fstest.MapFS {
	return fstest.MapFS{
		"surface.js":  &fstest.MapFile{Data: []byte(js)},
		"surface.css": &fstest.MapFile{Data: []byte(css)},
	}
}

func TestABundlerThatDeletesNoticesFailsTheSuite(t *testing.T) {
	bundled := []string{"lit", "tailwindcss"}

	carried := theBundle("(()=>{})();\n/*! Copyright (c) 2017 Google LLC. All rights reserved. */\n", "/*! Copyright (c) Tailwind Labs, Inc. */\n.a{}\n")
	if missing := NoticesTheBundleDoesNotCarry(carried, aNoticeNamingALibraryOfEachShippedFile, bundled); len(missing) != 0 {
		t.Errorf("the check reported %v against a bundle that carries each notice", missing)
	}

	deleted := theBundle("(()=>{})();\n", ".a{}\n")
	want := []MissingNotice{
		{Library: "lit", File: "surface.js", Line: "Copyright (c) 2017 Google LLC. All rights reserved."},
		{Library: "tailwindcss", File: "surface.css", Line: "Copyright (c) Tailwind Labs, Inc."},
	}
	if diff := cmp.Diff(want, NoticesTheBundleDoesNotCarry(deleted, aNoticeNamingALibraryOfEachShippedFile, bundled)); diff != "" {
		t.Errorf("what the bundle does not carry differs (-want +got):\n%s", diff)
	}
}

func TestANoticeCarriedByTheWrongShippedFileDoesNotCountForTheFileTheLibraryShipsIn(t *testing.T) {
	bundled := []string{"lit", "tailwindcss"}

	swapped := theBundle("(()=>{})();\n/*! Copyright (c) Tailwind Labs, Inc. */\n", "/*! Copyright (c) 2017 Google LLC. All rights reserved. */\n.a{}\n")
	want := []MissingNotice{
		{Library: "lit", File: "surface.js", Line: "Copyright (c) 2017 Google LLC. All rights reserved."},
		{Library: "tailwindcss", File: "surface.css", Line: "Copyright (c) Tailwind Labs, Inc."},
	}
	if diff := cmp.Diff(want, NoticesTheBundleDoesNotCarry(swapped, aNoticeNamingALibraryOfEachShippedFile, bundled)); diff != "" {
		t.Errorf("what the bundle does not carry differs (-want +got):\n%s", diff)
	}
}

func TestALibraryWhoseNoticeNamesNoShippedFileFailsTheCheck(t *testing.T) {
	notice := "==========================\n" +
		"lit 3.3.3\n" +
		"BSD 3-Clause License.\n" +
		"==========================\n\n" +
		"Copyright (c) 2017 Google LLC. All rights reserved.\n"

	want := []MissingNotice{{Library: "lit", Line: "Copyright (c) 2017 Google LLC. All rights reserved."}}
	if diff := cmp.Diff(want, NoticesTheBundleDoesNotCarry(theBundle("", ""), notice, []string{"lit"})); diff != "" {
		t.Errorf("what the bundle does not carry differs (-want +got):\n%s", diff)
	}
}

func TestALibraryTheNoticeStatesNoCopyrightLineForFailsTheCheck(t *testing.T) {
	want := []MissingNotice{{Library: "left-pad"}}
	if diff := cmp.Diff(want, NoticesTheBundleDoesNotCarry(theBundle("", ""), aNoticeNamingALibraryOfEachShippedFile, []string{"left-pad"})); diff != "" {
		t.Errorf("what the bundle does not carry differs (-want +got):\n%s", diff)
	}
}

func TestAnAssetInTheSourceButNotInTheBinaryNeedsNoEntry(t *testing.T) {
	tree := t.TempDir()
	write(t, tree, filepath.Join("web", "fonts", "IBMPlexSerif-400.woff2"), "a face nobody bundled")

	embedded := fstest.MapFS{
		"bundle/surface.js": &fstest.MapFile{Data: []byte("(()=>{})();")},
	}
	notice := "bundle/surface.js is the client, and it is socrates' own work.\n"

	unaccounted, err := EmbeddedFilesWithoutANotice(embedded, notice)
	if err != nil {
		t.Fatalf("reading what is embedded: %v", err)
	}

	if len(unaccounted) != 0 {
		t.Errorf("the check asked for an entry for %v, and what a binary distributes is what a notice is a condition of", unaccounted)
	}
}

const aNoticeNamingTwoModules = "==========================\n" +
	"golang.org/x/sync v0.21.0\n" +
	"Linked into cmd/socrates.\n" +
	"==========================\n\n" +
	"Copyright 2009 The Go Authors.\n\n" +
	"==========================\n" +
	"github.com/spf13/cobra v1.10.2\n" +
	"Linked into cmd/socrates.\n" +
	"==========================\n\n" +
	"Copyright 2013 Steve Francia\n"

func TestAModuleTheBinaryLinksAndTheNoticeDoesNotNameFailsTheCheck(t *testing.T) {
	linked := []string{"github.com/spf13/cobra", "golang.org/x/net", "golang.org/x/sync"}

	want := []string{"golang.org/x/net"}
	if diff := cmp.Diff(want, ModulesWithoutANotice(linked, aNoticeNamingTwoModules)); diff != "" {
		t.Errorf("the modules without an entry differ (-want +got):\n%s", diff)
	}
}

func TestAModuleWhosePathOnlyPrefixesANamedOneIsNotNamed(t *testing.T) {
	linked := []string{"github.com/spf13/cob"}

	want := []string{"github.com/spf13/cob"}
	if diff := cmp.Diff(want, ModulesWithoutANotice(linked, aNoticeNamingTwoModules)); diff != "" {
		t.Errorf("the modules without an entry differ (-want +got):\n%s", diff)
	}
}

func TestEveryLinkedModuleNamedInTheNoticeLeavesTheCheckQuiet(t *testing.T) {
	linked := []string{"github.com/spf13/cobra", "golang.org/x/sync"}

	if missing := ModulesWithoutANotice(linked, aNoticeNamingTwoModules); len(missing) != 0 {
		t.Errorf("the check reported %v against a notice naming every linked module", missing)
	}
	if stale := ModulesNoticedButNotLinked(linked, aNoticeNamingTwoModules); len(stale) != 0 {
		t.Errorf("the check reported %v against a notice naming only linked modules", stale)
	}
}

func TestAModuleTheNoticeNamesAndTheBinaryNoLongerLinksFailsTheCheck(t *testing.T) {
	linked := []string{"golang.org/x/sync"}

	want := []string{"github.com/spf13/cobra"}
	if diff := cmp.Diff(want, ModulesNoticedButNotLinked(linked, aNoticeNamingTwoModules)); diff != "" {
		t.Errorf("the modules noticed and not linked differ (-want +got):\n%s", diff)
	}
}

func TestAModuleNoticedAtAVersionOtherThanTheLinkedOneFailsTheCheck(t *testing.T) {
	linked := map[string]string{
		"github.com/spf13/cobra": "v1.10.2",
		"golang.org/x/sync":      "v0.22.0",
	}

	want := []VersionDisagreement{{Module: "golang.org/x/sync", Linked: "v0.22.0", Noticed: "v0.21.0"}}
	if diff := cmp.Diff(want, ModulesNoticedAtAnotherVersion(linked, aNoticeNamingTwoModules)); diff != "" {
		t.Errorf("the modules noticed at another version differ (-want +got):\n%s", diff)
	}
}

func TestEveryModuleNoticedAtTheLinkedVersionLeavesTheVersionCheckQuiet(t *testing.T) {
	linked := map[string]string{
		"github.com/spf13/cobra": "v1.10.2",
		"golang.org/x/sync":      "v0.21.0",
	}

	if disagreeing := ModulesNoticedAtAnotherVersion(linked, aNoticeNamingTwoModules); len(disagreeing) != 0 {
		t.Errorf("the check reported %v against a notice naming every module at the version linked", disagreeing)
	}
}

func TestAModuleNobodyNoticesIsLeftToTheCheckForMissingEntries(t *testing.T) {
	linked := map[string]string{"golang.org/x/net": "v0.61.0"}

	if disagreeing := ModulesNoticedAtAnotherVersion(linked, aNoticeNamingTwoModules); len(disagreeing) != 0 {
		t.Errorf("the version check reported %v for a module the notice does not name", disagreeing)
	}
}

const aLicenceTheModuleCarries = "Copyright 2009 The Go Authors.\n\nRedistribution and use in source and binary forms.\n"

const aPatentGrantTheModuleCarries = "Additional IP Rights Grant (Patents)\n\nGoogle hereby grants to You a licence.\n"

func aNoticeCarrying(licence string) string {
	return "==========================\n" +
		"golang.org/x/sync v0.22.0\n" +
		"Linked into cmd/socrates.\n" +
		"==========================\n\n" +
		licence + "\n" +
		"==========================\n" +
		"github.com/spf13/cobra v1.10.2\n" +
		"Linked into cmd/socrates.\n" +
		"==========================\n\n" +
		"Copyright 2013 Steve Francia\n"
}

func aModuleDirectoryHolding(t *testing.T, files map[string]string) string {
	t.Helper()

	directory := t.TempDir()
	for name, content := range files {
		write(t, directory, name, content)
	}
	return directory
}

func TestALicenceTextTheModuleCarriesAndTheNoticeDoesNotFailsTheCheck(t *testing.T) {
	directory := aModuleDirectoryHolding(t, map[string]string{"LICENSE": "Copyright 2009 The Go Authors.\n\nA clause the notice dropped.\n"})

	disagreeing, err := ModulesWhoseLicenceTheNoticeDoesNotCarry(map[string]string{"golang.org/x/sync": directory}, aNoticeCarrying(aLicenceTheModuleCarries))
	if err != nil {
		t.Fatalf("reading the licence files: %v", err)
	}

	want := []LicenceDisagreement{{Module: "golang.org/x/sync", File: "LICENSE"}}
	if diff := cmp.Diff(want, disagreeing); diff != "" {
		t.Errorf("the licences the notice does not carry differ (-want +got):\n%s", diff)
	}
}

func TestEveryLicenceFileTheModuleCarriesIsCarriedByTheNotice(t *testing.T) {
	directory := aModuleDirectoryHolding(t, map[string]string{
		"LICENSE": aLicenceTheModuleCarries,
		"PATENTS": aPatentGrantTheModuleCarries,
	})
	carried := aNoticeCarrying(aLicenceTheModuleCarries + "\n" + aPatentGrantTheModuleCarries)

	disagreeing, err := ModulesWhoseLicenceTheNoticeDoesNotCarry(map[string]string{"golang.org/x/sync": directory}, carried)
	if err != nil {
		t.Fatalf("reading the licence files: %v", err)
	}
	if len(disagreeing) != 0 {
		t.Errorf("the check reported %v against a notice carrying every licence file", disagreeing)
	}
}

func TestAPatentGrantTheNoticeDropsFailsTheCheck(t *testing.T) {
	directory := aModuleDirectoryHolding(t, map[string]string{
		"LICENSE": aLicenceTheModuleCarries,
		"PATENTS": aPatentGrantTheModuleCarries,
	})

	disagreeing, err := ModulesWhoseLicenceTheNoticeDoesNotCarry(map[string]string{"golang.org/x/sync": directory}, aNoticeCarrying(aLicenceTheModuleCarries))
	if err != nil {
		t.Fatalf("reading the licence files: %v", err)
	}

	want := []LicenceDisagreement{{Module: "golang.org/x/sync", File: "PATENTS"}}
	if diff := cmp.Diff(want, disagreeing); diff != "" {
		t.Errorf("the licences the notice does not carry differ (-want +got):\n%s", diff)
	}
}

func TestALicenceCheckedAgainstAnotherModulesEntryFailsTheCheck(t *testing.T) {
	directory := aModuleDirectoryHolding(t, map[string]string{"LICENSE": "Copyright 2013 Steve Francia\n"})

	disagreeing, err := ModulesWhoseLicenceTheNoticeDoesNotCarry(map[string]string{"golang.org/x/sync": directory}, aNoticeCarrying(aLicenceTheModuleCarries))
	if err != nil {
		t.Fatalf("reading the licence files: %v", err)
	}

	want := []LicenceDisagreement{{Module: "golang.org/x/sync", File: "LICENSE"}}
	if diff := cmp.Diff(want, disagreeing); diff != "" {
		t.Errorf("the licences the notice does not carry differ (-want +got):\n%s", diff)
	}
}

func TestALicenceStatedWithWindowsLineEndingsIsTheSameLicence(t *testing.T) {
	directory := aModuleDirectoryHolding(t, map[string]string{"LICENSE": strings.ReplaceAll(aLicenceTheModuleCarries, "\n", "\r\n")})

	disagreeing, err := ModulesWhoseLicenceTheNoticeDoesNotCarry(map[string]string{"golang.org/x/sync": directory}, aNoticeCarrying(aLicenceTheModuleCarries))
	if err != nil {
		t.Fatalf("reading the licence files: %v", err)
	}
	if len(disagreeing) != 0 {
		t.Errorf("the check reported %v against the same text with other line endings", disagreeing)
	}
}

func TestAModuleHoldingNoLicenceFileFailsTheCheckRatherThanPassing(t *testing.T) {
	directory := aModuleDirectoryHolding(t, map[string]string{"README.md": "no licence here"})

	disagreeing, err := ModulesWhoseLicenceTheNoticeDoesNotCarry(map[string]string{"golang.org/x/sync": directory}, aNoticeCarrying(aLicenceTheModuleCarries))
	if err != nil {
		t.Fatalf("reading the licence files: %v", err)
	}

	want := []LicenceDisagreement{{Module: "golang.org/x/sync"}}
	if diff := cmp.Diff(want, disagreeing); diff != "" {
		t.Errorf("the licences the notice does not carry differ (-want +got):\n%s", diff)
	}
}

func TestAModuleTheNoticeDoesNotNameIsLeftToTheCheckForMissingEntries(t *testing.T) {
	directory := aModuleDirectoryHolding(t, map[string]string{"LICENSE": aLicenceTheModuleCarries})

	disagreeing, err := ModulesWhoseLicenceTheNoticeDoesNotCarry(map[string]string{"golang.org/x/net": directory}, aNoticeCarrying(aLicenceTheModuleCarries))
	if err != nil {
		t.Fatalf("reading the licence files: %v", err)
	}
	if len(disagreeing) != 0 {
		t.Errorf("the licence check reported %v for a module the notice does not name", disagreeing)
	}
}

func TestAWorkspacePackageTheClientDependsOnIsNotALibraryItBundlesFromElsewhere(t *testing.T) {
	tree := t.TempDir()
	write(t, tree, filepath.Join(ClientPackages, "surface", "package.json"), `{
  "name": "@agora/surface",
  "dependencies": {
    "@agora/ui": "workspace:*"
  }
}
`)
	write(t, tree, filepath.Join(ClientPackages, "ui", "package.json"), `{
  "name": "@agora/ui",
  "dependencies": {
    "highlight.js": "^11.11.1"
  }
}
`)

	bundled, err := LibrariesTheClientBundles(tree)
	if err != nil {
		t.Fatalf("reading what the client bundles: %v", err)
	}

	if diff := cmp.Diff([]string{"highlight.js"}, bundled); diff != "" {
		t.Errorf("the libraries the client bundles differ (-want +got):\n%s", diff)
	}
}

func TestALibraryAddedToAnyWorkspacePackageIsOneTheClientBundles(t *testing.T) {
	tree := t.TempDir()
	write(t, tree, filepath.Join(ClientPackages, "surface", "package.json"), `{
  "name": "@agora/surface",
  "dependencies": {
    "@agora/ui": "workspace:*",
    "a-fabricated-library": "^1.0.0"
  }
}
`)
	write(t, tree, filepath.Join(ClientPackages, "ui", "package.json"), `{
  "name": "@agora/ui",
  "dependencies": {
    "highlight.js": "^11.11.1"
  }
}
`)

	bundled, err := LibrariesTheClientBundles(tree)
	if err != nil {
		t.Fatalf("reading what the client bundles: %v", err)
	}

	if diff := cmp.Diff([]string{"a-fabricated-library", "highlight.js"}, bundled); diff != "" {
		t.Errorf("the libraries the client bundles differ (-want +got):\n%s", diff)
	}
}
