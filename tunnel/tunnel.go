package tunnel

import (
	"encoding/hex"
	"fmt"
	"os"
	"regexp"

	"github.com/uhppoted/uhppoted-tunnel/log"
)

type In interface {
	Listen(relay func([]byte) []byte) error
	Close()
}

type Out interface {
	Listen() error
	Send([]byte) []byte
	Close()
}

type Tunnel struct {
	in  In
	out Out
}

func NewTunnel(in In, out Out) *Tunnel {
	return &Tunnel{
		in:  in,
		out: out,
	}
}

func (t *Tunnel) Run(interrupt chan os.Signal) {
	infof("%v", "uhppoted-tunnel::run")

	go func() {
		if err := t.out.Listen(); err != nil {
			fatalf("%v", err)
		}
	}()

	go func() {
		if err := t.in.Listen(t.redirect); err != nil {
			fatalf("%v", err)
		}
	}()

	<-interrupt

	t.in.Close()
	t.out.Close()
}

func (t *Tunnel) redirect(message []byte) []byte {
	return t.out.Send(message)
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
