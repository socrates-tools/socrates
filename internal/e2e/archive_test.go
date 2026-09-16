package e2e_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestAReleaseArchiveGivesWhatInstallingFromSourceGives(t *testing.T) {
	extracted := aReleaseArchiveExtracted(t)

	t.Run("An archive carries the terms and the notices", func(t *testing.T) {
		for _, stated := range []string{convention.LicenceFile, convention.NoticeFile} {
			held, err := os.ReadFile(filepath.Join(extracted, stated))
			if err != nil {
				t.Errorf("a release archive carries no %s, and an archive is a copy the terms travel with: %v", stated, err)
				continue
			}
			if len(bytes.TrimSpace(held)) == 0 {
				t.Errorf("the %s a release archive carries is empty", stated)
			}
		}
	})

	t.Run("An archive carries the tool", func(t *testing.T) {
		carried := []string{convention.Suite}
		for _, named := range carried {
			entry, err := os.Stat(filepath.Join(extracted, named))
			if err != nil {
				t.Errorf("a release archive carries no %s: %v", named, err)
				continue
			}
			if entry.Mode().Perm()&0o111 == 0 {
				t.Errorf("the %s a release archive carries is %s, which nobody can run", named, entry.Mode())
			}
		}

		entries, err := os.ReadDir(extracted)
		if err != nil {
			t.Fatalf("reading what a release archive carries: %v", err)
		}
		for _, entry := range entries {
			held, err := entry.Info()
			if err != nil || held.Mode().Perm()&0o111 == 0 || held.IsDir() {
				continue
			}
			if !slices.Contains(carried, entry.Name()) {
				t.Errorf("a release archive carries %s beside %v, and the tool answers to one name", entry.Name(), carried)
			}
		}
	})
}

func aReleaseArchiveExtracted(t *testing.T) string {
	t.Helper()

	repository := aCopyOfTheSource(t, requireMake(t))
	distribution := filepath.Join(repository, "dist")

	command := exec.Command("make", "release-snapshot")
	command.Dir = repository
	command.Env = append(os.Environ(), source.Toolchain()...)

	var failures bytes.Buffer
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("make release-snapshot into %s: %v: %s", distribution, err, failures.String())
	}

	archives, err := filepath.Glob(filepath.Join(distribution, "*_"+runtime.GOOS+"_"+runtime.GOARCH+".tar.gz"))
	if err != nil {
		t.Fatalf("reading what the release produced: %v", err)
	}
	if len(archives) != 1 {
		t.Fatalf("a release produced %d archives for %s_%s, and this platform is one the supported set names", len(archives), runtime.GOOS, runtime.GOARCH)
	}

	extracted := t.TempDir()
	if err := unpack(archives[0], extracted); err != nil {
		t.Fatalf("extracting %s: %v", archives[0], err)
	}
	return extracted
}

func unpack(archive, into string) error {
	held, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer func() { _ = held.Close() }()

	uncompressed, err := gzip.NewReader(held)
	if err != nil {
		return err
	}
	defer func() { _ = uncompressed.Close() }()

	carried := tar.NewReader(uncompressed)
	for {
		entry, err := carried.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if entry.Typeflag != tar.TypeReg {
			continue
		}

		name := filepath.Clean(filepath.FromSlash(entry.Name))
		if filepath.IsAbs(name) || strings.HasPrefix(name, "..") {
			return fmt.Errorf("the archive carries %s, which reaches outside where it is extracted", entry.Name)
		}

		written := filepath.Join(into, name)
		if err := os.MkdirAll(filepath.Dir(written), 0o755); err != nil {
			return err
		}
		file, err := os.OpenFile(written, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, entry.FileInfo().Mode().Perm())
		if err != nil {
			return err
		}
		if _, err := io.Copy(file, carried); err != nil {
			_ = file.Close()
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
	}
}
