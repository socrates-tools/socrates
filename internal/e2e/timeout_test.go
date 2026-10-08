package e2e_test

import (
	"strings"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/internal/testing/source"
)

const goTestsDefaultTimeout = 10 * time.Minute

func TestTheSuiteStatesATimeoutLongerThanGoTestsDefault(t *testing.T) {
	t.Parallel()
	planned, code := runMake(t, source.Root(t), "-n", "test")
	if code != 0 {
		t.Fatalf("make -n test exited %d: %s", code, planned)
	}

	invocation := ""
	for line := range strings.SplitSeq(planned, "\n") {
		if strings.Contains(line, "go test") {
			invocation = line
		}
	}
	if invocation == "" {
		t.Fatalf("make test runs no go test, so there is no timeout to read:\n%s", planned)
	}

	stated, found := timeoutStatedBy(strings.Fields(invocation))
	if !found {
		t.Fatalf("make test runs %q, which states no -timeout, so the end-to-end package fails as a panic in whichever test is under way at Go's default of %s", invocation, goTestsDefaultTimeout)
	}
	if stated <= goTestsDefaultTimeout {
		t.Errorf("make test states a timeout of %s, which is no longer than Go's default of %s that the end-to-end package outgrew", stated, goTestsDefaultTimeout)
	}
}

func timeoutStatedBy(arguments []string) (time.Duration, bool) {
	for position, argument := range arguments {
		value, joined := strings.CutPrefix(argument, "-timeout=")
		if !joined {
			if argument != "-timeout" || position+1 == len(arguments) {
				continue
			}
			value = arguments[position+1]
		}
		stated, err := time.ParseDuration(value)
		if err != nil {
			return 0, false
		}
		return stated, true
	}
	return 0, false
}
