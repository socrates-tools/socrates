package surface_test

import (
	"context"
	"net"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/socrates-devtool/socrates/internal/surface"
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

const theIdlePeriod = time.Hour

func TestAnIdleSessionIsCollectedAndTheOnesBeingReadAreNot(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		open := surface.NewSessions(time.Now())
		open.Open("beingread", "one", twoDecisions(), time.Now())
		open.Open("abandoned", "two", oneOtherDecision(), time.Now())

		read, standing := open.Exchange("beingread")
		if !standing {
			t.Fatal("a session just opened is not held")
		}
		read.Reading(1, time.Now())

		time.Sleep(theIdlePeriod + time.Minute)
		collected, empty := open.Collect(time.Now(), theIdlePeriod)

		if !slices.Equal(collected, []string{"abandoned"}) {
			t.Fatalf("the surface collected %v, and one session had a reader throughout while the other had none", collected)
		}
		if empty {
			t.Error("the surface reports itself empty while a session is being read")
		}
		if _, standing := open.Exchange("abandoned"); standing {
			t.Error("a session collected for having had no reader is still held")
		}

		still, standing := open.Exchange("beingread")
		if !standing {
			t.Fatal("the session being read was collected with the one nobody opened")
		}
		if _, err := still.Receive(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}}); err != nil {
			t.Errorf("the session being read no longer takes an answer: %v", err)
		}
	})
}

func TestAServerThatWasNeverGivenASessionStopsToo(t *testing.T) {
	// scenario: An idle server stops
	synctest.Test(t, func(t *testing.T) {
		open := surface.NewSessions(time.Now())

		if _, empty := open.Collect(time.Now(), theIdlePeriod); empty {
			t.Error("a surface reports itself idle the moment it comes up")
		}

		time.Sleep(theIdlePeriod + time.Minute)
		if _, empty := open.Collect(time.Now(), theIdlePeriod); !empty {
			t.Error("a surface that was never given a session is still standing after the idle period")
		}
	})
}

func TestAServerHoldingNoSessionStopsOnceItHasHeldNoneForTheIdlePeriod(t *testing.T) {
	// scenario: An idle server stops
	synctest.Test(t, func(t *testing.T) {
		open := surface.NewSessions(time.Now())
		open.Open("theonly", "one", twoDecisions(), time.Now())

		time.Sleep(theIdlePeriod + time.Minute)
		collected, empty := open.Collect(time.Now(), theIdlePeriod)
		if !slices.Equal(collected, []string{"theonly"}) {
			t.Fatalf("the surface collected %v, and the one session it held had no reader", collected)
		}
		if empty {
			t.Error("the surface stopped in the moment its last session was collected, rather than after holding none for the idle period")
		}

		time.Sleep(theIdlePeriod + time.Minute)
		if _, empty := open.Collect(time.Now(), theIdlePeriod); !empty {
			t.Error("the surface held no session for the idle period and is still standing")
		}
	})
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

func TestASessionACallerIsWaitingOnIsNotCollected(t *testing.T) {
	// scenario: An idle session is collected and the ones being read are not
	const idle = 200 * time.Millisecond

	held := serving(t, twoDecisions(), idle)

	waiting := make(chan surface.Answered, 1)
	go func() {
		defer close(waiting)
		read, arrived := held.await(t, 20*time.Second)
		if arrived {
			waiting <- read
		}
	}()

	time.Sleep(6 * idle)

	if holds := held.server.Holds(); len(holds) != 1 {
		t.Fatalf("the surface holds %v after %s with a caller waiting on its only session", holds, 6*idle)
	}

	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}})

	read, arrived := <-waiting
	if !arrived {
		t.Fatal("the caller that waited past the idle period was never answered")
	}
	if read.Decisions[0].ID != "c-1" {
		t.Errorf("what reached the caller that waited past the idle period is %+v", read)
	}
}
