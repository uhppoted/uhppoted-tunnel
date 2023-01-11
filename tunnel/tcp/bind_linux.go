package tcp

import (
	"syscall"
)

func bindToDevice(conn syscall.RawConn, device string, IPv4 bool) error {
	if device != "" {
		var operr error
		bind := func(fd uintptr) {
			operr = syscall.BindToDevice(int(fd), device)
		}

		if err := conn.Control(bind); err != nil {
			return err
		} else {
			return operr
		}
	}

	return nil
}
