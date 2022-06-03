package tunnel

import (
	"encoding/hex"
	"fmt"
	"os"
	"regexp"

	"github.com/uhppoted/uhppoted-tunnel/log"
)

type relay func(uint32, []byte)

type UDP interface {
	Close()
	Run(relay) error
	Send(uint32, []byte)
}

type TCP interface {
	Close()
	Run(relay) error
	Send(uint32, []byte)
}

type Tunnel struct {
	udp UDP
	tcp TCP
}

func NewTunnel(udp UDP, tcp TCP) *Tunnel {
	return &Tunnel{
		udp: udp,
		tcp: tcp,
	}
}

func (t *Tunnel) Run(interrupt chan os.Signal) {
	infof("%v", "uhppoted-tunnel::run")

	p := func(id uint32, message []byte) {
		t.tcp.Send(id, message)
	}

	q := func(id uint32, message []byte) {
		t.udp.Send(id, message)
	}

	go func() {
		if err := t.udp.Run(p); err != nil {
			fatalf("%v", err)
		}
	}()

	go func() {
		if err := t.tcp.Run(q); err != nil {
			fatalf("%v", err)
		}
	}()

	<-interrupt

	t.udp.Close()
	t.tcp.Close()
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
