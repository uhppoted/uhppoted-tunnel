package tunnel

import (
	"fmt"
	"io"
	"net"
	"time"
)

type tcpClient struct {
	addr          *net.TCPAddr
	maxRetries    int
	maxRetryDelay time.Duration
	timeout       time.Duration
	mode          Mode
	ch            chan message
}

const RETRY_MIN_DELAY = 5 * time.Second

func NewTCPClient(spec string, maxRetries int, maxRetryDelay time.Duration, mode Mode) (*tcpClient, error) {
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
		mode:          mode,
		ch:            make(chan message),
	}

	return &in, nil
}

func (tcp *tcpClient) Close() {

}

func (tcp *tcpClient) Run(relay relay) error {
	router := Switch{
		relay: relay,
	}

	return tcp.connect(&router)
}

func (tcp *tcpClient) Send(id uint32, msg []byte) {
	select {
	case tcp.ch <- message{id: id, message: msg}:
	default:
	}
}

func (tcp *tcpClient) connect(router *Switch) error {
	retryDelay := RETRY_MIN_DELAY
	retries := 0

	for tcp.maxRetries < 0 || retries < tcp.maxRetries {
		infof("TCP  connecting to %v", tcp.addr)

		if socket, err := net.Dial("tcp", fmt.Sprintf("%v", tcp.addr)); err != nil {
			warnf("TCP  %v", err)
		} else if socket == nil {
			warnf("TCP  connect %v failed (%v)", tcp.addr, socket)
		} else {
			retries = 0
			retryDelay = RETRY_MIN_DELAY
			eof := make(chan struct{})

			go func() {
				for {
					select {
					case msg := <-tcp.ch:
						infof("TCP  relaying message %v to %v", msg.id, socket.RemoteAddr())
						tcp.send(socket, msg.id, msg.message)

					case <-eof:
						return
					}
				}
			}()

			if err := tcp.listen(socket, router); err != nil {
				if err == io.EOF {
					warnf("TCP  connection to %v closed ", socket.RemoteAddr())
				} else {
					warnf("TCP  connection to %v error (%v)", tcp.addr, err)
				}
			}

			close(eof)
		}

		infof("TCP  connection failed ... retrying in %v", retryDelay)

		time.Sleep(retryDelay)

		retries++
		retryDelay *= 2
		if retryDelay > tcp.maxRetryDelay {
			retryDelay = tcp.maxRetryDelay
		}
	}

	return fmt.Errorf("Connect to %v failed (retry count exceeded %v)", tcp.addr, tcp.maxRetries)
}

func (tcp *tcpClient) listen(socket net.Conn, router *Switch) error {
	infof("TCP  connected  to %v", socket.RemoteAddr())

	defer socket.Close()

	for {
		buffer := make([]byte, 2048)
		N, err := socket.Read(buffer)
		if err != nil {
			return err
		}

		tcp.received(buffer[:N], router, socket)
	}
}

func (tcp *tcpClient) received(buffer []byte, router *Switch, socket net.Conn) {
	hex := dump(buffer, "                                ")
	debugf("TCP  received %v bytes from %v\n%s\n", len(buffer), socket.RemoteAddr(), hex)

	for len(buffer) > 0 {
		id, msg, remaining := depacketize(buffer)

		h := func(message []byte) {
			tcp.send(socket, id, message)
		}

		switch tcp.mode {
		case ModeNormal:
			router.request(id, msg, h)

		case ModeReverse:
			router.reply(id, msg)
		}

		buffer = remaining
	}
}

func (tcp *tcpClient) send(conn net.Conn, id uint32, msg []byte) []byte {
	packet := packetize(id, msg)

	if N, err := conn.Write(packet); err != nil {
		warnf("TCP  msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		warnf("TCP  msg %v  sent %v of %v bytes to %v", id, N, len(msg), conn.RemoteAddr())
	} else {
		infof("TCP  msg %v  sent %v bytes to %v", id, len(msg), conn.RemoteAddr())
	}

	return nil
}
