package e2e_test

import (
	"encoding/json"
	"fmt"
	"maps"
	"testing"
)

type sketch = map[string]any

func written(t *testing.T, document sketch) string {
	t.Helper()

	raw, err := json.Marshal(document)
	if err != nil {
		t.Fatalf("writing the document: %v", err)
	}
	return string(raw)
}

func leaveIt(id string) sketch {
	return sketch{"label": "leave it", "meaning": "acceptable here", "command": "agree " + id}
}

func aDecision(id string, stated sketch) sketch {
	decision := sketch{
		"id":       id,
		"question": "what should be done about " + id,
		"mode":     "single-choice",
		"options":  []sketch{leaveIt(id)},
	}
	maps.Copy(decision, stated)
	return decision
}

func anEntry(id, name, severity, audience, location string) sketch {
	entry := sketch{"id": id}
	for key, value := range map[string]string{"name": name, "severity": severity, "audience": audience, "location": location} {
		if value != "" {
			entry[key] = value
		}
	}
	return entry
}

func aDocument(title string, source []sketch, decisions ...sketch) sketch {
	document := sketch{"title": title}
	if len(source) > 0 {
		document["source"] = source
	}
	if len(decisions) > 0 {
		document["decisions"] = decisions
	}
	return document
}

func aCacheDecision(mode string) sketch {
	return sketch{
		"id":       "c-1",
		"severity": "error",
		"audience": "a person",
		"location": "internal/cache/read.go:41",
		"question": "should the cache be read before the lock is taken",
		"meaning":  "two callers arriving together read the same stale entry",
		"mode":     mode,
		"options": []sketch{
			{"label": "take the lock first", "meaning": "correctness first", "command": "patch c-1"},
			{"label": "leave it", "meaning": "the staleness is acceptable", "command": `agree c-1 --reason "acceptable"`},
		},
	}
}

func oneCacheDocument(mode string) sketch {
	return aDocument("one thing to settle", nil, aCacheDecision(mode))
}

func leastSevereFirstSource() []sketch {
	return []sketch{
		anEntry("c-3", "the log line says nothing", "info", "the author patches", "internal/cache/read.go:80"),
		anEntry("c-2", "the write counts twice", "warning", "the author patches", "internal/cache/write.go:12"),
		anEntry("c-1", "the read is unlocked", "error", "needs a person", "internal/cache/read.go:41"),
	}
}

func leastSevereFirst() sketch {
	return aDocument("three things to settle", leastSevereFirstSource(), aDecision("c-3", nil), aDecision("c-2", nil), aDecision("c-1", nil))
}

func citingThreeEntries() sketch {
	return aDocument("one thing to settle", []sketch{
		anEntry("c-1", "the read is unlocked", "error", "a person", "internal/cache/read.go:41"),
		anEntry("c-2", "the write counts twice", "warning", "the author patches", "internal/cache/write.go:12"),
		anEntry("c-3", "", "info", "the author patches", "internal/cache/read.go:80"),
	}, aCacheDecision("single-choice"))
}

type revisions struct {
	t       *testing.T
	held    *surfaced
	page    *reader
	reached int
}

func showing(t *testing.T, first sketch) *revisions {
	t.Helper()

	held := standing(t)
	return &revisions{t: t, held: held, page: opening(t, held.serve(written(t, first))), reached: 1}
}

func (r *revisions) show(next sketch) {
	r.t.Helper()

	r.held.serveAgain(written(r.t, next))
	r.reached++
	r.page.waitFor(fmt.Sprintf(`document.body.dataset.revision === '%d'`, r.reached))
}

const severityOfAnElement = `((held) => { const mark = held.querySelector('.mark.severity'); return mark ? mark.dataset.severity : ''; })`
