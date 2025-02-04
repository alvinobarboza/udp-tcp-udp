package tcp

import "net"

type TCPClient interface {
	Write([]byte) ([]byte, error)
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
	tcpAddr *net.TCPAddr
}

func (tcp *tcpClient) Write(datagram []byte) ([]byte, error) {
	conn, errD := net.DialTCP("tcp", nil, tcp.tcpAddr)
	if errD != nil {
		return nil, errD
	}

	_, err1 := conn.Write(datagram)
	if err1 != nil {
		return nil, err1
	}

	reply := make([]byte, 0)

	_, err2 := conn.Read(reply)
	if err2 != nil {
		return nil, err2
	}

	conn.Close()
	return reply, nil
}
