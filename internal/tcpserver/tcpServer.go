package tcpserver

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/alvinobarboza/udp-tcp-udp/internal/args"
	"github.com/alvinobarboza/udp-tcp-udp/internal/utils"
)

type TCPServer interface {
	Listen() error
}

type tcpServer struct {
	wk      Worker
	ipAddr  string
	counter uint32
	mu      sync.Mutex
}

func NewTCPServer(ipaddr string, wk Worker) TCPServer {
	return &tcpServer{
		wk:     wk,
		ipAddr: ipaddr,
	}
}

func (ts *tcpServer) Listen() error {
	tcpAddr, errt := net.ResolveTCPAddr("tcp", ts.ipAddr)
	if errt != nil {
		return errt
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Println("Listening on:", ts.ipAddr)
	go ts.wk.Start()

	for {
		conn, errC := listener.AcceptTCP()
		if errC != nil {
			return errC
		}
		if err := conn.SetReadBuffer(1024 * 1024 * 20); err != nil {
			return err
		}

		go ts.handlRequest(conn)
	}
}

func (ts *tcpServer) handlRequest(conn *net.TCPConn) {
	defer conn.Close()
	ts.mu.Lock()
	ts.counter++
	local := ts.counter
	ts.mu.Unlock()

	header := make([]byte, 8)
	_, errH := conn.Read(header)

	if errH != nil {
		fmt.Println()
		fmt.Println()
		log.Println(errH, local, header)
		fmt.Println()
		return
	}

	if eofFromClient(header) {
		fmt.Println()
		log.Println("Closed early by client", local)
		fmt.Println()
		return
	}

	tcpBuf := &utils.TCPBuffData{
		Counter: binary.LittleEndian.Uint64(header),
	}

	reply := []byte("ok")
	data := make([]byte, args.MPEGTS_PKT_DEFAULT)
	for {
		dRead, errR := conn.Read(data)
		if errR != nil {
			fmt.Println()
			log.Println(errR, "req: ", tcpBuf.Counter)
			fmt.Println()
			return
		}
		if dRead == 5 && eofFromClient(data) {
			fmt.Println()
			log.Println("Closed betwen transmission", local, tcpBuf.Counter)
			fmt.Println()
			return
		}
		if dRead == 2 {
			break
		}
		tcpBuf.Data = append(tcpBuf.Data, data[:dRead]...)
		conn.Write(reply)
	}
	log.Println(len(tcpBuf.Data), "enqueue")
	ts.wk.Enqueue(tcpBuf)

	conn.Write(reply)
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
