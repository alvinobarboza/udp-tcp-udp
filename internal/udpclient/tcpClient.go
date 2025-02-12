package udpclient

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/alvinobarboza/udp-tcp-udp/internal/args"
)

type TCPClient interface {
	Write(net.Conn, []byte, chan error, chan string)
	GetConn() (net.Conn, error)
	Close(net.Conn) error
}

func NewTCPClient(servAddr string) (TCPClient, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return nil, err
	}

	return &tcpClient{
		tcpAddr: tcpAddr,
	}, nil
}

type tcpClient struct {
	mu         sync.Mutex
	tcpAddr    *net.TCPAddr
	pktCounter uint16
}

func (tcp *tcpClient) GetConn() (net.Conn, error) {
	conn, err := net.DialTCP("tcp", nil, tcp.tcpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (tcp *tcpClient) Close(conn net.Conn) error {
	_, err2 := conn.Write([]byte{0xff, 0xff, 0xff, 0xff, 0xff})
	if err2 != nil {
		return err2
	}
	return conn.Close()
}

func (tcp *tcpClient) Write(
	conn net.Conn, datagram []byte,
	err chan error, t chan string,
) {
	defer conn.Close()

	go func(c net.Conn, terminator chan string) {
		value := <-terminator
		if value == "" {
			tcp.Close(c)
		}
	}(conn, t)

	tcp.mu.Lock()
	tcp.pktCounter++
	local := tcp.pktCounter
	tcp.mu.Unlock()
	fmt.Println("Conn n:", local)

	header := append(intToByte(local), intToByte(uint16(args.MPEGTS_PKT_DEFAULT))...)

	_, err1 := conn.Write(header)
	if err1 != nil {
		err <- err1
		return
	}

	pktToSend := make([]byte, 0)
	for i, data := range datagram {
		if i%args.MPEGTS_PKT_DEFAULT == 0 && len(pktToSend) > 0 {
			_, err1 := conn.Write(pktToSend)
			if err1 != nil {
				err <- err1
				return
			}
			// if local%2 == 0 {
			// 	fmt.Printf("\t\t")
			// }
			// fmt.Printf("Curr %02d %d\r", local, len(pktToSend))
			pktToSend = make([]byte, 0)
		}
		pktToSend = append(pktToSend, data)
	}

	_, err2 := conn.Write([]byte{0xff, 0xff})
	if err2 != nil {
		err <- err2
		return
	}

	reply := make([]byte, 1024)

	_, err3 := conn.Read(reply)
	if err3 != nil {
		err <- err3
		return
	}
	t <- "ok"
}

func intToByte(n uint16) []byte {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, n)
	return bs
}
