package tunnel

import (
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/uhppoted/uhppoted-tunnel/log"
	"github.com/uhppoted/uhppoted-tunnel/router"
)

type UDP interface {
	Close()
	Run(*router.Switch) error
	Send(uint32, []byte)
}

type TCP interface {
	Close()
	Run(*router.Switch) error
	Send(uint32, []byte)
}

type Message struct {
	id      uint32
	message []byte
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

func (t *Tunnel) Run(interrupt chan os.Signal) error {
	infof("", "%v", "uhppoted-tunnel::run")

	q := router.NewSwitch(func(id uint32, message []byte) {
		t.udp.Send(id, message)
	})

	p := router.NewSwitch(func(id uint32, message []byte) {
		t.tcp.Send(id, message)
	})

	u := make(chan error)
	v := make(chan error)

	go func() {
		if err := t.udp.Run(&p); err != nil {
			u <- err // fatalf("%v", err)
		}
	}()

	go func() {
		if err := t.tcp.Run(&q); err != nil {
			v <- err // fatalf("%v", err)
		}
	}()

	select {
	case <-interrupt:

	case err := <-u:
		return err

	case err := <-v:
		return err
	}

	infof("", "closing")

	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		router.Close()
	}()

	go func() {
		defer wg.Done()
		t.udp.Close()
	}()

	go func() {
		defer wg.Done()
		t.tcp.Close()
	}()

	wg.Wait()
	infof("", "closed")

	return nil
}

func dump(m []byte, prefix string) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), prefix+"$1"))
}

func dumpf(tag string, message []byte, format string, args ...any) {
	hex := dump(message, "                                  ")
	preamble := fmt.Sprintf(format, args...)

	debugf(tag, "%v\n%s", preamble, hex)
}

func debugf(tag, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Debugf(f, args...)
}

func infof(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Infof(f, args...)
}

func warnf(tag, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Warnf(f, args...)
}

func errorf(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Errorf(f, args...)
}

func fatalf(format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", "", format)

	log.Fatalf(f, args...)
}
