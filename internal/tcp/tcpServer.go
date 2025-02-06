package tcp

import (
	"bufio"
	"log"
	"net"
)

type PktWriter interface {
	Write([]byte, chan error)
	CloseConn()
}

type TCPServer interface {
	Listen() error
}

type tcpServer struct {
	pktWriter PktWriter
	ipAddr    string
}

func NewTCPServer(ipaddr string, pktWriter PktWriter) TCPServer {
	return &tcpServer{
		pktWriter: pktWriter,
		ipAddr:    ipaddr,
	}
}

func (ts *tcpServer) Listen() error {
	defer ts.pktWriter.CloseConn()

	listener, err := net.Listen("tcp", ts.ipAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Println("Listening on:", ts.ipAddr)

	errChann := make(chan error)

	for {
		select {
		case errC := <-errChann:
			return errC
		default:
			conn, errC := listener.Accept()
			if errC != nil {
				return errC
			}

			go ts.handlRequest(conn, errChann)
		}
	}
}

func (ts *tcpServer) handlRequest(conn net.Conn, err chan error) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	body := make([]byte, 65_800)
	read, errM := reader.Read(body)

	if errM != nil {
		conn.Close()
		err <- errM
		return
	}

	ts.pktWriter.Write(body[0:read], err)

	conn.Write([]byte("Received"))
}
