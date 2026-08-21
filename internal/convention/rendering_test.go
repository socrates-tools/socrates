package convention_test

import (
	"path/filepath"
	"testing"

	"github.com/rui-ktei/socrates/internal/agent"
	"github.com/rui-ktei/socrates/internal/convention"
	"github.com/rui-ktei/socrates/internal/reviewer"
	"github.com/rui-ktei/socrates/internal/runlog"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

func TestOnlyOnePlaceTurnsAnEventIntoALinePeopleRead(t *testing.T) {
	kinds := []string{
		runlog.EventRotated,
		agent.EventStarted,
		agent.EventCompleted,
		agent.EventFailed,
		agent.EventToolServersLoaded,
		reviewer.EventRunStarted,
		reviewer.EventRunEnded,
		reviewer.EventRoundStarted,
		reviewer.EventRoundParked,
		reviewer.EventResponse,
	}
	defining := []string{
		filepath.Join("internal", "runlog"),
		filepath.Join("internal", "agent"),
		filepath.Join("internal", "reviewer"),
	}

	naming, err := convention.FilesNamingEventKinds(source.Root(t), kinds, defining, convention.EventKindsNamed)
	if err != nil {
		t.Fatalf("reading which files name an event kind: %v", err)
	}

	for _, path := range naming.Unnamed {
		t.Errorf("%s names an event kind, and only the renderer may turn one into a line: name it in convention.EventKindsNamed with why it is not a second rendering", path)
	}
	for _, path := range naming.Stale {
		t.Errorf("convention.EventKindsNamed names %s, which names no event kind any more", path)
	}
	if len(naming.Named) == 0 {
		t.Error("no file names an event kind, so this check is reading nothing")
	}
}
