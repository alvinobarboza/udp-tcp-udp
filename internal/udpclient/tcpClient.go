package udpclient

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/alvinobarboza/udp-tcp-udp/internal/args"
)

type TCPClient interface {
	GetConn() (net.Conn, error)
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
	defer conn.Close()

	tcp.mu.Lock()
	tcp.pktCounter++
	local := tcp.pktCounter
	tcp.mu.Unlock()
	fmt.Println("\nNEW============>", tcp.pktCounter)

	_, err1 := conn.Write(intToByte(local))
	if err1 != nil {
		err <- err1
		return
	}

	pktToSend := make([]byte, 0)
	for i, data := range datagram {
		if i%args.TS_PACKET_DEFAULT == 0 && len(pktToSend) > 0 {
			// _, err1 := conn.Write(pktToSend)
			// if err1 != nil {
			// 	err <- err1
			// 	return
			// }
			fmt.Printf("Curr %02d size %02d\r", local, i)
			pktToSend = make([]byte, 0)
		}
		pktToSend = append(pktToSend, data)
	}

	// size := len(datagram)
	// tcp.pktCounter++

	// fmt.Printf("%02d %06d        \r", tcp.pktCounter, size)

	reply := make([]byte, 1024)

	_, err2 := conn.Read(reply)
	if err2 != nil {
		err <- err2
		return
	}

	// err <- fmt.Errorf("endend")
}

func intToByte(n uint16) []byte {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, n)
	return bs
}
