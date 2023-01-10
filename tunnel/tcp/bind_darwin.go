package tcp

import (
	"syscall"
)

// FIXME: https://djangocas.dev/blog/linux/linux-SO_BINDTODEVICE-and-mac-IP_BOUND_IF-to-bind-socket-to-a-network-interface/
func bindToDevice(conn syscall.RawConn, device string) error {
	return nil
}
