package tcp

import (
	"fmt"
	"net"
	"syscall"
)

const IP_BOUND_IF = 25
const IPV6_BOUND_IF = 125

// Ref. https://djangocas.dev/blog/linux/linux-SO_BINDTODEVICE-and-mac-IP_BOUND_IF-to-bind-socket-to-a-network-interface
func bindToDevice(conn syscall.RawConn, device string, IPv4 bool) error {
	if device != "" {

		if ifaces, err := net.Interfaces(); err != nil {
			return err
		} else {
			for _, iface := range ifaces {
				if iface.Name == device {
					var operr error
					bind := func(fd uintptr) {
						if IPv4 {
							operr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, IP_BOUND_IF, iface.Index)
						} else {
							operr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, IPV6_BOUND_IF, iface.Index)
						}
					}

					if err := conn.Control(bind); err != nil {
						return err
					} else {
						return operr
					}
				}
			}

			return fmt.Errorf("network interface '%v' not found", device)
		}
	}

	return nil
}
