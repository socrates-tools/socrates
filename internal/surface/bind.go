package surface

import (
	"net"

	"github.com/rui-ktei/socrates/internal/refusal"
)

const DefaultAddress = "127.0.0.1:0"

const localhostName = "localhost"

func CheckLoopback(address string) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return refusal.Wrapping(err, addressIsNotHostAndPort, address, err)
	}
	if port == "" {
		return refusal.New(addressNamesNoPort, address)
	}
	if host == "" {
		return refusal.New(addressNamesEveryInterface, address)
	}
	if host == localhostName {
		return nil
	}

	stated := net.ParseIP(host)
	if stated == nil {
		return refusal.New(addressIsNotAnAddress, address, host)
	}
	if stated.IsUnspecified() {
		return refusal.New(addressNamesEveryInterface, address)
	}
	if !stated.IsLoopback() {
		return refusal.New(addressIsNotLoopback, address)
	}
	return nil
}

func hostIsTheSurface(claimed, bound string) error {
	if claimed == "" {
		return refusal.New(requestNamesNoHost)
	}

	claimedHost, claimedPort, err := net.SplitHostPort(claimed)
	if err != nil {
		return refusal.New(requestNamesTheHost, claimed)
	}
	_, boundPort, err := net.SplitHostPort(bound)
	if err != nil {
		return refusal.New(requestNamesTheHost, claimed)
	}
	if claimedPort != boundPort {
		return refusal.New(requestNamesTheHost, claimed)
	}
	if claimedHost == localhostName {
		return nil
	}
	if stated := net.ParseIP(claimedHost); stated != nil && stated.IsLoopback() {
		return nil
	}
	return refusal.New(requestNamesTheHost, claimed)
}
