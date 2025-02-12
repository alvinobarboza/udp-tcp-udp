package tcpserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/alvinobarboza/udp-tcp-udp/internal/filehandler"
)

type TCPServer interface {
	Listen() error
}

type tcpServer struct {
	udp    UDPSender
	file   filehandler.FileHandler
	ipAddr string
}

func NewTCPServer(ipaddr string, udp UDPSender, file filehandler.FileHandler) TCPServer {
	return &tcpServer{
		udp:    udp,
		file:   file,
		ipAddr: ipaddr,
	}
}

func (ts *tcpServer) Listen() error {
	defer ts.udp.CloseConn()

	listener, err := net.Listen("tcp", ts.ipAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Println("Listening on:", ts.ipAddr)

	for {
		conn, errC := listener.Accept()
		if errC != nil {
			return errC
		}

		go ts.handlRequest(conn)
	}
}

func (ts *tcpServer) handlRequest(conn net.Conn, err chan error) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(time.Second * 3))

	reader := bufio.NewReader(conn)
	body := make([]byte, 1500)
	read, errM := reader.Read(body)

	if errM != nil {
		conn.Close()
		err <- errM
		return
	}
	fmt.Printf("Size: %02d body: %02d\n", read, body[:read])

	// ts.pktWriter.Write(body[0:read], err)

	conn.Write([]byte("Received"))
}
