package surface

import (
	"net"
	"net/url"
	"strconv"

	"github.com/ktei/socrates/internal/refusal"
)

const (
	windowTopPort  = 18018
	windowPorts    = 10
	windowHost     = "127.0.0.1"
	windowBasePort = windowTopPort - windowPorts + 1
)

var DefaultAddress = inTheWindow(windowTopPort)

const localhostName = "localhost"

func inTheWindow(port int) string {
	return net.JoinHostPort(windowHost, strconv.Itoa(port))
}

func WindowAddresses() []string {
	addresses := make([]string, 0, windowPorts)
	for port := windowTopPort; port >= windowBasePort; port-- {
		addresses = append(addresses, inTheWindow(port))
	}
	return addresses
}

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

func CheckFree(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return refusal.New(portIsHeldByAnother, address)
	}
	return listener.Close()
}

func FirstFreeIn(addresses []string) (string, error) {
	for _, address := range addresses {
		if CheckFree(address) == nil {
			return address, nil
		}
	}
	return "", refusal.New(everyPortInTheWindowIsHeld, addresses[len(addresses)-1], addresses[0])
}

func CheckOpenable(address string) error {
	read, err := url.Parse(address)
	if err != nil {
		return refusal.Wrapping(err, addressIsNotAURL, address, err)
	}
	if read.Scheme != "http" && read.Scheme != "https" {
		return refusal.New(addressIsNotOpenable, address)
	}
	if read.Host == "" {
		return refusal.New(addressNamesNoHost, address)
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
