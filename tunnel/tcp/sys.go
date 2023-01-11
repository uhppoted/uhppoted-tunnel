package tcp

import (
	"bytes"
	"net"
)

// Ref. https://stackoverflow.com/questions/22751035/golang-distinguish-ipv4-ipv6
func isIPv4(ip net.IP) bool {
	if len(ip) == net.IPv4len {
		return true
	}

	prefix := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}
	if len(ip) == net.IPv6len && bytes.Equal(ip[0:12], prefix) {
		return true
	}

	return false
}
