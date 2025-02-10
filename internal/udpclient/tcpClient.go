package udpclient

import (
	"fmt"
	"net"
)

type TCPClient interface {
	Write([]byte, chan error)
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
	tcpAddr    *net.TCPAddr
	pktCounter uint8
}

func (tcp *tcpClient) Write(datagram []byte, err chan error) {
	conn, errD := net.DialTCP("tcp", nil, tcp.tcpAddr)
	if errD != nil {
		err <- errD
		return
	}
	defer conn.Close()

	size := len(datagram)
	tcp.pktCounter++

	fmt.Printf("%02d %06d        \r", tcp.pktCounter, size)

	_, err1 := conn.Write(datagram)
	if err1 != nil {
		err <- err1
		return
	}

	reply := make([]byte, 0)

	_, err2 := conn.Read(reply)
	if err2 != nil {
		err <- err2
		return
	}

	// err <- fmt.Errorf("Endend")
}
