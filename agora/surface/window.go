package surface

import "strconv"

const (
	windowHost   = "127.0.0.1"
	windowTop    = 18018
	windowSize   = 10
	windowBottom = windowTop - windowSize + 1
)

var DefaultAddress = inTheWindow(windowTop)

func Addresses() []string {
	addresses := make([]string, 0, windowSize)
	for port := windowTop; port >= windowBottom; port-- {
		addresses = append(addresses, inTheWindow(port))
	}
	return addresses
}

func inTheWindow(port int) string {
	return windowHost + ":" + strconv.Itoa(port)
}
