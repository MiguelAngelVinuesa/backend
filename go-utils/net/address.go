package net

import (
	"strconv"
	"strings"
)

// SplitHostPort splits the given address into host and port.
// If no port is given in the address, the port will be zero.
func SplitHostPort(address string) (string, int) {
	list := strings.Split(address, ":")
	if len(list) == 1 {
		return list[0], 0
	}
	port, _ := strconv.Atoi(list[1])
	return list[0], port
}

// BuildAddress builds a full address from the given host and port.
func BuildAddress(host string, port int) string {
	if host == "" {
		return ""
	}
	if port == 0 {
		return host
	}
	return host + ":" + strconv.Itoa(port)
}
