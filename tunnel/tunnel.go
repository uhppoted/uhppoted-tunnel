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

type Tunnel struct {
	udp []net.Conn
	tcp []net.Conn
}

func NewTunnel() *Tunnel {
	return &Tunnel{
		udp: []net.Conn{},
		tcp: []net.Conn{},
	}
}

func (t *Tunnel) Run(remote string, interrupt chan os.Signal) {
	infof("%v", "uhppoted-tunnel::run")

	if remote != "" {
		t.connect(remote, interrupt)
	} else {

		udp := make(chan []byte)

		go func() {
			socket, err := net.Listen("tcp", "0.0.0.0:8080")
			if err != nil {
				fatalf("%v", err)
			}

			infof("TCP server listening on %v", socket.Addr())

			for {
				if client, err := socket.Accept(); err != nil {
					errorf("%v", err)
				} else {
					t.accept(client)
				}
			}
		}()

		go func() {
			for {
				msg := <-udp
				t.redirect(msg)
			}
		}()

		if err := t.listen("0.0.0.0:60000", udp, interrupt); err != nil {
			log.Errorf("%v", err)
		}
	}
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

					if err := t.broadcast(buffer[ix+2 : ix+2+int(size)]); err != nil {
						warnf("%v", err)
					}

					ix += 2 + int(size)
				}
			}
		}
	}()

	<-interrupt

	return nil
}

func (t *Tunnel) listen(bind string, ch chan []byte, interrupt chan os.Signal) error {
	addr, err := net.ResolveUDPAddr("udp", bind)
	if err != nil {
		return err
	} else if addr == nil {
		return fmt.Errorf("unable to resolve UDP listen address '%v'", bind)
	} else if addr.Port == 0 {
		return fmt.Errorf("'listen' requires a non-zero UDP port")
	}

	socket, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("Error opening UDP listen socket (%v)", err)
	} else if socket == nil {
		return fmt.Errorf("Failed to open UDP socket (%v)", socket)
	}

	defer socket.Close()

	go func() {
		infof("UDP  listening on %v", bind)
		buffer := make([]byte, 2048)

		for {
			N, remote, err := socket.ReadFromUDP(buffer)
			if err != nil {
				debugf("%v", err)
				break
			}

			hex := dump(buffer[:N], "                           ")
			debugf("UDP  received %v bytes from %v\n%s\n", N, remote, hex)

			ch <- buffer[:N]
		}
	}()

	<-interrupt

	return nil
}

func (t *Tunnel) accept(c net.Conn) {
	infof("TCP  incoming connection (%v)", c.RemoteAddr())

	if socket, ok := c.(*net.TCPConn); !ok {
		errorf("%v", "invalid TCP socket")
	} else {
		t.tcp = append(t.tcp, c)

		go func() {
			buffer := make([]byte, 2048)
			for {
				if N, err := socket.Read(buffer); err != nil {
					warnf("%v", err)
					break
				} else {
					hex := dump(buffer[:N], "                           ")
					debugf("TCP  received %v bytes from %v\n%s\n", N, socket.RemoteAddr(), hex)

					t.broadcast(buffer[:N])
				}
			}

			socket.Close()
		}()
	}
}

func (t *Tunnel) redirect(message []byte) {
	debugf("redirect (%v connections)", len(t.tcp))

	packet := make([]byte, len(message)+2)

	packet[0] = byte((len(message) >> 8) & 0x00ff)
	packet[1] = byte((len(message) >> 0) & 0x00ff)
	copy(packet[2:], message)

	for _, c := range t.tcp {
		if N, err := c.Write(packet); err != nil {
			warnf("error redirecting message to %v (%v)", c.RemoteAddr(), err)
		} else if N != len(packet) {
			warnf("redirected %v of %v bytes to %v", N, len(message), c.RemoteAddr())
		} else {
			infof("redirected %v bytes to %v", len(message), c.RemoteAddr())
		}
	}
}

func (t *Tunnel) broadcast(message []byte) error {
	hex := dump(message, "                           ")
	debugf("broadcast%v\n%s\n", "", hex)

	addr, err := net.ResolveUDPAddr("udp", "192.168.1.255:60000")
	if err != nil {
		return err
	}

	bind, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	if err != nil {
		return err
	}

	socket, err := net.ListenUDP("udp", bind)
	if err != nil {
		return err
	} else if socket == nil {
		return fmt.Errorf("invalid UDP socket (%v)", socket)
	}

	defer socket.Close()

	if err := socket.SetWriteDeadline(time.Now().Add(1000 * time.Millisecond)); err != nil {
		return err
	}

	if err := socket.SetReadDeadline(time.Now().Add(5000 * time.Millisecond)); err != nil {
		return err
	}

	if N, err := socket.WriteToUDP(message, addr); err != nil {
		return err
	} else {
		debugf(" ... sent %v bytes to %v\n", N, addr)
	}

	reply := make([]byte, 2048)

	if N, remote, err := socket.ReadFromUDP(reply); err != nil {
		return err
	} else {
		debugf(" ... received %v bytes from %v\n%s", N, remote, dump(reply[:N], " ...          "))
	}

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
