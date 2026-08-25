package surface_test

import (
	"context"
	"net"
	"testing"
	"testing/synctest"
	"time"

	"github.com/rui-ktei/socrates/internal/surface"
)

func TestNobodyAnswers(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		exchange := surface.NewExchange(twoDecisions(), time.Now())
		started := time.Now()

		read, arrived := exchange.Await(context.Background(), surface.PollDeadline)

		if arrived {
			t.Fatalf("nobody answered and something was delivered: %+v", read)
		}
		if waited := time.Since(started); waited != surface.PollDeadline {
			t.Errorf("the wait ended after %s, and it is bounded at %s", waited, surface.PollDeadline)
		}
	})
}

func TestAWaitThatIsAnsweredReturnsBeforeItsBound(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		exchange := surface.NewExchange(twoDecisions(), time.Now())
		started := time.Now()

		go func() {
			time.Sleep(time.Hour)
			if _, err := exchange.Receive(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}}); err != nil {
				t.Errorf("sending an answer: %v", err)
			}
		}()

		read, arrived := exchange.Await(context.Background(), surface.PollDeadline)

		if !arrived {
			t.Fatal("an answer sent within the bound was not delivered")
		}
		if read.Decisions[0].ID != "c-1" {
			t.Errorf("what the wait returned is %+v", read)
		}
		if waited := time.Since(started); waited != time.Hour {
			t.Errorf("a wait answered after an hour returned after %s", waited)
		}
	})
}

func TestTheBoundIsAPersonsRatherThanAProcesss(t *testing.T) {
	// scenario: Nobody answers
	if surface.PollDeadline < time.Hour {
		t.Errorf("the wait is bounded at %s, and it waits on a person who may be at lunch", surface.PollDeadline)
	}
}

func TestAnIdleServerStops(t *testing.T) {
	held := serving(t, twoDecisions(), 200*time.Millisecond)

	select {
	case <-held.server.Stopping():
	case <-time.After(10 * time.Second):
		t.Fatal("no session was open for the idle period and the server is still running")
	}
	held.server.Wait()

	listener, err := net.Listen("tcp", held.server.Address())
	if err != nil {
		t.Fatalf("the port the stopped server held is not free: %v", err)
	}
	listener.Close()
}
