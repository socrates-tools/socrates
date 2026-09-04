package surface_test

import (
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/surface"
)

func portOf(t *testing.T, address string) int {
	t.Helper()

	_, named, err := net.SplitHostPort(address)
	if err != nil {
		t.Fatalf("reading the port of %q: %v", address, err)
	}
	port, err := strconv.Atoi(named)
	if err != nil {
		t.Fatalf("the port of %q is not a number: %v", address, err)
	}
	return port
}

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

func TestTheWindowIsWalkedFromTheTopDown(t *testing.T) {
	addresses := surface.WindowAddresses()

	if len(addresses) != 10 {
		t.Fatalf("the window holds %d ports: %v", len(addresses), addresses)
	}
	if addresses[0] != surface.DefaultAddress {
		t.Errorf("the top of the window is %q, and the address bound where nothing states one is %q", addresses[0], surface.DefaultAddress)
	}

	top := portOf(t, addresses[0])
	for at, address := range addresses {
		if port := portOf(t, address); port != top-at {
			t.Errorf("the window holds %q at %d, and counting down from %d reaches %d there", address, at, top, top-at)
		}
		if host, _, _ := net.SplitHostPort(address); net.ParseIP(host) == nil || !net.ParseIP(host).IsLoopback() {
			t.Errorf("the window holds %q, and the surface binds a loopback address only", address)
		}
	}
}
