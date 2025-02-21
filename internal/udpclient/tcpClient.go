package udpclient

import (
	"fmt"
	"net"

	"github.com/alvinobarboza/udp-tcp-udp/internal/args"
	"github.com/alvinobarboza/udp-tcp-udp/internal/utils"
)

type TCPClient interface {
	Write(*utils.TCPBuffData, chan error, chan string)
	GetConn() (*net.TCPConn, error)
	Close(*net.TCPConn) error
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

func (tcp *tcpClient) GetConn() (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp", nil, tcp.tcpAddr)
	if err != nil {
		return nil, err
	}
	if err := conn.SetWriteBuffer(1024 * 1024 * 20); err != nil {
		return nil, err
	}
	return conn, nil
}

func (tcp *tcpClient) Close(conn *net.TCPConn) error {
	_, err2 := conn.Write([]byte{0xff, 0xff, 0xff, 0xff, 0xff})
	if err2 != nil {
		return err2
	}
	return conn.Close()
}

func (tcp *tcpClient) Write(
	datagram *utils.TCPBuffData,
	err chan error, t chan string,
) {
	conn, errC := tcp.GetConn()
	if errC != nil {
		err <- errC
		return
	}
	defer conn.Close()

	go func(c *net.TCPConn, terminator chan string) {
		<-terminator
		tcp.Close(c)
	}(conn, t)

	header := headerData(datagram.Counter)
	fmt.Println("Conn:", datagram.Counter, len(header), len(datagram.Data))

	_, err1 := conn.Write(header)
	if err1 != nil {
		err <- err1
		return
	}

	reply := make([]byte, 128)
	pktToSend := make([]byte, 0)
	for _, data := range datagram.Data {
		pktToSend = append(pktToSend, data)
		if len(pktToSend) == args.MPEGTS_PKT_DEFAULT {
			_, err1 := conn.Write(pktToSend)
			if err1 != nil {
				err <- err1
				return
			}
			pktToSend = make([]byte, 0)
			_, err3 := conn.Read(reply)
			if err3 != nil {
				err <- err3
				return
			}
		}
	}

	if len(pktToSend) > 0 && len(pktToSend) < args.MPEGTS_PKT_DEFAULT {
		_, err1 := conn.Write(pktToSend)
		if err1 != nil {
			err <- err1
			return
		}
	}

	_, err2 := conn.Write([]byte{0xff, 0xff})
	if err2 != nil {
		err <- err2
		return
	}

	_, err3 := conn.Read(reply)
	if err3 != nil {
		err <- err3
		return
	}
	fmt.Println("ended")
}

func headerData(count uint64) []byte {
	return utils.Int64ToByte(count)
}
