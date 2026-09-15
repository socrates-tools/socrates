package runstore

import (
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRunsSortByWhenTheyStarted(t *testing.T) {
	key, err := RepositoryKey(t.TempDir())
	if err != nil {
		t.Fatalf("deriving the repository key: %v", err)
	}

	started := time.Date(2026, 8, 19, 4, 15, 30, 0, time.UTC)
	inOrder := []string{
		Name(started, key),
		Name(started.Add(time.Second), key),
		Name(started.Add(time.Hour), key),
	}

	shuffled := []string{inOrder[2], inOrder[0], inOrder[1]}
	slices.Sort(shuffled)

	if diff := cmp.Diff(inOrder, shuffled); diff != "" {
		t.Errorf("names do not sort into the order they started (-started +sorted):\n%s", diff)
	}

	latest, found := Latest(shuffled, key)
	if !found || latest != inOrder[2] {
		t.Errorf("the latest run is %q, want %q", latest, inOrder[2])
	}
}

func TestTwoRepositoriesDoNotCollide(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "first")
	second := filepath.Join(root, "second")

	firstKey, err := RepositoryKey(first)
	if err != nil {
		t.Fatalf("deriving the first key: %v", err)
	}
	secondKey, err := RepositoryKey(second)
	if err != nil {
		t.Fatalf("deriving the second key: %v", err)
	}

	sameSecond := time.Date(2026, 8, 19, 4, 15, 30, 0, time.UTC)
	firstName := Name(sameSecond, firstKey)
	secondName := Name(sameSecond, secondKey)

	if firstName == secondName {
		t.Fatalf("two repositories started in the same second share the name %q", firstName)
	}

	names := []string{firstName, secondName}

	latest, found := Latest(names, firstKey)
	if !found || latest != firstName {
		t.Errorf("the first repository's run is %q, want %q", latest, firstName)
	}

	latest, found = Latest(names, secondKey)
	if !found || latest != secondName {
		t.Errorf("the second repository's run is %q, want %q", latest, secondName)
	}
}

func TestTheSamePathAlwaysGetsTheSameKey(t *testing.T) {
	root := t.TempDir()

	first, err := RepositoryKey(root)
	if err != nil {
		t.Fatalf("deriving the key: %v", err)
	}
	again, err := RepositoryKey(root)
	if err != nil {
		t.Fatalf("deriving the key again: %v", err)
	}
	if first != again {
		t.Errorf("one repository root produced %q and then %q", first, again)
	}

	unclean, err := RepositoryKey(filepath.Join(root, "elsewhere", ".."))
	if err != nil {
		t.Fatalf("deriving the key from an uncleaned path: %v", err)
	}
	if unclean != first {
		t.Errorf("the same root spelled differently produced %q, want %q", unclean, first)
	}

	other, err := RepositoryKey(filepath.Join(root, "other"))
	if err != nil {
		t.Fatalf("deriving the other key: %v", err)
	}
	if other == first {
		t.Errorf("two repository roots share the key %q", first)
	}
}

func TestLatestFindsNothingForAnUnknownRepository(t *testing.T) {
	key, err := RepositoryKey(t.TempDir())
	if err != nil {
		t.Fatalf("deriving the repository key: %v", err)
	}

	names := []string{Name(time.Now(), "aaaaaaaa"), Name(time.Now(), "bbbbbbbb")}
	if latest, found := Latest(names, key); found {
		t.Errorf("found %q for a repository with no runs", latest)
	}
}
