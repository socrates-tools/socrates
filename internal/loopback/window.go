package loopback

import "strconv"

const (
	host     = "127.0.0.1"
	topPort  = 18018
	ports    = 10
	basePort = topPort - ports + 1
)

var DefaultAddress = at(topPort)

func Addresses() []string {
	addresses := make([]string, 0, ports)
	for port := topPort; port >= basePort; port-- {
		addresses = append(addresses, at(port))
	}
	return addresses
}

func at(port int) string {
	return host + ":" + strconv.Itoa(port)
}
