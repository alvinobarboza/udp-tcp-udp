package main

import (
	"fmt"

	"github.com/alvinobarboza/udp-tcp-udp/internal/tcp"
	"github.com/alvinobarboza/udp-tcp-udp/internal/udp"
)

func main() {
	//TODO: Args handling
	defer fmt.Println("Exit                     ")
	packetSize := 2000
	tcpMultiplierBuf := 50
	timerSeconds := 10

	// file := filehandler.NewFileHandler()
	// if err := file.NewFile("teste.bin"); err != nil {
	// 	panic(err)
	// }

	tcpCon, errc := tcp.NewTCPClient("localhost:3002")
	if errc != nil {
		panic(errc)
	}

	udpListener := udp.NewUDPListener(tcpMultiplierBuf, packetSize, timerSeconds, tcpCon)

	if err := udpListener.SetUpListener("Ethernet", "234.50.99.2:6000"); err != nil {
		panic(err)
	}

	if errU := udpListener.Listen(); errU != nil {
		panic(errU)
	}
}
