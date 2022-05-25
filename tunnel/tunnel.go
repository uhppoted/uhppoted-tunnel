package tunnel

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"regexp"

	"github.com/uhppoted/uhppoted-tunnel/log"
)

type Tunnel struct {
}

func (t Tunnel) Run(interrupt chan os.Signal) {
	infof("%v", "uhppoted-tunnel::run")

	if err := t.listen("0.0.0.0:60000", interrupt); err != nil {
		log.Errorf("%v", err)
	}
}

func (t Tunnel) listen(bind string, interrupt chan os.Signal) error {
	addr, err := net.ResolveUDPAddr("udp", bind)
	if err != nil {
		return err
	} else if addr == nil {
		return fmt.Errorf("unable to resolve UDP listen address '%v'", bind)
	} else if addr.Port == 0 {
		return fmt.Errorf("tunnel::listen requires a non-zero UDP port")
	}

	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("Error opening UDP listen socket (%v)", err)
	} else if sock == nil {
		return fmt.Errorf("Failed to open UDP socket (%v)", sock)
	}

	go func() {
		debugf("listening on %v", bind)
		buffer := make([]byte, 2048)

		for {

			N, remote, err := sock.ReadFromUDP(buffer)
			if err != nil {
				debugf("%v", err)
				break
			}

			hex := dump(buffer[:N], "                           ")

			debugf("received %v bytes from %v\n%s\n", N, remote, hex)
		}
	}()

	<-interrupt

	sock.Close()

	return nil
}

func dump(m []byte, prefix string) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), prefix+"$1"))
}

func debugf(format string, args ...any) {
	log.Debugf(format, args...)
}

func infof(format string, args ...any) {
	log.Infof(format, args...)
}

func warnf(format string, args ...any) {
	log.Warnf(format, args...)
}

func errorf(format string, args ...any) {
	log.Errorf(format, args...)
}

func fatalf(format string, args ...any) {
	log.Fatalf(format, args...)
}
