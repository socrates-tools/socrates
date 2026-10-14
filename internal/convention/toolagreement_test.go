package convention_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/dispatcher"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func mountedNamespaces() []string {
	var code int
	var mounted []string
	for _, command := range dispatcher.NewRootCommand(&code).Commands() {
		mounted = append(mounted, command.Name())
	}
	return mounted
}

func TestEveryNamespaceTheDispatcherMountsIsAToolOrTheSuites(t *testing.T) {
	for _, namespace := range convention.MountedWithoutATool(mountedNamespaces(), convention.Tools) {
		t.Errorf("the dispatcher mounts %s, which is neither the suite's own nor on convention.Tools: no tool imports another tool, and that rule cannot see %s", namespace, namespace)
	}
}

func TestEveryToolOnTheListIsMountedByTheDispatcher(t *testing.T) {
	for _, namespace := range convention.ToolsNotMounted(mountedNamespaces(), convention.Tools) {
		t.Errorf("convention.Tools lists %s, which the dispatcher does not mount", namespace)
	}
}

func TestEveryToolOnTheListHoldsACommandPackage(t *testing.T) {
	lacking, err := convention.ToolsWithoutACommandPackage(source.Root(t), convention.Tools)
	if err != nil {
		t.Fatalf("reading the tree: %v", err)
	}
	for _, namespace := range lacking {
		t.Errorf("convention.Tools lists %s, whose directory holds no command package", namespace)
	}
}

func TestEveryTopLevelDirectoryHoldingACommandPackageIsAToolOnTheList(t *testing.T) {
	unlisted, err := convention.CommandDirectoriesOutsideTheList(source.Root(t), convention.Tools)
	if err != nil {
		t.Fatalf("reading the tree: %v", err)
	}
	for _, directory := range unlisted {
		t.Errorf("%s holds a command package and is not on convention.Tools: no tool imports another tool, and that rule cannot see %s", directory, directory)
	}
}

func TestAMountedNamespaceMissingFromTheListIsReported(t *testing.T) {
	tools := []convention.Tool{{Namespace: "review", Directory: "review"}}

	got := convention.MountedWithoutATool([]string{"review", "agora", "version", "help"}, tools)

	if diff := cmp.Diff([]string{"agora"}, got); diff != "" {
		t.Errorf("the namespaces mounted and not listed differ (-want +got):\n%s", diff)
	}
}

func TestTheSuitesOwnNamespacesNeedNoTool(t *testing.T) {
	got := convention.MountedWithoutATool([]string{"version", "skill", "help"}, nil)

	if len(got) != 0 {
		t.Errorf("the suite's own namespaces are reported as unlisted: %v", got)
	}
}

func TestAListedToolTheDispatcherDoesNotMountIsReported(t *testing.T) {
	tools := []convention.Tool{{Namespace: "review", Directory: "review"}, {Namespace: "agora", Directory: "agora"}}

	got := convention.ToolsNotMounted([]string{"review", "version"}, tools)

	if diff := cmp.Diff([]string{"agora"}, got); diff != "" {
		t.Errorf("the tools listed and not mounted differ (-want +got):\n%s", diff)
	}
}

func TestAListedToolWithNoDirectoryIsReported(t *testing.T) {
	root := t.TempDir()
	tools := []convention.Tool{{Namespace: "review", Directory: "review"}}

	got, err := convention.ToolsWithoutACommandPackage(root, tools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	if diff := cmp.Diff([]string{"review"}, got); diff != "" {
		t.Errorf("the tools without a command package differ (-want +got):\n%s", diff)
	}
}

func TestAListedToolWhoseDirectoryHoldsNoCommandPackageIsReported(t *testing.T) {
	root := t.TempDir()
	tools := []convention.Tool{{Namespace: "review", Directory: "review"}, {Namespace: "agora", Directory: "agora"}}
	mkdir(t, filepath.Join(root, "review", "config"))
	write(t, filepath.Join(root, "review", "config", "config.go"), "package config\n")
	mkdir(t, filepath.Join(root, "agora", "command"))
	write(t, filepath.Join(root, "agora", "command", "root.go"), "package command\n")

	got, err := convention.ToolsWithoutACommandPackage(root, tools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	if diff := cmp.Diff([]string{"review"}, got); diff != "" {
		t.Errorf("the tools without a command package differ (-want +got):\n%s", diff)
	}
}

func TestACommandDirectoryHoldingNoGoFileIsNotACommandPackage(t *testing.T) {
	root := t.TempDir()
	tools := []convention.Tool{{Namespace: "review", Directory: "review"}}
	mkdir(t, filepath.Join(root, "review", "command"))

	got, err := convention.ToolsWithoutACommandPackage(root, tools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	if diff := cmp.Diff([]string{"review"}, got); diff != "" {
		t.Errorf("the tools without a command package differ (-want +got):\n%s", diff)
	}
}

func TestATopLevelDirectoryHoldingACommandPackageAndMissingFromTheListIsReported(t *testing.T) {
	root := t.TempDir()
	tools := []convention.Tool{{Namespace: "review", Directory: "review"}}
	mkdir(t, filepath.Join(root, "review", "command"))
	write(t, filepath.Join(root, "review", "command", "root.go"), "package command\n")
	mkdir(t, filepath.Join(root, "scribe", "command"))
	write(t, filepath.Join(root, "scribe", "command", "root.go"), "package command\n")
	mkdir(t, filepath.Join(root, "cmd", "socrates"))
	write(t, filepath.Join(root, "cmd", "socrates", "main.go"), "package main\n")

	got, err := convention.CommandDirectoriesOutsideTheList(root, tools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	if diff := cmp.Diff([]string{"scribe"}, got); diff != "" {
		t.Errorf("the directories holding a command package outside the list differ (-want +got):\n%s", diff)
	}
}
