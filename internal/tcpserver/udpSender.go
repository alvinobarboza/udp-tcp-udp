package tcpserver

import "net"

type UDPSender interface {
	Write([]byte) error
	CloseConn()
}

func NewUDPSender(local, remote string) (UDPSender, error) {
	localAddr, err1 := net.ResolveUDPAddr("udp", local)
	if err1 != nil {
		return nil, err1
	}

	remoteAddr, err2 := net.ResolveUDPAddr("udp", remote)
	if err2 != nil {
		return nil, err2
	}

	conn, err3 := net.DialUDP("udp", localAddr, remoteAddr)
	if err3 != nil {
		return nil, err3
	}

	return &udpSender{
		local:  localAddr,
		remote: remoteAddr,
		conn:   conn,
	}, nil
}

type udpSender struct {
	local  *net.UDPAddr
	remote *net.UDPAddr
	conn   *net.UDPConn
}

func (ud *udpSender) Write(data []byte) error {
	_, err := ud.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (us *udpSender) CloseConn() {
	us.conn.Close()
}
