package tcpserver

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/alvinobarboza/udp-tcp-udp/internal/filehandler"
)

type TCPServer interface {
	Listen() error
}

type tcpServer struct {
	udp     UDPSender
	file    filehandler.FileHandler
	ipAddr  string
	counter uint32
	mu      sync.Mutex
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

func (ts *tcpServer) handlRequest(conn net.Conn) {
	defer conn.Close()
	ts.mu.Lock()
	ts.counter++
	local := ts.counter
	ts.mu.Unlock()

	header := make([]byte, 5)
	_, errH := conn.Read(header)

	if errH != nil {
		log.Println(errH, local, header)
		return
	}

	if eofFromClient(header) {
		log.Println("Closed early by client", local)
		return
	}

	pCount := binary.LittleEndian.Uint16(header[0:2])
	pSize := binary.LittleEndian.Uint16(header[2:4])
	fmt.Printf("\nSize: %02d counter: %02d - %d\n", pSize, pCount, header)

	data := make([]byte, pSize)
	for {
		dRead, errR := conn.Read(data)

		if errR != nil {
			log.Println(errR, "req: ", pCount)
			return
		}
		if dRead == 5 && eofFromClient(data) {
			log.Println("Closed early by client", local, pCount)
			return
		}
		if dRead == 2 {
			fmt.Println("\nClosed", data[:dRead])
			break
		}
		fmt.Printf("read: %02d\r", dRead)
	}

	// ts.pktWriter.Write(body[0:read], err)

	conn.Write([]byte("Received"))
}

func eofFromClient(data []byte) bool {
	checkBound := 5
	counter := 0
	for i, d := range data {
		if i > checkBound {
			break
		}
		if d == 0xff {
			counter++
		}
	}
	return counter == checkBound
}
