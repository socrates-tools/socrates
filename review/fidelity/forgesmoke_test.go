//go:build forgesmoke

package fidelity_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/exit"
	"github.com/socrates-tools/socrates/internal/state"
	"github.com/socrates-tools/socrates/internal/testing/simulator"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/dialogue"
	"github.com/socrates-tools/socrates/review/snapshot"
)

const (
	theRepositoryUnderReview = "socrates-tools/repository-under-review"
	theFixturePullRequest    = 1
)

const (
	whileTheRoundIsParked        = exit.Working
	onceTheFindingsAreHandedOver = exit.SomethingAsked
)

const whatAReviewerFoundInTheFixture = `{"findings":[` +
	`{"id":"f1","severity":"error","file":"settlement/rate.go","line":6,"description":"Rounding by adding half a basis point rounds a negative amount away from zero, so a refund of one and a half cents settles as two.","action":"auto-fix"},` +
	`{"id":"f2","severity":"error","file":"settlement/ledger.go","line":22,"description":"Statement keeps its own truncating copy of the conversion, so a statement prints line items that do not sum to the total it prints beside them.","action":"ask-user"},` +
	`{"id":"f3","severity":"warning","file":"settlement/ledger.go","line":23,"description":"Statement concatenates in a loop, which reallocates on every entry.","action":"auto-fix"}` +
	`],"risk_level":"medium","risk_rationale":"one changed line answers wrongly for a negative amount without erroring"}`

var insideTheChange = map[string]bool{"f1": true, "f2": false, "f3": false}

type reportedFinding struct {
	ID          string              `json:"id"`
	File        string              `json:"file"`
	Line        int                 `json:"line"`
	Description string              `json:"description"`
	Resolution  dialogue.Resolution `json:"resolution"`
	Inside      *bool               `json:"inside_the_change"`
}

type reportedRound struct {
	Findings []reportedFinding `json:"findings"`
}

