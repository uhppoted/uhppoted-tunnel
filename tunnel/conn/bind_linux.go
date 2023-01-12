package conn

import (
	"syscall"
)

func BindToDevice(connection syscall.RawConn, device string, IPv4 bool, c Conn) error {
	if device != "" {
		var operr error
		bind := func(fd uintptr) {
			c.Infof("binding to interface %v", device)
			operr = syscall.BindToDevice(int(fd), device)
		}

		if err := connection.Control(bind); err != nil {
			return err
		} else {
			return operr
		}
	}

	return nil
}
