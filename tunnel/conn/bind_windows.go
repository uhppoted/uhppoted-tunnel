package conn

import (
	"syscall"
)

func BindToDevice(conn syscall.RawConn, device string, IPv4 bool, c Conn) error {
	c.Warnf("bind to interface not supported for Microsoft Windows")

	return nil
}
