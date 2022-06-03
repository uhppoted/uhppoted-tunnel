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
	ch            chan []byte
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
		ch:            make(chan []byte),
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

func (tcp *tcpClient) Send(id uint32, message []byte) {
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

			// go func() {
			// 	for {
			// 		msg := <-tcp.ch
			// 		infof("TCP  relaying %v bytes to %v", len(msg), socket.RemoteAddr())
			// 		tcp.send(socket, msg)
			// 	}
			// }()

			if err := tcp.listen(socket, router); err != nil {
				if err == io.EOF {
					warnf("TCP  connection to %v closed ", socket.RemoteAddr())
				} else {
					warnf("TCP  connection to %v error (%v)", tcp.addr, err)
				}
			}
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

	buffer := make([]byte, 2048)

	for {
		N, err := socket.Read(buffer)
		if err != nil {
			return err
		}

		hex := dump(buffer[:N], "                                ")
		debugf("TCP  received %v bytes from %v\n%s\n", N, socket.RemoteAddr(), hex)

		ix := 0
		for ix < N {
			size := uint(buffer[ix])
			size <<= 8
			size += uint(buffer[ix+1])

			id, message := depacketize(buffer[ix:])

			h := func(message []byte) {
				tcp.send(socket, id, message)
			}

			router.request(id, message, h)

			// if reply := relay(id, message); reply != nil && len(reply) > 0 {
			// 	packet := packetize(id, reply)
			//
			// 	if N, err := socket.Write(packet); err != nil {
			// 		warnf("error relaying reply to %v (%v)", socket.RemoteAddr(), err)
			// 	} else if N != len(packet) {
			// 		warnf("relayed reply with %v of %v bytes to %v", N, len(reply), socket.RemoteAddr())
			// 	} else {
			// 		infof("relayed reply with %v bytes to %v", len(reply), socket.RemoteAddr())
			// 	}
			// }

			ix += 6 + int(size)
		}
	}
}

func (tcp *tcpClient) send(conn net.Conn, id uint32, message []byte) []byte {
	packet := packetize(id, message)

	if N, err := conn.Write(packet); err != nil {
		warnf("error sending message to %v (%v)", conn.RemoteAddr(), err)
	} else if N != len(packet) {
		warnf("TCP  sent %v of %v bytes to %v", N, len(message), conn.RemoteAddr())
	} else {
		infof("TCP  sent %v bytes to %v", len(message), conn.RemoteAddr())
	}

	return nil
}
