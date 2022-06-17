package tunnel

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/types"
)

type tcpClient struct {
	addr          *net.TCPAddr
	maxRetries    int
	maxRetryDelay time.Duration
	timeout       time.Duration
	ch            chan types.Message
	closing       chan struct{}
	closed        chan struct{}
}

const RETRY_MIN_DELAY = 5 * time.Second

func NewTCPClient(spec string, maxRetries int, maxRetryDelay time.Duration) (*tcpClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	}

	in := tcpClient{
		addr:          addr,
		maxRetries:    maxRetries,
		maxRetryDelay: maxRetryDelay,
		timeout:       5 * time.Second,
		ch:            make(chan types.Message, 16),
		closing:       make(chan struct{}),
		closed:        make(chan struct{}),
	}

	return &in, nil
}

func (tcp *tcpClient) Close() {
	infof("TCP", "closing")
	close(tcp.closing)

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-tcp.closed:
		infof("TCP", "closed")

	case <-timeout.C:
		infof("TCP", "close timeout")
	}
}

func (tcp *tcpClient) Run(router *router.Switch) error {
	tcp.connect(router)
	tcp.closed <- struct{}{}

	return nil
}

func (tcp *tcpClient) Send(id uint32, msg []byte) {
	select {
	case tcp.ch <- types.Message{ID: id, Message: msg}:
	default:
	}
}

func (tcp *tcpClient) connect(router *router.Switch) {
	retryDelay := RETRY_MIN_DELAY
	retries := 0
	closing := false

	for !closing {
		infof("TCP", "connecting to %v", tcp.addr)

		if socket, err := net.Dial("tcp", fmt.Sprintf("%v", tcp.addr)); err != nil {
			warnf("TCP", "%v", err)
		} else if socket == nil {
			warnf("TCP", "connect %v failed (%v)", tcp.addr, socket)
		} else {
			retries = 0
			retryDelay = RETRY_MIN_DELAY
			eof := make(chan struct{})

			go func() {
				for {
					select {
					case msg := <-tcp.ch:
						infof("TCP", "msg %v  relaying to %v", msg.ID, socket.RemoteAddr())
						tcp.send(socket, msg.ID, msg.Message)

					case <-eof:
						return

					case <-tcp.closing:
						closing = true
						socket.Close()
						return
					}
				}
			}()

			if err := tcp.listen(socket, router); err != nil {
				if err == io.EOF {
					warnf("TCP", "connection to %v closed ", socket.RemoteAddr())
				} else {
					warnf("TCP", "connection to %v error (%v)", tcp.addr, err)
				}
			}

			close(eof)
		}

		if closing {
			break
		}

		// ... retry
		retries++
		if tcp.maxRetries >= 0 && retries > tcp.maxRetries {
			warnf("TCP", "Connect to %v failed (retry count exceeded %v)", tcp.addr, tcp.maxRetries)
			return
		}

		infof("TCP", "connection failed ... retrying in %v", retryDelay)
		time.Sleep(retryDelay)

		retryDelay *= 2
		if retryDelay > tcp.maxRetryDelay {
			retryDelay = tcp.maxRetryDelay
		}
	}
}

func (tcp *tcpClient) listen(socket net.Conn, router *router.Switch) error {
	infof("TCP", "connected  to %v", socket.RemoteAddr())

	defer socket.Close()

	for {
		buffer := make([]byte, 2048)
		N, err := socket.Read(buffer)
		if err != nil {
			return err
		}

		go func() {
			tcp.received(buffer[:N], router, socket)
		}()
	}
}

func (tcp *tcpClient) received(buffer []byte, router *router.Switch, socket net.Conn) {
	dumpf("TCP", buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

	for len(buffer) > 0 {
		id, msg, remaining := depacketize(buffer)
		buffer = remaining

		router.Received(id, msg, func(message []byte) {
			tcp.send(socket, id, message)
		})
	}
}

func (tcp *tcpClient) send(conn net.Conn, id uint32, msg []byte) []byte {
	packet := packetize(id, msg)

	if N, err := conn.Write(packet); err != nil {
		warnf("TCP", "msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		warnf("TCP", "msg %v  sent %v of %v bytes to %v", id, N, len(msg), conn.RemoteAddr())
	} else {
		infof("TCP", "msg %v  sent %v bytes to %v", id, len(msg), conn.RemoteAddr())
	}

	return nil
}
