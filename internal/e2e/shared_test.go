package e2e_test

import (
	"strings"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/process"
	"github.com/rui-ktei/socrates/internal/runstore"
)

const cleanScenario = `
default:
  structured: '{"findings":[],"risk_level":"low","risk_rationale":"a small and self-contained change"}'
`

const escalatingScenario = `
default:
  structured: '{"findings":[{"id":"f1","severity":"error","file":"internal/added.go","line":1,"description":"The added file answers wrongly for an empty input.","action":"no-op"},{"id":"f2","severity":"warning","file":"internal/added.go","line":1,"description":"The retry loop never backs off.","action":"auto-fix"},{"id":"f3","severity":"info","file":"internal/added.go","line":1,"description":"This package now owns two spellings of one idea."}],"risk_level":"medium","risk_rationale":"one changed path answers wrongly without erroring"}'
`

const informationalScenario = `
default:
  structured: '{"findings":[{"id":"f1","severity":"info","file":"internal/added.go","line":1,"description":"The new name reads better than the old one.","action":"no-op"}],"risk_level":"low","risk_rationale":"nothing in this change can answer wrongly"}'
`

const slowScenario = "default:\n  events: 40\n  interval: 60ms\n"

func reason(run *runstore.Run) string {
	if run.Verdict == nil {
		return "no verdict"
	}
	return run.Verdict.Reason
}

func firstLine(printed string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(printed), "\n")
	return line
}

func waitForDeath(t *testing.T, name string, process1, process2 int) {
	t.Helper()

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if !process.Alive(process1) && !process.Alive(process2) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Errorf("a process of run %s outlived the cancellation: reviewer %d alive=%v, agent %d alive=%v",
		name, process1, process.Alive(process1), process2, process.Alive(process2))
}
