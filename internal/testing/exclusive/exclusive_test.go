package exclusive_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/internal/testing/exclusive"
)

func TestTestsHoldingTheSameMachineResourceNeverRunTogether(t *testing.T) {
	var holding, overlapped atomic.Int32

	t.Run("holders", func(t *testing.T) {
		for range 4 {
			t.Run("one", func(t *testing.T) {
				t.Parallel()
				exclusive.Machine(t, "the-exclusive-package-itself")
				if holding.Add(1) > 1 {
					overlapped.Add(1)
				}
				time.Sleep(50 * time.Millisecond)
				holding.Add(-1)
			})
		}
	})

	if overlapped.Load() != 0 {
		t.Errorf("%d holders ran while another held the same machine resource", overlapped.Load())
	}
}
