package surface_test

import (
	"net"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/loopback"
)

func aWindowOf(t *testing.T, held, free int) []string {
	t.Helper()

	addresses := make([]string, 0, held+free)
	for range held {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("holding a port the window is walked past: %v", err)
		}
		t.Cleanup(func() { _ = listener.Close() })
		addresses = append(addresses, listener.Addr().String())
	}
	for range free {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("choosing a port the window can take: %v", err)
		}
		address := listener.Addr().String()
		if err := listener.Close(); err != nil {
			t.Fatalf("releasing the port the window can take: %v", err)
		}
		addresses = append(addresses, address)
	}
	return addresses
}

func TestTheHighestFreePortOfTheWindowIsTaken(t *testing.T) {
	addresses := aWindowOf(t, 2, 2)

	free, err := surface.FirstFreeIn(addresses)
	if err != nil {
		t.Fatalf("taking the highest free port of the window: %v", err)
	}
	if free != addresses[2] {
		t.Errorf("the first two ports of the window are held and it resolved to %q, and the first free one is %q", free, addresses[2])
	}
}

func TestAWholeWindowHeldIsARefusalThatNamesTheWindow(t *testing.T) {
	addresses := aWindowOf(t, 4, 0)

	free, err := surface.FirstFreeIn(addresses)
	if err == nil {
		t.Fatalf("every port of the window is held and it resolved to %q", free)
	}
	for _, named := range []string{addresses[0], addresses[len(addresses)-1]} {
		if !strings.Contains(err.Error(), named) {
			t.Errorf("refusing a whole window held does not name %s:\n%s", named, err)
		}
	}
}

func TestEveryAddressOfTheWindowIsOneTheSurfaceWillBind(t *testing.T) {
	for _, address := range loopback.Addresses() {
		if err := surface.CheckLoopback(address); err != nil {
			t.Errorf("the window states %q, and the surface binds a loopback address only: %v", address, err)
		}
	}
}