type pullRequest struct {
	Head struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"base"`
	State string `json:"state"`
}

type postedComment struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	Path     string `json:"path"`
	Line     int    `json:"line"`
	Position *int   `json:"position"`
	HTMLURL  string `json:"html_url"`
}

func TestAForgePutsAForwardedFindingWhereSocratesSaidItGoes(t *testing.T) {
	pull := theFixture(t)
	if pull.State != "open" {
		t.Fatalf("pull request %d of %s is %s, and the smoke test posts to the one that stays open", theFixturePullRequest, theRepositoryUnderReview, pull.State)
	}

	clone := theirChangeCheckedOut(t, pull)
	binary := socratesReview(t)
	simulated := simulator.InstallFor(t, agent.TheDefault())
	simulated.Scenario("default:\n  structured: '" + whatAReviewerFoundInTheFixture + "'\n")

	root := t.TempDir()
	t.Setenv(state.HomeVariable, root)
	t.Cleanup(func() { _ = snapshot.Unseal(state.HomeUnder(root)) })

	name := reviewed(t, binary, clone, pull)
	reported := whatSocratesReported(t, binary, clone, name, whileTheRoundIsParked)

	for _, finding := range reported.Findings {
		if finding.Inside == nil {
			t.Fatalf("socrates-review says nothing about whether %s is inside the change, and where its comment goes is read from that", finding.ID)
		}
		if *finding.Inside != insideTheChange[finding.ID] {
			t.Fatalf("socrates-review reports %s at %s:%d as inside the change: %t, and the pull request's diff holds that line: %t",
				finding.ID, finding.File, finding.Line, *finding.Inside, insideTheChange[finding.ID])
		}
	}

	for _, finding := range reported.Findings {
		socrates(t, binary, clone, whileTheRoundIsParked, "respond", name, "--finding", finding.ID, "--forwarded")
	}
	socrates(t, binary, clone, onceTheFindingsAreHandedOver, "continue", name)

	forwarded := whatSocratesReported(t, binary, clone, name, onceTheFindingsAreHandedOver)
	var inside, outside []reportedFinding
	for _, finding := range forwarded.Findings {
		if finding.Resolution != dialogue.Forwarded {
			t.Fatalf("%s reads as %q rather than forwarded, and only what was forwarded is posted", finding.ID, finding.Resolution)
		}
		if *finding.Inside {
			inside = append(inside, finding)
			continue
		}
		outside = append(outside, finding)
	}
	if len(inside) == 0 || len(outside) < 2 {
		t.Fatalf("the fixture reported %d findings inside the change and %d outside it, and both halves of the rule need exercising", len(inside), len(outside))
	}

	t.Run("a finding inside the change lands on its own line", func(t *testing.T) {
		for _, finding := range inside {
			posted := commentOnALine(t, pull, finding)
			t.Cleanup(func() { deleteReviewComment(t, posted.ID) })

			if posted.Path != finding.File || posted.Line != finding.Line {
				t.Errorf("%s was forwarded for %s:%d and the forge put it at %s:%d", finding.ID, finding.File, finding.Line, posted.Path, posted.Line)
			}
			if posted.Position == nil {
				t.Errorf("%s landed on %s:%d and the forge holds it against no position in the diff, so it reads as outdated rather than as a comment on the change", finding.ID, posted.Path, posted.Line)
			}
			if posted.Body != finding.Description {
				t.Errorf("what reached the pull request reads %q, and socrates-review recorded %q", posted.Body, finding.Description)
			}
			t.Logf("%s landed at %s", finding.ID, posted.HTMLURL)
		}
	})

	t.Run("a finding outside the change is refused a line of its own", func(t *testing.T) {
		refused, code := commentOnALineAllowingRefusal(t, pull, outside[0])
		if code == 0 {
			deleteReviewComment(t, refused.ID)
			t.Fatalf("the forge accepted a comment on %s:%d, which the pull request's diff does not hold, so nothing about the placement rule is load bearing", outside[0].File, outside[0].Line)
		}
	})

	t.Run("everything else lands as one general comment", func(t *testing.T) {
		posted := oneGeneralComment(t, pull, outside)
		t.Cleanup(func() { deleteIssueComment(t, posted.ID) })

		for _, finding := range outside {
			if !strings.Contains(posted.Body, finding.Description) {
				t.Errorf("the general comment does not carry %s, which was forwarded: %q", finding.ID, posted.Body)
			}
		}
		t.Logf("%d findings outside the change landed as one comment at %s", len(outside), posted.HTMLURL)
	})
}

func theFixture(t *testing.T) pullRequest {
	t.Helper()

	var read pullRequest
	if err := json.Unmarshal([]byte(gh(t, "api", pullPath(t))), &read); err != nil {
		t.Fatalf("reading what the forge says about pull request %d: %v", theFixturePullRequest, err)
	}
	return read
}

func theirChangeCheckedOut(t *testing.T, pull pullRequest) string {
	t.Helper()

	clone := filepath.Join(t.TempDir(), "repository-under-review")
	gh(t, "repo", "clone", theRepositoryUnderReview, clone, "--", "--quiet")
	git(t, clone, "switch", "--quiet", "--detach", pull.Head.SHA)

	if head := git(t, clone, "rev-parse", "HEAD"); head != pull.Head.SHA {
		t.Fatalf("the clone is at %s, and the pull request's head is %s", head, pull.Head.SHA)
	}
	return clone
}

func socratesReview(t *testing.T) string {
	t.Helper()

	binary := filepath.Join(t.TempDir(), "socrates-review")
	source.Build(t, "./review/cmd/socrates-review", binary)
	return binary
}

func reviewed(t *testing.T, binary, clone string, pull pullRequest) string {
	t.Helper()

	printed := socrates(t, binary, clone, whileTheRoundIsParked, "begin", "--wait",
		"--base", "origin/"+pull.Base.Ref,
		"--not-my-change",
		"--intent", theAuthorsAccount(t))

	name, _, _ := strings.Cut(strings.TrimSpace(printed), "\n")
	if name == "" {
		t.Fatal("socrates-review begin printed no run identifier")
	}
	return name
}

func theAuthorsAccount(t *testing.T) string {
	t.Helper()
	return gh(t, "pr", "view", strconv.Itoa(theFixturePullRequest), "--repo", theRepositoryUnderReview, "--json", "title,body", "--template", "{{.title}}\n\n{{.body}}")
}

func whatSocratesReported(t *testing.T, binary, clone, name string, want int) reportedRound {
	t.Helper()

	var read reportedRound
	printed := socrates(t, binary, clone, want, "findings", name, "--json")
	if err := json.Unmarshal([]byte(printed), &read); err != nil {
		t.Fatalf("reading what socrates-review findings --json reported: %v\n%s", err, printed)
	}
	if len(read.Findings) == 0 {
		t.Fatal("socrates-review reported no findings, so there is nothing to forward")
	}
	return read
}

func commentOnALine(t *testing.T, pull pullRequest, finding reportedFinding) postedComment {
	t.Helper()

	posted, code := commentOnALineAllowingRefusal(t, pull, finding)
	if code != 0 {
		t.Fatalf("the forge refused a comment on %s:%d, which socrates-review reports as inside the change", finding.File, finding.Line)
	}
	return posted
}

func commentOnALineAllowingRefusal(t *testing.T, pull pullRequest, finding reportedFinding) (postedComment, int) {
	t.Helper()

	printed, code := ghAllowingFailure(t, "api", pullPath(t)+"/comments",
		"--method", "POST",
		"--field", "body="+finding.Description,
		"--field", "commit_id="+pull.Head.SHA,
		"--field", "path="+finding.File,
		"--field", "line="+strconv.Itoa(finding.Line),
		"--field", "side=RIGHT")
	if code != 0 {
		return postedComment{}, code
	}

	var read postedComment
	if err := json.Unmarshal([]byte(printed), &read); err != nil {
		t.Fatalf("reading what the forge said about the comment it accepted: %v\n%s", err, printed)
	}
	return read, 0
}

func oneGeneralComment(t *testing.T, pull pullRequest, findings []reportedFinding) postedComment {
	t.Helper()

	var body strings.Builder
	body.WriteString("An independent reviewer read this change. These are the findings whose location this pull request does not show.\n")
	for _, finding := range findings {
		fmt.Fprintf(&body, "\n- %s:%d: %s\n", finding.File, finding.Line, finding.Description)
	}

	printed := gh(t, "api",
		fmt.Sprintf("repos/%s/issues/%d/comments", theRepositoryUnderReview, theFixturePullRequest),
		"--method", "POST",
		"--field", "body="+body.String())

	var read postedComment
	if err := json.Unmarshal([]byte(printed), &read); err != nil {
		t.Fatalf("reading what the forge said about the comment it accepted: %v\n%s", err, printed)
	}
	return read
}

func deleteReviewComment(t *testing.T, id int) {
	t.Helper()
	gh(t, "api", fmt.Sprintf("repos/%s/pulls/comments/%d", theRepositoryUnderReview, id), "--method", "DELETE")
}

func deleteIssueComment(t *testing.T, id int) {
	t.Helper()
	gh(t, "api", fmt.Sprintf("repos/%s/issues/comments/%d", theRepositoryUnderReview, id), "--method", "DELETE")
}

func pullPath(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("repos/%s/pulls/%d", theRepositoryUnderReview, theFixturePullRequest)
}

func gh(t *testing.T, arguments ...string) string {
	t.Helper()

	printed, code := ghAllowingFailure(t, arguments...)
	if code != 0 {
		t.Fatalf("gh %s exited %d", strings.Join(arguments, " "), code)
	}
	return printed
}

func ghAllowingFailure(t *testing.T, arguments ...string) (string, int) {
	t.Helper()
	return run(t, "gh", "", arguments...)
}

func git(t *testing.T, directory string, arguments ...string) string {
	t.Helper()

	printed, code := run(t, "git", directory, arguments...)
	if code != 0 {
		t.Fatalf("git %s exited %d", strings.Join(arguments, " "), code)
	}
	return printed
}

func socrates(t *testing.T, binary, clone string, want int, arguments ...string) string {
	t.Helper()

	printed, code := run(t, binary, clone, arguments...)
	if code != want {
		t.Fatalf("socrates-review %s exited %d, want %d: %s", strings.Join(arguments, " "), code, want, printed)
	}
	return printed
}

func run(t *testing.T, command, directory string, arguments ...string) (string, int) {
	t.Helper()

	running := exec.Command(command, arguments...)
	running.Dir = directory

	var out, failures bytes.Buffer
	running.Stdout = &out
	running.Stderr = &failures

	if err := running.Run(); err != nil {
		var exited *exec.ExitError
		if !errors.As(err, &exited) {
			t.Fatalf("running %s %s: %v: %s", command, strings.Join(arguments, " "), err, failures.String())
		}
		t.Logf("%s %s exited %d: %s", command, strings.Join(arguments, " "), exited.ExitCode(), strings.TrimSpace(failures.String()))
		return strings.TrimSpace(out.String()), exited.ExitCode()
	}
	return strings.TrimSpace(out.String()), 0
}
