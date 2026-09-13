package loopback_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/loopback"
)

func TestTheWindowIsWalkedFromTheTopDown(t *testing.T) {
	addresses := loopback.Addresses()

	if len(addresses) != 10 {
		t.Fatalf("the window holds %d ports: %v", len(addresses), addresses)
	}
	if addresses[0] != loopback.DefaultAddress {
		t.Errorf("the top of the window is %q, and the address bound where nothing states one is %q", addresses[0], loopback.DefaultAddress)
	}

	top := portOf(t, addresses[0])
	for at, address := range addresses {
		if port := portOf(t, address); port != top-at {
			t.Errorf("the window holds %q at %d, and counting down from %d reaches %d there", address, at, top, top-at)
		}
	}
}

func portOf(t *testing.T, address string) int {
	t.Helper()

	_, named, found := strings.Cut(address, ":")
	if !found {
		t.Fatalf("%q names no port", address)
	}
	port, err := strconv.Atoi(named)
	if err != nil {
		t.Fatalf("the port of %q is not a number: %v", address, err)
	}
	return port
}
