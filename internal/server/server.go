package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var serverCounter = 0

func main() {
	listener, errL := net.Listen("tcp", "localhost:3002")
	if errL != nil {
		log.Println(errL)
		return
	}
	defer listener.Close()

	local, err := net.ResolveUDPAddr("udp", "10.50.0.113:1234")
	if err != nil {
		log.Println(err)
		return
	}

	remote, err := net.ResolveUDPAddr("udp", "234.50.100.11:6000")
	if err != nil {
		log.Println(err)
		return
	}

	udpConn, err1 := net.DialUDP("udp", local, remote)

	if err1 != nil {
		log.Println(err1)
		return
	}
	defer udpConn.Close()

	defer fmt.Print("\n")

	for {
		conn, errC := listener.Accept()
		if errC != nil {
			log.Println(errC)
			return
		}
		go handlRequest(conn, udpConn)
	}

}

func handlRequest(conn net.Conn, udpConn *net.UDPConn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	body := make([]byte, 65_800)
	size, errM := reader.Read(body)

	if errM != nil {
		conn.Close()
		log.Print("Request closed\n\n")
		return
	}
	// log.Println(string(body), "\n")
	go sendUdpData(udpConn, body[0:size])
	// saveToFile(body[0:size])

	fmt.Printf("Request size: %d Served: %d\r", size, serverCounter)

	serverCounter++
	conn.Write([]byte("Received" + fmt.Sprint(serverCounter)))
}

func sendUdpData(udpConn *net.UDPConn, body []byte) {
	start := 0
	for i := range 65_800 {
		if i%1316 == 0 {
			_, errUDP := udpConn.Write(body[start:i])

			if errUDP != nil {
				log.Println(errUDP)
			}
		}
	}
}

func saveToFile(data []byte) error {
	file, err := os.OpenFile("teste.bin", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, errW := file.Write(data)
	if errW != nil {
		return errW
	}
	return nil
}
