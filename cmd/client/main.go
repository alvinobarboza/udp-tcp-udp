package main

import (
	"fmt"

	"github.com/alvinobarboza/udp-tcp-udp/internal/filehandler"
	"github.com/alvinobarboza/udp-tcp-udp/internal/udp"
)

func main() {
	defer fmt.Println("Exit                     ")
	packetSize := 2000
	tcpMultiplierBuf := 50
	timerSeconds := 10

	file := filehandler.NewFileHandler()
	if err := file.NewFile("as/teste.bin"); err != nil {
		panic(err)
	}

	udpListener := udp.NewUDPListener(tcpMultiplierBuf, packetSize, timerSeconds, file)

	if err := udpListener.SetUpListener("Ethernet", "234.50.99.3:6000"); err != nil {
		panic(err)
	}

	if errU := udpListener.Listen(); errU != nil {
		panic(errU)
	}
}
