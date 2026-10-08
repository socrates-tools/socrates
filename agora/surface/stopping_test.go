package surface_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/agora/surface"
)

func TestWaitDoesNotReturnWhileAResponseIsStillOwedToACaller(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	request, err := http.NewRequest(http.MethodGet, held.under("/answer?within=1s"), nil)
	if err != nil {
		t.Fatalf("addressing the surface: %v", err)
	}
	request.Header.Set(surface.TokenHeader, held.token)

	answered := make(chan error, 1)
	go func() {
		response, err := http.DefaultClient.Do(request)
		if err == nil {
			_ = response.Body.Close()
		}
		answered <- err
	}()

	time.Sleep(200 * time.Millisecond)
	stopped := time.Now()
	go held.server.Stop()
	held.server.Wait()

	if waited := time.Since(stopped); waited < 500*time.Millisecond {
		t.Errorf("Wait returned %s after Stop while a caller was still owed its response, and a process that exits then drops it", waited)
	}
	if err := <-answered; err != nil {
		t.Errorf("the caller waiting on the surface when it stopped was not answered: %v", err)
	}
}
