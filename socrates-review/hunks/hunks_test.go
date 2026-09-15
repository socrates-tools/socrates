package hunks_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/socrates-review/hunks"
)

const oneChange = `diff --git a/internal/added.go b/internal/added.go
index 1234567..89abcde 100644
--- a/internal/added.go
+++ b/internal/added.go
@@ -1,4 +1,7 @@
 package internal

+func Answer(document string) int {
+	return len(document)
+}
+
 func Ask() {}
diff --git a/internal/gone.go b/internal/gone.go
deleted file mode 100644
index 89abcde..0000000
--- a/internal/gone.go
+++ /dev/null
@@ -1,3 +0,0 @@
-package internal
-
-func Gone() {}
`

func TestALineInsideAHunkIsInsideTheChange(t *testing.T) {
	change := hunks.Of(oneChange)

	for line := 1; line <= 7; line++ {
		if !change.Holds("internal/added.go", line) {
			t.Errorf("line %d of internal/added.go reads as outside the change, and the hunk covers lines 1 to 7", line)
		}
	}
	for _, outside := range []int{0, 8, 41} {
		if change.Holds("internal/added.go", outside) {
			t.Errorf("line %d of internal/added.go reads as inside the change, and the hunk covers lines 1 to 7", outside)
		}
	}
}

func TestAFileTheChangeNeverTouchedIsOutsideIt(t *testing.T) {
	change := hunks.Of(oneChange)

	if change.Holds("internal/untouched.go", 1) {
		t.Error("a file the diff never names reads as inside the change")
	}
	if change.Holds("internal/gone.go", 1) {
		t.Error("a deleted file reads as inside the change, and there is no line of it left to comment on")
	}
}

func TestASecondHunkOfTheSameFileIsRead(t *testing.T) {
	const twoHunks = `diff --git a/a.go b/a.go
--- a/a.go
+++ b/a.go
@@ -1,2 +1,3 @@
 one
+two
 three
@@ -40,2 +41,3 @@
 forty
+forty one
 forty two
`

	change := hunks.Of(twoHunks)

	want := map[string][]hunks.Span{"a.go": {{First: 1, Last: 3}, {First: 41, Last: 43}}}
	if diff := cmp.Diff(want, change.Touched); diff != "" {
		t.Errorf("the lines the change touches differ (-want +got):\n%s", diff)
	}
	for _, inside := range []int{2, 42} {
		if !change.Holds("a.go", inside) {
			t.Errorf("line %d reads as outside the change", inside)
		}
	}
	if change.Holds("a.go", 20) {
		t.Error("a line between two hunks reads as inside the change")
	}
}

func TestAHunkHeaderInsideAHunkIsNotAHunkHeader(t *testing.T) {
	// scenario: A diff of a diff does not move the change's lines
	const diffOfADiff = `diff --git a/recorded.diff b/recorded.diff
--- a/recorded.diff
+++ b/recorded.diff
@@ -1,3 +1,5 @@
+++ b/somewhere/else.go
+@@ -900,5 +900,9 @@
 kept
 also kept
 third
@@ -40,2 +42,3 @@
 forty
+forty one
 forty two
`

	change := hunks.Of(diffOfADiff)

	want := map[string][]hunks.Span{"recorded.diff": {{First: 1, Last: 5}, {First: 42, Last: 44}}}
	if diff := cmp.Diff(want, change.Touched); diff != "" {
		t.Errorf("the lines the change touches differ (-want +got):\n%s", diff)
	}
	if change.Holds("somewhere/else.go", 901) {
		t.Error("a path named on a line inside a hunk reads as a file the change touched")
	}
}

func TestAChangeWithASingleLineHunk(t *testing.T) {
	const oneLine = `diff --git a/a.go b/a.go
--- a/a.go
+++ b/a.go
@@ -1 +1 @@
-was
+is
`

	change := hunks.Of(oneLine)

	if !change.Holds("a.go", 1) {
		t.Error("the one line the change replaced reads as outside it")
	}
	if change.Holds("a.go", 2) {
		t.Error("the line after the one the change replaced reads as inside it")
	}
}

func TestAChangeSurvivesBeingWrittenAndReadBack(t *testing.T) {
	change := hunks.Of(oneChange)

	written, err := json.Marshal(change)
	if err != nil {
		t.Fatalf("writing what the change touches: %v", err)
	}

	var read hunks.Change
	if err := json.Unmarshal(written, &read); err != nil {
		t.Fatalf("reading what the change touches: %v", err)
	}
	if diff := cmp.Diff(change, read); diff != "" {
		t.Errorf("what was written differs from what was read back (-written +read):\n%s", diff)
	}
}

func TestNothingIsInsideAChangeNobodyRead(t *testing.T) {
	for _, unread := range []hunks.Change{{}, hunks.Of("")} {
		if unread.Holds("internal/added.go", 1) {
			t.Error("a line reads as inside a change nothing was read from")
		}
	}
}
