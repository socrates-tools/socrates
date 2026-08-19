package e2e_test

import "github.com/rui-ktei/socrates/internal/runstore"

func reason(run *runstore.Run) string {
	if run.Verdict == nil {
		return "no verdict"
	}
	return run.Verdict.Reason
}
