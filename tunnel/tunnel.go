package tunnel

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"

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

	// if remote != "" {
	// 	t.connect(remote, interrupt)
	// } else {

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

func (t *Tunnel) connect(addr string, interrupt chan os.Signal) error {
	socket, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("Error connecting to  %v (%v)", addr, err)
	} else if socket == nil {
		return fmt.Errorf("Failed to create TCP connection to %v (%v)", addr, socket)
	}

	defer socket.Close()

	infof("TCP  connected to %v", addr)

	go func() {
		buffer := make([]byte, 2048)
		for {
			if N, err := socket.Read(buffer); err != nil {
				warnf("%v", err)
				break
			} else {
				hex := dump(buffer[:N], "                           ")
				debugf("TCP  received %v bytes from %v\n%s\n", N, socket.RemoteAddr(), hex)

				ix := 0
				for ix < N {
					size := uint(buffer[ix])
					size <<= 8
					size += uint(buffer[ix+1])

					if reply, err := t.broadcast(buffer[ix+2 : ix+2+int(size)]); err != nil {
						warnf("%v", err)
					} else if reply == nil || len(reply) == 0 {
						warnf("empty reply (%v)", reply)
					} else {
						packet := make([]byte, len(reply)+2)

						packet[0] = byte((len(reply) >> 8) & 0x00ff)
						packet[1] = byte((len(reply) >> 0) & 0x00ff)
						copy(packet[2:], reply)

						if N, err := socket.Write(packet); err != nil {
							warnf("error redirecting reply to %v (%v)", socket.RemoteAddr(), err)
						} else if N != len(packet) {
							warnf("replied with %v of %v bytes to %v", N, len(reply), socket.RemoteAddr())
						} else {
							infof("replied with %v bytes to %v", len(reply), socket.RemoteAddr())
						}
					}

					ix += 2 + int(size)
				}
			}
		}
	}()

	<-interrupt

	return nil
}

func (t *Tunnel) redirect(message []byte) []byte {
	return t.out.Send(message)
}

func (t *Tunnel) broadcast(message []byte) ([]byte, error) {
	hex := dump(message, "                           ")
	debugf("broadcast%v\n%s\n", "", hex)

	addr, err := net.ResolveUDPAddr("udp", "192.168.1.255:60000")
	if err != nil {
		return nil, err
	}

	bind, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	if err != nil {
		return nil, err
	}

	socket, err := net.ListenUDP("udp", bind)
	if err != nil {
		return nil, err
	} else if socket == nil {
		return nil, fmt.Errorf("invalid UDP socket (%v)", socket)
	}

	defer socket.Close()

	if err := socket.SetWriteDeadline(time.Now().Add(1000 * time.Millisecond)); err != nil {
		return nil, err
	}

	if err := socket.SetReadDeadline(time.Now().Add(5000 * time.Millisecond)); err != nil {
		return nil, err
	}

	if N, err := socket.WriteToUDP(message, addr); err != nil {
		return nil, err
	} else {
		debugf(" ... sent %v bytes to %v\n", N, addr)
	}

	reply := make([]byte, 2048)

	if N, remote, err := socket.ReadFromUDP(reply); err != nil {
		return nil, err
	} else {
		debugf(" ... received %v bytes from %v\n%s", N, remote, dump(reply[:N], " ...          "))

		return reply[:N], nil
	}
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
