package tunnel

import (
	"fmt"
	"net"
	"time"
)

type tcpInClient struct {
	addr          *net.TCPAddr
	maxRetries    int
	maxRetryDelay time.Duration
}

const RETRY_MIN_DELAY = 5 * time.Second

func NewTCPIn(spec string, maxRetries int, maxRetryDelay time.Duration) (*tcpInClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	}

	in := tcpInClient{
		addr:          addr,
		maxRetries:    maxRetries,
		maxRetryDelay: maxRetryDelay,
	}

	return &in, nil
}

func (tcp *tcpInClient) Listen(relay func([]byte) []byte) error {
	retryDelay := RETRY_MIN_DELAY
	retries := 0

	for tcp.maxRetries < 0 || retries < tcp.maxRetries {
		infof("connecting to %v", tcp.addr)

		if socket, err := net.Dial("tcp", fmt.Sprintf("%v", tcp.addr)); err != nil {
			warnf("%v", err)
		} else if socket == nil {
			warnf("Failed to create TCP connection to %v (%v)", tcp.addr, socket)
		} else {
			retries = 0
			retryDelay = RETRY_MIN_DELAY

			if err := tcp.listen(socket, relay); err != nil {
				warnf("%v", err)
			}
		}

		infof("connection failed ... retrying in %v", retryDelay)

		time.Sleep(retryDelay)

		retries++
		retryDelay *= 2
		if retryDelay > tcp.maxRetryDelay {
			retryDelay = tcp.maxRetryDelay
		}
	}

	return fmt.Errorf("Connect to %v failed (retry count exceeded %v)", tcp.addr, tcp.maxRetries)
}

func (tcp *tcpInClient) Close() {

}

func (tcp *tcpInClient) listen(socket net.Conn, relay func([]byte) []byte) error {
	infof("TCP  connected to %v", socket.RemoteAddr())

	defer socket.Close()

	buffer := make([]byte, 2048)

	for {
		N, err := socket.Read(buffer)
		if err != nil {
			return err
		}

		hex := dump(buffer[:N], "                           ")
		debugf("TCP  received %v bytes from %v\n%s\n", N, socket.RemoteAddr(), hex)

		ix := 0
		for ix < N {
			size := uint(buffer[ix])
			size <<= 8
			size += uint(buffer[ix+1])

			message := depacketize(buffer[ix : ix+2+int(size)])

			if reply := relay(message); reply != nil && len(reply) > 0 {
				packet := packetize(reply)

				if N, err := socket.Write(packet); err != nil {
					warnf("error relaying reply to %v (%v)", socket.RemoteAddr(), err)
				} else if N != len(packet) {
					warnf("relayed reply with %v of %v bytes to %v", N, len(reply), socket.RemoteAddr())
				} else {
					infof("relayed reply with %v bytes to %v", len(reply), socket.RemoteAddr())
				}
			}

			ix += 2 + int(size)
		}
	}
}
