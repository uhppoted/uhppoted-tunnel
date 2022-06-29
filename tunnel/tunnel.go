package tunnel

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/uhppoted/uhppoted-tunnel/log"
	"github.com/uhppoted/uhppoted-tunnel/router"
)

type Conn interface {
	Close()
	Run(*router.Switch) error
	Send(uint32, []byte)
}

type Message struct {
	id      uint32
	message []byte
}

type Tunnel struct {
	in  Conn
	out Conn
	ctx context.Context
}

func NewTunnel(in Conn, out Conn, ctx context.Context) *Tunnel {
	return &Tunnel{
		in:  in,
		out: out,
		ctx: ctx,
	}
}

func (t *Tunnel) Run(interrupt chan os.Signal) error {
	infof("", "%v", "uhppoted-tunnel::run")

	p := router.NewSwitch(func(id uint32, message []byte) {
		t.out.Send(id, message)
	})

	q := router.NewSwitch(func(id uint32, message []byte) {
		t.in.Send(id, message)
	})

	u := make(chan error)
	v := make(chan error)

	go func() {
		if err := t.in.Run(&p); err != nil {
			u <- err
		}
	}()

	go func() {
		if err := t.out.Run(&q); err != nil {
			v <- err
		}
	}()

	select {
	case <-t.ctx.Done():
		// closing

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
		t.in.Close()
	}()

	go func() {
		defer wg.Done()
		t.out.Close()
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
