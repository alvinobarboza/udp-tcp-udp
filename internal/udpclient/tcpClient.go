package udpclient

import (
	"fmt"
	"net"

	"github.com/alvinobarboza/udp-tcp-udp/internal/args"
	"github.com/alvinobarboza/udp-tcp-udp/internal/utils"
)

type TCPClient interface {
	Write(*utils.TCPBuffData, chan error, chan string)
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
	tcpAddr *net.TCPAddr
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
	datagram *utils.TCPBuffData,
	err chan error, t chan string,
) {
	conn, errC := tcp.GetConn()
	if errC != nil {
		err <- errC
		return
	}
	defer conn.Close()

	go func(c net.Conn, terminator chan string) {
		<-terminator
		tcp.Close(c)
	}(conn, t)

	header := headerData(
		datagram.Counter,
		datagram.MS,
	)
	fmt.Println("Conn:", datagram.Counter, datagram.MS, len(header), len(datagram.Data))

	_, err1 := conn.Write(header)
	if err1 != nil {
		err <- err1
		return
	}

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

	reply := make([]byte, 1024)

	_, err3 := conn.Read(reply)
	if err3 != nil {
		err <- err3
		return
	}
}

func headerData(count uint64, ms uint32) []byte {
	return append(
		utils.Int64ToByte(count),
		utils.Int32ToByte(ms)...,
	)
}
